package ipc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/wrestle-R/Spotify-Tui/internal/auth"
	"github.com/wrestle-R/Spotify-Tui/internal/config"
)

// OfflineProbe isolates only this app's playback processes, then restores the
// normal online helper. It does not download media or change host networking.
func (c *Client) OfflineProbe(ctx context.Context, cfg config.Config, p config.Paths, vault auth.Vault, uri string, out io.Writer) error {
	state, err := c.State(ctx)
	if err != nil {
		return err
	}
	if !state.Connected || !state.LoggedIn {
		return errors.New("offline probe requires a paired, working online session first; no processes were stopped")
	}
	if _, err := exec.LookPath("unshare"); err != nil {
		return errors.New("unshare is needed for app-scoped offline testing")
	}
	if _, err := exec.LookPath("ip"); err != nil {
		return errors.New("iproute2 is needed for the isolated loopback setup")
	}
	preflight := exec.CommandContext(ctx, "unshare", "--user", "--map-root-user", "--net", "true")
	if preflight.Run() != nil {
		return errors.New("unprivileged network namespaces are unavailable; no host network settings were changed")
	}
	helper, err := resolveHelper(cfg.HelperPath)
	if err != nil {
		return err
	}
	provider, err := exec.LookPath(cfg.SoloistPath)
	if err != nil {
		return err
	}
	key, err := vault.Get("soloist")
	if err != nil {
		return err
	}
	if err = c.Stop(ctx); err != nil {
		return errors.New("online helper did not stop; offline probe canceled")
	}
	defer func() {
		restoreCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if e := c.Ensure(restoreCtx, cfg, p, vault); e != nil {
			fmt.Fprintln(out, "Online helper restore requires attention:", e)
		} else {
			fmt.Fprintln(out, "Normal online helper started; use doctor to check pairing.")
		}
	}()
	dir, err := os.MkdirTemp(p.Runtime, "offline-probe-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	probe := &Client{Socket: filepath.Join(dir, "player.sock")}
	args := []string{"--user", "--map-root-user", "--net", "--", "sh", "-c", "ip link set lo up && exec \"$@\"", "spotify-tui-offline", helper, "serve", "--socket", probe.Socket, "--soloist-path", provider, "--provider-data", filepath.Join(p.Data, "soloist"), "--provider-cache", filepath.Join(p.Cache, "soloist")}
	cmd := exec.Command("unshare", args...)
	cmd.Stdin = strings.NewReader(key + "\n")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err = cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
	}()
	fmt.Fprintln(out, "Playback processes started in a separate network namespace with loopback only; no external interface or default route.")
	ready := false
	for n := 0; n < 150; n++ {
		if e := probe.Call(ctx, "hello", nil, nil); e == nil {
			ready = true
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			done <- nil
			return errors.New("OFFLINE COLD START FAILED: provider exited before becoming ready without external connectivity")
		case <-time.After(100 * time.Millisecond):
		}
	}
	if !ready {
		return errors.New("OFFLINE COLD START FAILED: provider did not become ready without external connectivity")
	}
	state, err = probe.State(ctx)
	if err != nil {
		return err
	}
	for n := 0; !state.LoggedIn && n < 100; n++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
		state, err = probe.State(ctx)
		if err != nil {
			return err
		}
	}
	if !state.LoggedIn {
		return errors.New("OFFLINE COLD START FAILED: stored provider session could not authenticate without external connectivity")
	}
	if err = probe.Call(ctx, "play_track", map[string]any{"uri": uri}, nil); err != nil {
		return fmt.Errorf("OFFLINE PLAYBACK FAILED: %w", err)
	}
	fmt.Fprintln(out, "Provider accepted a cached-track command. This alone does not establish full offline download support.")
	return errors.New("OFFLINE REQUIREMENT STILL UNSATISFIED: no supported music-download operation exists for this provider")
}
