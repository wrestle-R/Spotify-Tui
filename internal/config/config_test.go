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
