package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func packageFixture(t *testing.T, version string) []byte {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	binary, err := os.ReadFile(self)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	gz := gzip.NewWriter(&out)
	tw := tar.NewWriter(gz)
	root := strings.TrimSuffix(archiveName(version, runtime.GOOS, runtime.GOARCH), ".tar.gz") + "/"
	for _, name := range []string{"spotify-tui", "spotify-tui-playback"} {
		if err = tw.WriteHeader(&tar.Header{Name: root + name, Mode: 0755, Size: int64(len(binary))}); err != nil {
			t.Fatal(err)
		}
		_, _ = tw.Write(binary)
	}
	if err = tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err = gz.Close(); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}
func installed(t *testing.T) (string, string) {
	t.Helper()
	prefix := filepath.Join(t.TempDir(), "prefix with spaces")
	stage := filepath.Join(prefix, "lib", "spotify-tui", "build.old")
	if err := os.MkdirAll(stage, 0700); err != nil {
		t.Fatal(err)
	}
	_ = os.MkdirAll(filepath.Join(prefix, "bin"), 0700)
	for _, name := range []string{"spotify-tui", "spotify-tui-playback"} {
		if err := os.WriteFile(filepath.Join(stage, name), []byte("old build"), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(stage, name), filepath.Join(prefix, "bin", name)); err != nil {
			t.Fatal(err)
		}
	}
	canonical, err := filepath.EvalSymlinks(prefix)
	if err != nil {
		t.Fatal(err)
	}
	return canonical, filepath.Join(canonical, "lib", "spotify-tui", "build.old", "spotify-tui")
}
func fixtureClient(t *testing.T, corrupt bool) (*Client, string) {
	prefix, exe := installed(t)
	archive := packageFixture(t, "0.3.0")
	checksum := fmt.Sprintf("%x  %s\n", sha256.Sum256(archive), archiveName("0.3.0", runtime.GOOS, runtime.GOARCH))
	if corrupt {
		archive = append(archive, 'x')
	}
	c := &Client{Current: "0.2.1", Executable: exe, Fetch: func(ctx context.Context, url string, w io.Writer, limit int64) error {
		var data []byte
		if strings.HasSuffix(url, "SHA256SUMS") {
			data = []byte(checksum)
		} else {
			data = archive
		}
		_, err := w.Write(data)
		return err
	}, Run: func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return []byte(filepath.Base(name) + " 0.3.0\n"), nil
	}}
	return c, prefix
}
func TestVerifiedPairActivatesTogetherAndRetainsOldBuild(t *testing.T) {
	c, prefix := fixtureClient(t, false)
	old := filepath.Dir(c.Executable)
	result, err := c.Install(context.Background(), Release{Tag: "v0.3.0"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Version != "0.3.0" || result.Executable != filepath.Join(prefix, "bin", "spotify-tui") {
		t.Fatal(result)
	}
	var stage string
	for _, name := range []string{"spotify-tui", "spotify-tui-playback"} {
		link := filepath.Join(prefix, "bin", name)
		target, err := filepath.EvalSymlinks(link)
		if err != nil {
			t.Fatal(err)
		}
		if stage != "" && stage != filepath.Dir(target) {
			t.Fatal("mixed pair")
		}
		stage = filepath.Dir(target)
		raw, _ := os.Readlink(link)
		if !strings.Contains(raw, "/current/") {
			t.Fatal(raw)
		}
	}
	if stage == old {
		t.Fatal("old build still active")
	}
	if _, err = os.Stat(filepath.Join(old, "spotify-tui")); err != nil {
		t.Fatal("rollback build lost")
	}
	if temp, _ := filepath.Glob(filepath.Join(prefix, "lib", "spotify-tui", ".update-*")); len(temp) != 0 {
		t.Fatal("temporary downloads leaked")
	}
}
func TestCorruptionLeavesInstallationUntouched(t *testing.T) {
	c, prefix := fixtureClient(t, true)
	before, _ := os.Readlink(filepath.Join(prefix, "bin", "spotify-tui"))
	if _, err := c.Install(context.Background(), Release{Tag: "v0.3.0"}); err == nil {
		t.Fatal("accepted corrupt archive")
	}
	after, _ := os.Readlink(filepath.Join(prefix, "bin", "spotify-tui"))
	if after != before {
		t.Fatal("changed installation on failure")
	}
}
func TestUnsafeArchivesRejected(t *testing.T) {
	for _, name := range []string{"root/../escape", "/root/file", "other/file", "root/./file"} {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)
			tw := tar.NewWriter(gz)
			_ = tw.WriteHeader(&tar.Header{Name: name, Typeflag: tar.TypeReg, Size: 0})
			_ = tw.Close()
			_ = gz.Close()
			p := filepath.Join(t.TempDir(), "bad.gz")
			_ = os.WriteFile(p, b.Bytes(), 0600)
			if err := extract(p, t.TempDir(), "root"); err == nil {
				t.Fatal("unsafe path accepted")
			}
		})
	}
	for _, kind := range []byte{tar.TypeSymlink, tar.TypeLink, tar.TypeChar} {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		tw := tar.NewWriter(gz)
		_ = tw.WriteHeader(&tar.Header{Name: "root/file", Typeflag: kind, Linkname: "/tmp/outside"})
		_ = tw.Close()
		_ = gz.Close()
		p := filepath.Join(t.TempDir(), "bad.gz")
		_ = os.WriteFile(p, b.Bytes(), 0600)
		if extract(p, t.TempDir(), "root") == nil {
			t.Fatal("unsafe file accepted")
		}
	}
}
func TestStableCheckAndNumericVersions(t *testing.T) {
	if !Newer("v0.10.0", "0.9.9") || Newer("0.3.0-beta.1", "0.2.1") || Newer("0.2.1", "0.3.0") {
		t.Fatal("version ordering")
	}
	c := New("0.2.1")
	c.Fetch = func(ctx context.Context, url string, w io.Writer, limit int64) error {
		_, err := fmt.Fprintf(w, `{"tag_name":"v0.3.0","assets":[{"name":%q},{"name":"SHA256SUMS"}]}`, archiveName("0.3.0", runtime.GOOS, runtime.GOARCH))
		return err
	}
	r, err := c.Check(context.Background())
	if err != nil || r == nil || r.Tag != "v0.3.0" {
		t.Fatal(r, err)
	}
	c.Current = "0.3.0"
	r, err = c.Check(context.Background())
	if err != nil || r != nil {
		t.Fatal("same release offered")
	}
	c.Current = "0.2.1"
	c.Fetch = func(ctx context.Context, url string, w io.Writer, limit int64) error {
		_, err := io.WriteString(w, `{"tag_name":"v9.0.0","prerelease":true}`)
		return err
	}
	if _, err = c.Check(context.Background()); err == nil {
		t.Fatal("preview offered")
	}
}
func TestUnmanagedAndMixedInstallsAreRefused(t *testing.T) {
	if _, err := ManagedPrefix("/bin/true"); err == nil {
		t.Fatal("unmanaged allowed")
	}
	prefix, exe := installed(t)
	link := filepath.Join(prefix, "bin", "spotify-tui-playback")
	_ = os.Remove(link)
	_ = os.Symlink("/bin/true", link)
	if _, err := ManagedPrefix(exe); err == nil {
		t.Fatal("unrelated symlink allowed")
	}
}
