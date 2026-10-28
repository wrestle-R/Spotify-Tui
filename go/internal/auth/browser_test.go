package auth

import "testing"

func TestBrowserLauncherUsesNativeMacCommand(t *testing.T) {
	if browserExecutable("darwin") != "open" || browserExecutable("linux") != "xdg-open" {
		t.Fatal("wrong browser command")
	}
}
