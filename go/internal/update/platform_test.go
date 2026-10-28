package update

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPlatformArchivesAndNativeExecutableValidation(t *testing.T) {
	for _, target := range []struct{ os, arch, archive string }{
		{"linux", "amd64", "spotify-tui_0.4.0_linux_amd64.tar.gz"},
		{"darwin", "arm64", "spotify-tui_0.4.0_darwin_arm64.tar.gz"},
	} {
		if !supportedTarget(target.os, target.arch) || archiveName("0.4.0", target.os, target.arch) != target.archive {
			t.Fatal(target)
		}
	}
	if supportedTarget("darwin", "amd64") || supportedTarget("windows", "amd64") {
		t.Fatal("unsupported playback platform advertised")
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err = verifyExecutable(executable, runtime.GOOS, runtime.GOARCH); err != nil {
		t.Fatal(err)
	}
	wrongOS := "darwin"
	wrongArch := "arm64"
	if runtime.GOOS == "darwin" {
		wrongOS = "linux"
		wrongArch = "amd64"
	}
	if verifyExecutable(executable, wrongOS, wrongArch) == nil {
		t.Fatal("foreign executable accepted")
	}
	fake := filepath.Join(t.TempDir(), "fake")
	if err = os.WriteFile(fake, []byte("not an executable"), 0700); err != nil {
		t.Fatal(err)
	}
	if verifyExecutable(fake, runtime.GOOS, runtime.GOARCH) == nil {
		t.Fatal("invalid executable accepted")
	}
}
