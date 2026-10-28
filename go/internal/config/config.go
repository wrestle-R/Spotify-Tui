package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type Config struct {
	ClientID       string `json:"client_id"`
	RedirectURI    string `json:"redirect_uri"`
	SoloistPath    string `json:"soloist_path"`
	HelperPath     string `json:"helper_path"`
	AllowFileToken bool   `json:"allow_file_token"`
}

type Paths struct{ Config, Data, Cache, State, Runtime string }

func PathsForUser() (Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Paths{}, err
	}
	xdg := func(key, fallback string) string {
		if p := os.Getenv(key); filepath.IsAbs(p) {
			return filepath.Join(p, "spotify-tui")
		}
		return filepath.Join(home, fallback, "spotify-tui")
	}
	p := Paths{Config: xdg("XDG_CONFIG_HOME", ".config"), Data: xdg("XDG_DATA_HOME", ".local/share"), Cache: xdg("XDG_CACHE_HOME", ".cache"), State: xdg("XDG_STATE_HOME", ".local/state")}
	if runtimeDir := os.Getenv("XDG_RUNTIME_DIR"); filepath.IsAbs(runtimeDir) {
		p.Runtime = filepath.Join(runtimeDir, "spotify-tui")
	} else if runtime.GOOS == "darwin" {
		// Darwin Unix sockets have a short path limit; HOME can be very long.
		p.Runtime = filepath.Join(os.TempDir(), "stui-"+strconv.Itoa(os.Getuid()))
	} else {
		p.Runtime = filepath.Join(p.State, "run")
	}
	return p, nil
}

func (p Paths) Ensure() error {
	for _, dir := range []string{p.Config, p.Data, p.Cache, p.State, p.Runtime} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
		info, err := os.Lstat(dir)
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("private directory is not a real directory: %s", dir)
		}
		if err := os.Chmod(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func Load(p Paths) (Config, error) {
	c := Config{RedirectURI: "http://127.0.0.1:43821/callback", SoloistPath: "soloist", HelperPath: "spotify-tui-playback"}
	b, err := os.ReadFile(filepath.Join(p.Config, "config.json"))
	if os.IsNotExist(err) {
		return c, nil
	}
	if err != nil {
		return c, err
	}
	if err = json.Unmarshal(b, &c); err != nil {
		return c, fmt.Errorf("invalid config: %w", err)
	}
	return c, nil
}

func Save(p Paths, c Config) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWrite(filepath.Join(p.Config, "config.json"), append(b, '\n'))
}

func AtomicWrite(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".write-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	if err = f.Chmod(0600); err == nil {
		_, err = f.Write(b)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(name, path)
}
