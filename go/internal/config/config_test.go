package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWriteReplacesFileAndUsesPrivatePermissions(t *testing.T) {
	p := filepath.Join(t.TempDir(), "private", "config.json")
	if err := AtomicWrite(p, []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := AtomicWrite(p, []byte("second")); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p)
	if err != nil || string(b) != "second" {
		t.Fatal("replacement failed")
	}
	st, err := os.Stat(p)
	if err != nil || st.Mode().Perm() != 0600 {
		t.Fatal("file is not private")
	}
}

func TestEnsureRejectsSymlinkDirectory(t *testing.T) {
	dir := t.TempDir()
	link := filepath.Join(dir, "link")
	if err := os.Symlink(t.TempDir(), link); err != nil {
		t.Fatal(err)
	}
	p := Paths{Config: link, Data: dir, Cache: dir, State: dir, Runtime: dir}
	if err := p.Ensure(); err == nil {
		t.Fatal("symlink accepted")
	}
}

func TestRuntimePathHonoursExplicitXDGAndStaysPrivate(t *testing.T) {
	runtimeDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", runtimeDir)
	paths, err := PathsForUser()
	if err != nil {
		t.Fatal(err)
	}
	if paths.Runtime != filepath.Join(runtimeDir, "spotify-tui") {
		t.Fatal(paths.Runtime)
	}
	if err = paths.Ensure(); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(paths.Runtime)
	if err != nil || info.Mode().Perm() != 0700 {
		t.Fatal("runtime directory not private")
	}
}
