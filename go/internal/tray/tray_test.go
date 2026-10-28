package tray

import (
	"github.com/godbus/dbus/v5"
	"github.com/wrestle-R/Spotify-Tui/go/internal/ipc"
	"testing"
)

func TestTrayPixelsAndMenuMatchHostSignatures(t *testing.T) {
	p, err := pixmaps()
	if err != nil {
		t.Fatal(err)
	}
	if len(p) < 4 {
		t.Fatal("missing tray sizes")
	}
	for _, icon := range p {
		if len(icon.Data) != int(icon.Width*icon.Height*4) {
			t.Fatal("bad ARGB32 size")
		}
		var visible bool
		for i := 0; i < len(icon.Data); i += 4 {
			visible = visible || icon.Data[i] > 0
		}
		if !visible {
			t.Fatal("empty icon")
		}
	}
	if got := dbus.SignatureOf(p).String(); got != "a(iiay)" {
		t.Fatal(got)
	}
	s := &Service{state: ipc.State{Status: "playing"}}
	_, layout, e := s.GetLayout(0, -1, nil)
	if e != nil {
		t.Fatal(e)
	}
	if got := dbus.SignatureOf(layout).String(); got != "(ia{sv}av)" {
		t.Fatal(got)
	}
	if len(layout.Children) != 6 || s.label() != "Pause" {
		t.Fatal("playing menu")
	}
	s.state.Status = "paused"
	if s.label() != "Resume" {
		t.Fatal("paused menu")
	}
}
