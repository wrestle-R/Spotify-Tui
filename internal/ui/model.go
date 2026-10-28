package ui

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/wrestle-R/Spotify-Tui/internal/auth"
	"github.com/wrestle-R/Spotify-Tui/internal/config"
	"github.com/wrestle-R/Spotify-Tui/internal/ipc"
	"github.com/wrestle-R/Spotify-Tui/internal/spotify"
)

type Library interface {
	Playlists(context.Context) ([]spotify.Playlist, error)
	Tracks(context.Context, string) (spotify.Collection, error)
	Search(context.Context, string) ([]spotify.Track, error)
}
type Player interface {
	Call(context.Context, string, any, any) error
	State(context.Context) (ipc.State, error)
	Ensure(context.Context, config.Config, config.Paths, auth.Vault) error
}
type rowsMsg struct {
	generation uint64
	tracks     []spotify.Track
	stale      bool
	warning    string
	err        error
}
type playlistsMsg struct {
	items []spotify.Playlist
	err   error
}
type playerMsg struct {
	state ipc.State
	err   error
}
type actionMsg struct {
	message string
	err     error
}
type loginMsg struct{ err error }
type tickMsg time.Time
type Model struct {
	Cfg                     config.Config
	Paths                   config.Paths
	Auth                    *auth.Manager
	Library                 Library
	Player                  Player
	width, height           int
	playlists               []spotify.Playlist
	tracks                  []spotify.Track
	selected, sideSelection int
	sidebar                 bool
	view, title, playlistID string
	generation              uint64
	loading                 bool
	input                   string
	inputMode               string
	status                  string
	player                  ipc.State
	help                    bool
	loginCancel             context.CancelFunc
	dataCancel              context.CancelFunc
}

var accent = lipgloss.Color("114")

