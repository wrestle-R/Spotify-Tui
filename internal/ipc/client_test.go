package ipc

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestMismatchedResponseIsRejected(t *testing.T) {
	socket := filepath.Join(t.TempDir(), "socket")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	go func() {
		conn, e := listener.Accept()
		if e != nil {
			return
		}
		defer conn.Close()
		var req any
		_ = json.NewDecoder(conn).Decode(&req)
		_ = json.NewEncoder(conn).Encode(map[string]any{"version": 1, "id": "wrong", "result": true})
	}()
	if err := (&Client{Socket: socket}).Call(context.Background(), "hello", nil, nil); err == nil {
		t.Fatal("mismatched ID accepted")
	}
}

func TestRealRustHelperProtocol(t *testing.T) {
	helper := os.Getenv("SPOTIFY_TUI_TEST_HELPER")
	if helper == "" {
		t.Skip("set SPOTIFY_TUI_TEST_HELPER to the built Rust binary for cross-language checks")
	}
	socket := filepath.Join(t.TempDir(), "player.sock")
	cmd := exec.Command(helper, "serve", "--socket", socket, "--provider-url", "ws://127.0.0.1:9")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
	c := &Client{Socket: socket}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ready := false
	for i := 0; i < 50; i++ {
		if c.Call(ctx, "hello", nil, nil) == nil {
			ready = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !ready {
		t.Fatal("helper never became ready")
	}
	var caps struct {
		Downloads bool   `json:"downloads"`
		Offline   bool   `json:"offline_start"`
		Provider  string `json:"provider"`
	}
	if err := c.Call(ctx, "capabilities", nil, &caps); err != nil || caps.Downloads || caps.Offline || caps.Provider != "soloist" {
		t.Fatalf("incorrect capabilities: %#v %v", caps, err)
	}
	if err := c.Call(ctx, "download_start", nil, nil); err == nil {
		t.Fatal("download unexpectedly succeeded")
	}
	s, err := c.State(ctx)
	if err != nil || s.Connected {
		t.Fatal("disconnected provider reported as connected")
	}
	if err = c.Call(ctx, "shutdown", nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestSharedProtocolFixture(t *testing.T) {
	request, err := os.ReadFile("../../protocol/fixtures/hello-request.json")
	if err != nil {
		t.Fatal(err)
	}
	response, err := os.ReadFile("../../protocol/fixtures/hello-response.json")
	if err != nil {
		t.Fatal(err)
	}
	socket := filepath.Join(t.TempDir(), "socket")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	done := make(chan error, 1)
	go func() {
		conn, e := listener.Accept()
		if e != nil {
			done <- e
			return
		}
		defer conn.Close()
		var got, want map[string]any
		e = json.NewDecoder(conn).Decode(&got)
		if e == nil {
			e = json.Unmarshal(request, &want)
		}
		if e == nil {
			for _, key := range []string{"version", "id", "method"} {
				if got[key] != want[key] {
					t.Errorf("fixture mismatch: %s", key)
				}
			}
		}
		if e == nil {
			_, e = conn.Write(response)
		}
		done <- e
	}()
	var caps map[string]any
	if err := (&Client{Socket: socket}).Call(context.Background(), "hello", nil, &caps); err != nil {
		t.Fatal(err)
	}
	if caps["downloads"] != false {
		t.Fatal("invalid capability result")
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
