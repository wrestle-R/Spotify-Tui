package ui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wrestle-R/Spotify-Tui/go/internal/update"
	"time"
)

type Updater interface {
	Check(context.Context) (*update.Release, error)
	Install(context.Context, update.Release) (update.Installation, error)
}
type updateCheckedMsg struct {
	release *update.Release
	err     error
}
type updateInstalledMsg struct {
	installation update.Installation
	err          error
}

func (m Model) WithUpdater(updater Updater, version string) Model {
	m.updater = updater
	m.version = version
	return m
}
func (m Model) RestartPath() string { return m.restartPath }
func (m Model) StopRequested() bool { return m.stopRequested }
func (m Model) checkUpdates() tea.Cmd {
	if m.demo || m.updater == nil {
		return nil
	}
	client := m.updater
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		r, err := client.Check(ctx)
		return updateCheckedMsg{r, err}
	}
}
func (m *Model) requestUpdate() tea.Cmd {
	if m.demo {
		m.notify("Updates are disabled in the demo.", "info")
		return nil
	}
	if m.updater == nil {
		m.notify("Update service unavailable.", "error")
		return nil
	}
	if m.updating {
		m.notify("Update in progress. You can keep browsing.", "info")
		return nil
	}
	if m.busy || m.connecting {
		m.notify("Finish the player action, then press Ctrl-U to update.", "info")
		return nil
	}
	if m.latest == nil {
		m.installAfterCheck = true
		m.notify("Checking for a stable update…", "info")
		return m.checkUpdates()
	}
	m.updating = true
	m.installAfterCheck = false
	m.updateNotice = "Installing " + m.latest.Tag + "…"
	m.notify("Downloading and verifying update. The app restarts when ready.", "info")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	m.updateCancel = cancel
	client, r := m.updater, *m.latest
	return func() tea.Msg {
		defer cancel()
		installation, err := client.Install(ctx, r)
		return updateInstalledMsg{installation, err}
	}
}

func (m Model) WithResume(uri string, position int64) Model {
	m.resumeURI = uri
	m.resumePosition = max(int64(0), position)
	return m
}
func (m *Model) resumeAfterUpdate() tea.Cmd {
	if m.resumeURI == "" {
		return nil
	}
	return m.perform("Restore playback", playerCommand{"play_track", map[string]any{"uri": m.resumeURI}}, playerCommand{"seek", map[string]any{"position_ms": m.resumePosition}})
}
