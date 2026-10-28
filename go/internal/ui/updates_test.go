package ui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wrestle-R/Spotify-Tui/go/internal/ipc"
	"github.com/wrestle-R/Spotify-Tui/go/internal/update"
	"testing"
)

type updateStub struct{ checks, installs int }

func (u *updateStub) Check(context.Context) (*update.Release, error) {
	u.checks++
	return &update.Release{Tag: "v0.3.0"}, nil
}
func (u *updateStub) Install(context.Context, update.Release) (update.Installation, error) {
	u.installs++
	return update.Installation{Version: "0.3.0", Executable: "/installed/spotify-tui"}, nil
}
func TestCtrlUChecksThenInstallsEvenWhileSearching(t *testing.T) {
	u := &updateStub{}
	m := fixture().WithUpdater(u, "0.2.1")
	m.inputMode = "search"
	m.input = "keep this query"
	m, cmd := key(m, tea.KeyCtrlU)
	if cmd == nil {
		t.Fatal("missing check")
	}
	next, cmd := m.Update(cmd())
	m = next.(Model)
	if !m.updating || cmd == nil || m.input != "keep this query" {
		t.Fatal("Ctrl-U edited input or failed to start installation")
	}
	next, quit := m.Update(cmd())
	m = next.(Model)
	if quit == nil || m.RestartPath() != "/installed/spotify-tui" || u.checks != 1 || u.installs != 1 {
		t.Fatal("update did not restart once")
	}
}
func TestOlderPlayerSnapshotCannotReverseDisplayedPause(t *testing.T) {
	m := fixture()
	m.player = ipc.State{Instance: "session", Revision: 20, Status: "paused"}
	next, _ := m.Update(playerMsg{state: ipc.State{Instance: "session", Revision: 19, Status: "playing"}})
	if next.(Model).player.Status != "paused" {
		t.Fatal("old poll reversed paused status")
	}
	next, _ = next.(Model).Update(playerMsg{state: ipc.State{Instance: "new", Revision: 1, Status: "playing"}})
	if next.(Model).player.Status != "playing" {
		t.Fatal("helper restart snapshot rejected")
	}
}
func TestSpaceDelegatesToggleInsteadOfUsingStaleUIState(t *testing.T) {
	m := fixture()
	m.player = ipc.State{URI: "spotify:track:fixture", Status: "paused"}
	m, cmd := key(m, tea.KeySpace)
	if cmd == nil {
		t.Fatal("no toggle")
	}
	_ = cmd()
	p := m.Player.(*playerStub)
	if len(p.calls) != 1 || p.calls[0] != "toggle_pause" {
		t.Fatal(p.calls)
	}
}