func New(c config.Config, p config.Paths, a *auth.Manager, library Library, player Player) Model {
	m := Model{Cfg: c, Paths: p, Auth: a, Library: library, Player: player, width: 100, height: 30, title: "Liked songs", view: "liked", status: "L sign in · c connect player · ? help", sidebar: true}
	if c.ClientID == "" {
		m.inputMode = "setup"
		m.status = "Create a Spotify Developer app, then paste its public client ID below."
	}
	return m
}
func tick() tea.Cmd { return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }) }
func (m Model) Init() tea.Cmd {
	if m.Cfg.ClientID == "" {
		return tick()
	}
	return tea.Batch(tick(), m.loadPlaylists(), m.loadRows("liked", "", m.generation), m.pollPlayer())
}
func (m Model) loadPlaylists() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		items, e := m.Library.Playlists(ctx)
		return playlistsMsg{items, e}
	}
}
func (m Model) loadRows(view, value string, generation uint64) tea.Cmd {
	return m.loadRowsContext(context.Background(), view, value, generation)
}
func (m Model) loadRowsContext(parent context.Context, view, value string, generation uint64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(parent, 2*time.Minute)
		defer cancel()
		if view == "search" {
			items, e := m.Library.Search(ctx, value)
			return rowsMsg{generation: generation, tracks: items, err: e}
		}
		items, e := m.Library.Tracks(ctx, value)
		return rowsMsg{generation: generation, tracks: items.Tracks, stale: items.Stale, warning: items.Warning, err: e}
	}
}
func (m Model) pollPlayer() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		s, e := m.Player.State(ctx)
		return playerMsg{s, e}
	}
}
func (m Model) action(method string, params any) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()
		e := m.Player.Call(ctx, method, params, nil)
		return actionMsg{message: "Command accepted; waiting for player state", err: e}
	}
}
func (m *Model) open(view, title, value string) tea.Cmd {
	if m.dataCancel != nil {
		m.dataCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.dataCancel = cancel
	m.view = view
	m.title = title
	m.playlistID = value
	m.selected = 0
	m.tracks = nil
	m.generation++
	m.loading = true
	m.sidebar = false
	m.status = "Loading..."
	return m.loadRowsContext(ctx, view, value, m.generation)
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		return m, tea.Batch(tick(), m.pollPlayer())
	case playerMsg:
		if msg.err == nil {
			m.player = msg.state
		} else {
			m.player.Connected = false
			m.player.Status = "unavailable"
			m.player.Error = msg.err.Error()
		}
	case rowsMsg:
		if msg.generation != m.generation {
			return m, nil
		}
		m.loading = false
		if msg.err != nil {
			m.status = msg.err.Error()
		} else {
			m.tracks = msg.tracks
			m.status = fmt.Sprintf("%d tracks", len(msg.tracks))
			if msg.stale {
				m.status = "Cached metadata · " + msg.warning
			}
			m.selected = min(m.selected, max(0, len(m.tracks)-1))
		}
	case playlistsMsg:
		if msg.err == nil {
			m.playlists = msg.items
		} else if len(m.playlists) == 0 {
			m.status = msg.err.Error()
		}
	case actionMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
		} else {
			m.status = msg.message
		}
		return m, m.pollPlayer()
	case loginMsg:
		m.loginCancel = nil
		if msg.err != nil {
			m.status = msg.err.Error()
		} else {
			m.status = "Signed in. Web API login is separate from Soloist pairing."
			if resetter, ok := m.Library.(interface{ ResetAccount() }); ok {
				resetter.ResetAccount()
			}
			m.playlists = nil
			m.tracks = nil
			m.view = "liked"
			m.title = "Liked songs"
			m.playlistID = ""
			m.loading = true
			m.generation++
			return m, tea.Batch(m.loadPlaylists(), m.loadRows("liked", "", m.generation))
		}
	case tea.KeyMsg:
		if msg.Type == tea.KeyRunes && len(msg.Runes) > 1 && m.inputMode == "" {
			var commands []tea.Cmd
			for _, r := range msg.Runes {
				next, command := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
				m = next.(Model)
				if command != nil {
					commands = append(commands, command)
				}
			}
			return m, tea.Batch(commands...)
		}
		key := msg.String()
		if key == "ctrl+c" {
			if m.dataCancel != nil {
				m.dataCancel()
			}
			if m.loginCancel != nil {
				m.loginCancel()
			}
			return m, tea.Quit
		}
		if m.help {
			m.help = false
			return m, nil
		}
		if m.inputMode != "" {
			switch key {
			case "esc":
				if m.inputMode != "setup" {
					m.inputMode = ""
					m.input = ""
				}
			case "backspace", "ctrl+h":
				if len(m.input) > 0 {
					_, size := utf8.DecodeLastRuneInString(m.input)
					m.input = m.input[:len(m.input)-size]
				}
			case "enter":
				value := strings.TrimSpace(m.input)
				if m.inputMode == "setup" {
					if len(value) != 32 || strings.IndexFunc(value, func(r rune) bool { return !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') }) >= 0 {
						m.status = "Client ID must be the 32-character public ID from Spotify Developer settings"
						return m, nil
					}
					m.Cfg.ClientID = value
					m.Auth.Config = m.Cfg
					if e := config.Save(m.Paths, m.Cfg); e != nil {
						m.status = e.Error()
						return m, nil
					}
					m.inputMode = ""
					m.input = ""
					m.status = "Client ID saved. Press L for browser sign-in."
				} else if value != "" {
					m.inputMode = ""
					m.input = ""
					return m, m.open("search", "Search: "+value, value)
				}
			default:
				if msg.Type == tea.KeyRunes && len(m.input) < 512 {
					m.input += spotify.Clean(string(msg.Runes))
				} else if key == " " || key == "space" {
					m.input += " "
				}
			}
			return m, nil
		}
		switch key {
		case "q":
			if m.dataCancel != nil {
				m.dataCancel()
			}
			if m.loginCancel != nil {
				m.loginCancel()
			}
			return m, tea.Quit
		case "Q":
			if m.dataCancel != nil {
				m.dataCancel()
			}
			if m.loginCancel != nil {
				m.loginCancel()
			}
			return m, tea.Sequence(m.action("shutdown", nil), tea.Quit)
		case "?":
			m.help = true
		case "/":
			m.inputMode = "search"
			m.input = ""
		case "tab", "shift+tab":
			m.sidebar = !m.sidebar
		case "j", "down":
			if m.sidebar {
				m.sideSelection = min(m.sideSelection+1, len(m.playlists)+2)
			} else {
				m.selected = min(m.selected+1, max(0, len(m.tracks)-1))
			}
		case "k", "up":
			if m.sidebar {
				m.sideSelection = max(0, m.sideSelection-1)
			} else {
				m.selected = max(0, m.selected-1)
			}
		case "g", "home":
			if m.sidebar {
				m.sideSelection = 0
			} else {
				m.selected = 0
			}
		case "G", "end":
			if m.sidebar {
				m.sideSelection = len(m.playlists) + 2
			} else {
				m.selected = max(0, len(m.tracks)-1)
			}
		case "pgdown":
			m.selected = min(m.selected+max(1, m.height-12), max(0, len(m.tracks)-1))
		case "pgup":
			m.selected = max(0, m.selected-max(1, m.height-12))
		case "enter":
			if m.sidebar {
				switch m.sideSelection {
				case 0:
					return m, m.open("liked", "Liked songs", "")
				case 1:
					m.status = "Choose a playlist below and press Enter"
				case 2:
					m.view = "downloads"
					m.title = "Downloads"
					m.tracks = nil
					m.generation++
					m.loading = false
					m.sidebar = false
				default:
					i := m.sideSelection - 3
					if i >= 0 && i < len(m.playlists) {
						p := m.playlists[i]
						return m, m.open("playlist", p.Name, p.ID)
					}
				}
			} else if len(m.tracks) > m.selected {
				t := m.tracks[m.selected]
				if !t.CanPlay() {
					m.status = "This item is not playable through this provider"
				} else {
					return m, m.action("play_track", map[string]any{"uri": t.URI})
				}
			}
		case "P":
			if m.view == "playlist" {
				return m, m.action("play_collection", map[string]any{"uri": "spotify:playlist:" + m.playlistID})
			}
		case "a":
			if len(m.tracks) > m.selected && m.tracks[m.selected].CanPlay() {
				return m, m.action("add_to_queue", map[string]any{"uri": m.tracks[m.selected].URI})
			}
		case " ", "space":
			if m.player.Status == "playing" {
				return m, m.action("pause", nil)
			}
			return m, m.action("resume", nil)
		case "n":
			return m, m.action("next", nil)
		case "p":
			return m, m.action("previous", nil)
		case "h", "l":
			delta := int64(5000)
			if key == "h" {
				delta = -delta
			}
			return m, m.action("seek", map[string]any{"position_ms": max(int64(0), currentPosition(m.player)+delta)})
		case "+", "=", "-":
			delta := 5
			if key == "-" {
				delta = -5
			}
			return m, m.action("set_volume", map[string]any{"volume": max(0, min(100, m.player.Volume+delta))})
		case "s":
			return m, m.action("set_shuffle", map[string]any{"enabled": !m.player.Shuffle})
		case "r":
			mode := "context"
			if m.player.Repeat == "context" {
				mode = "track"
			} else if m.player.Repeat == "track" {
				mode = "off"
			}
			return m, tea.Sequence(m.action("repeat_track", map[string]any{"enabled": false}), m.action("repeat_context", map[string]any{"enabled": mode == "context"}), m.action("repeat_track", map[string]any{"enabled": mode == "track"}))
		case "d", "D":
			m.view = "downloads"
			m.title = "Downloads"
			m.generation++
			m.loading = false
			m.tracks = nil
			m.sidebar = false
		case "R":
			if m.view == "liked" {
				return m, m.open("liked", "Liked songs", "")
			} else if m.view == "playlist" {
				return m, m.open("playlist", m.title, m.playlistID)
			}
		case "L":
			if m.loginCancel != nil {
				m.status = "Browser sign-in already pending"
				return m, nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			m.loginCancel = cancel
			m.status = "Finish Spotify authorization in your browser. Ctrl-C cancels."
			a := m.Auth
			return m, func() tea.Msg { defer cancel(); return loginMsg{a.Login(ctx, auth.OpenBrowser, io.Discard)} }
		case "c":
			cfg, paths, vault := m.Cfg, m.Paths, m.Auth.Vault
			return m, func() tea.Msg {
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				err := m.Player.Ensure(ctx, cfg, paths, vault)
				return actionMsg{message: "Helper ready. Pair Spotify TUI Local once if prompted.", err: err}
			}
		}
	}
	return m, nil
}

func clip(s string, width int) string {
	if width < 1 {
		return ""
	}
	return ansi.Truncate(spotify.Clean(s), width, "…")
}
func duration(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	return fmt.Sprintf("%d:%02d", ms/60000, (ms/1000)%60)
}
func currentPosition(state ipc.State) int64 {
	position := state.PositionMS
	if state.Status == "playing" && state.PositionUpdatedAt > 0 {
		position += max(int64(0), time.Now().UnixMilli()-state.PositionUpdatedAt)
	}
	if state.DurationMS > 0 {
		position = min(position, state.DurationMS)
	}
	return max(int64(0), position)
}
func (m Model) View() string {
	w, h := m.width, m.height
	if w < 35 || h < 10 {
		return "Resize terminal to at least 35×10. q quits; Ctrl-C exits.\n"
	}
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	active := lipgloss.NewStyle().Foreground(accent).Bold(true)
	if m.help {
		return lipgloss.NewStyle().Width(w-4).Padding(1, 2).Render("music · keys\n\nTab focus · arrows/j/k move · Enter play/open · P play playlist from start\n/ search · Space pause · n/p next/previous · h/l seek · +/- volume\na queue selected · s shuffle · r repeat · R refresh\nL browser login · c connect Soloist · d downloads status\nq closes the TUI (audio continues) · Q stops the helper and exits\n\nSpotify downloads remain unsupported. Premium does not add a download API.\n\nAny key closes help.")
	}
	var lines []string
	lines = append(lines, active.Render("music")+muted.Render("  ·  Go + Rust  ·  "+m.player.Status))
	if m.inputMode == "setup" {
		return strings.Join(append(lines, "", "Spotify Developer public client ID:", clip(m.input+"_", w-2), "", clip(m.status, w-2), "", "Register redirect: "+m.Cfg.RedirectURI, "Press Enter to save. Ctrl-C exits."), "\n")
	}
	search := "/ search songs"
	if m.inputMode == "search" {
		search = "Search: " + m.input + "_"
	}
	lines = append(lines, clip(search, w), strings.Repeat("─", w))
	bodyH := max(1, h-10)
	sideW := min(24, w/3)
	if w < 70 {
		sideW = 0
	}
	labels := []string{"Liked songs", "Playlists", "Downloads"}
	for _, p := range m.playlists {
		labels = append(labels, p.Name)
	}
	start := 0
	if m.selected >= bodyH-1 {
		start = m.selected - bodyH + 2
	}
	for row := 0; row < bodyH; row++ {
		left := ""
		if sideW > 0 {
			sideStart := max(0, m.sideSelection-bodyH+1)
			i := row + sideStart
			if i < len(labels) {
				prefix := "  "
				if i == m.sideSelection {
					prefix = "> "
				}
				left = clip(prefix+labels[i], sideW-2)
				if i == m.sideSelection && m.sidebar {
					left = active.Render(left)
				}
			}
			left = lipgloss.NewStyle().Width(sideW).Render(left) + muted.Render("│ ")
		}
		rightW := w - sideW
		if sideW > 0 {
			rightW -= 2
		}
		right := ""
		if row == 0 {
			right = active.Render(clip(m.title, rightW))
		} else if m.view == "downloads" {
			messages := []string{"Offline Spotify downloads: unsupported", "A playback cache is not an offline library.", "No download will be reported as complete.", "See docs/FEASIBILITY.md for the blocker."}
			if row-1 < len(messages) {
				right = clip(messages[row-1], rightW)
			}
		} else if i := start + row - 1; i < len(m.tracks) {
			t := m.tracks[i]
			prefix := "  "
			if i == m.selected {
				prefix = "> "
			}
			nameW := max(5, (rightW-12)*3/5)
			artistW := max(0, rightW-nameW-12)
			name := lipgloss.NewStyle().Width(nameW).Render(clip(t.Name, nameW))
			artist := lipgloss.NewStyle().Width(artistW).Render(clip(t.Artist(), artistW))
			right = prefix + name + "  " + artist + " " + duration(t.DurationMS)
			right = clip(right, rightW)
			if i == m.selected && !m.sidebar {
				right = active.Render(right)
			}
		} else if row == 1 && m.loading {
			right = "Loading..."
		} else if row == 1 && len(m.tracks) == 0 {
			right = "No tracks loaded. L signs in; R refreshes."
		}
		lines = append(lines, left+right)
	}
	lines = append(lines, strings.Repeat("─", w))
	title := m.player.Title
	if title == "" {
		title = "Nothing playing"
	}
	lines = append(lines, clip("Playing: "+title+" · "+m.player.Artist, w))
	position := currentPosition(m.player)
	lines = append(lines, clip(fmt.Sprintf("%s / %s   volume %d%%   shuffle %t   repeat %s", duration(position), duration(m.player.DurationMS), m.player.Volume, m.player.Shuffle, m.player.Repeat), w))
	lines = append(lines, muted.Render(clip(m.status, w)))
	if m.player.Error != "" {
		lines = append(lines, muted.Render(clip(m.player.Error, w)))
	} else {
		lines = append(lines, "")
	}
	lines = append(lines, muted.Render(clip("Tab focus  / search  Enter play  Space pause  n next  L login  ? help  q detach", w)))
	return strings.Join(lines, "\n")
}
