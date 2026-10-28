package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wrestle-R/Spotify-Tui/go/internal/auth"
	"github.com/wrestle-R/Spotify-Tui/go/internal/config"
	"github.com/wrestle-R/Spotify-Tui/go/internal/ipc"
	"github.com/wrestle-R/Spotify-Tui/go/internal/spotify"
	"github.com/wrestle-R/Spotify-Tui/go/internal/ui"
	"golang.org/x/term"
)

var version = "0.2.0"

const helpText = `Spotify TUI (Go interface + Rust helper)

Commands:
  --demo                   Try the UI with generic sample songs and no account
  setup --client-id ID [--soloist-path PATH] [--helper-path PATH]
        [--allow-file-token]  Explicit owner-only plaintext credential fallback
  login [--no-browser]      Browser PKCE authorization for library/search
  provider-key             Save a personal Soloist API key securely
  connect                  Start/connect the local playback helper
  doctor                   Show installation and provider status
  check-account            Check real library/search access
  check-playback           Briefly test local audio controls, then pause
  check-offline            Probe isolated offline startup, then restore online
  stop                     Stop this project's playback helper
  logout                   Stop helper and remove this app's credentials/cache
  version                  Print version

Run without a command to open the TUI. Offline Spotify downloads are unsupported.
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
func run() error {
	if len(os.Args) > 1 && (os.Args[1] == "--demo" || os.Args[1] == "demo") {
		_, err := tea.NewProgram(ui.NewDemo(), tea.WithAltScreen()).Run()
		return err
	}
	p, err := config.PathsForUser()
	if err != nil {
		return err
	}
	if err = p.Ensure(); err != nil {
		return err
	}
	cfg, err := config.Load(p)
	if err != nil {
		return err
	}
	manager := auth.New(cfg, p)
	player := ipc.New(p)
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "--version", "version":
			fmt.Println("spotify-tui", version)
			return nil
		case "--help", "help":
			fmt.Print(helpText)
			return nil
		case "setup":
			fs := flag.NewFlagSet("setup", flag.ContinueOnError)
			fs.StringVar(&cfg.ClientID, "client-id", cfg.ClientID, "Public Spotify Developer client ID")
			fs.StringVar(&cfg.SoloistPath, "soloist-path", cfg.SoloistPath, "Official Soloist executable")
			fs.StringVar(&cfg.HelperPath, "helper-path", cfg.HelperPath, "Rust helper executable")
			fs.BoolVar(&cfg.AllowFileToken, "allow-file-token", cfg.AllowFileToken, "Allow an owner-only plaintext credentials file if keyring unavailable")
			if e := fs.Parse(args[1:]); e != nil {
				return e
			}
			if len(cfg.ClientID) != 32 || strings.IndexFunc(cfg.ClientID, func(r rune) bool { return !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') }) >= 0 {
				return errors.New("client ID must be 32 hexadecimal characters")
			}
			if cfg.AllowFileToken {
				fmt.Println("Credential fallback enabled: owner-only files are plaintext, not encrypted.")
			}
			if err = config.Save(p, cfg); err != nil {
				return err
			}
			fmt.Println("Saved config. Register redirect:", cfg.RedirectURI)
			fmt.Println("Next: spotify-tui login")
			return nil
		case "login":
			fs := flag.NewFlagSet("login", flag.ContinueOnError)
			noBrowser := fs.Bool("no-browser", false, "Print an authorization URL for your normal browser")
			if err := fs.Parse(args[1:]); err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			open := auth.OpenBrowser
			if *noBrowser {
				open = func(link string) error { fmt.Fprintln(os.Stdout, "Open in your normal browser:\n"+link); return nil }
			}
			if err = manager.Login(ctx, open, os.Stdout); err != nil {
				return err
			}
			fmt.Println("Signed in for library/search. Soloist pairing is separate.")
			return nil
		case "provider-key":
			fmt.Print("Personal Soloist API key (hidden): ")
			b, e := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if e != nil {
				return e
			}
			key := strings.TrimSpace(string(b))
			if key == "" || len(key) > 4096 {
				return errors.New("invalid empty/oversized key")
			}
			if e = manager.Vault.Put("soloist", key); e != nil {
				return e
			}
			fmt.Println("Key stored. In the TUI press c to start Soloist, then pair Spotify TUI Local once.")
			return nil
		case "stop":
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return player.Stop(ctx)
		case "connect":
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			if err := player.Ensure(ctx, cfg, p, manager.Vault); err != nil {
				return err
			}
			fmt.Println("Playback helper started. Use doctor to check pairing; select Spotify TUI Local in Spotify if unpaired.")
			return nil
		case "logout":
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if e := player.Stop(ctx); e != nil {
				return errors.New("helper has not stopped; credentials were left intact, retry logout after stopping it")
			}
			if e := manager.Vault.Delete("oauth"); e != nil {
				return e
			}
			if e := manager.Vault.Delete("soloist"); e != nil {
				return e
			}
			if e := os.RemoveAll(filepath.Join(p.Cache, "metadata")); e != nil {
				return e
			}
			if e := os.RemoveAll(filepath.Join(p.Data, "soloist")); e != nil {
				return e
			}
			fmt.Println("This app's credentials and cached library removed.")
			return nil
		case "doctor":
			fmt.Println("spotify-tui", version)
			fmt.Println("Config:", filepath.Join(p.Config, "config.json"))
			fmt.Println("Client ID configured:", cfg.ClientID != "")
			_, e := exec.LookPath(cfg.SoloistPath)
			fmt.Println("Soloist executable available:", e == nil)
			fmt.Println("Socket:", player.Socket)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			s, e := player.State(ctx)
			fmt.Println("Helper reachable:", e == nil)
			if e == nil {
				fmt.Println("Provider connected:", s.Connected, "paired:", s.LoggedIn)
			}
			fmt.Println("Offline Spotify downloads: UNSUPPORTED (mandatory release blocker)")
			return nil
		case "check-account":
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			api := spotify.New(manager, p)
			if _, err := api.Account(ctx); err != nil {
				return err
			}
			fmt.Println("Account identity: OK")
			lists, err := api.Playlists(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("Playlist listing: OK (%d playlists)\n", len(lists))
			for _, list := range lists {
				if !api.CanReadPlaylistItems(list) {
					continue
				}
				items, e := api.Tracks(ctx, list.ID)
				if e != nil {
					fmt.Println("Owned/collaborative playlist contents: unavailable")
				} else {
					fmt.Printf("Owned/collaborative playlist contents: OK (%d tracks, cached=%t)\n", len(items.Tracks), items.Stale)
				}
				break
			}
			liked, err := api.Tracks(ctx, "")
			if err != nil {
				return err
			}
			fmt.Printf("Liked songs: OK (%d tracks, cached=%t)\n", len(liked.Tracks), liked.Stale)
			results, err := api.Search(ctx, "Radiohead")
			if err != nil {
				return err
			}
			fmt.Printf("Catalog search: OK (%d results)\n", len(results))
			fmt.Println("This check verifies metadata access, not audio or offline playback.")
			return nil
		case "check-playback":
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			api := spotify.New(manager, p)
			liked, err := api.Tracks(ctx, "")
			if err != nil {
				return err
			}
			for _, track := range liked.Tracks {
				if track.CanPlay() {
					return player.CheckPlayback(ctx, track.URI, os.Stdout)
				}
			}
			return errors.New("no playable liked track available for the audio check")
		case "check-offline":
			ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
			defer cancel()
			api := spotify.New(manager, p)
			liked, err := api.Tracks(ctx, "")
			if err != nil {
				return err
			}
			for _, track := range liked.Tracks {
				if track.CanPlay() {
					return player.OfflineProbe(ctx, cfg, p, manager.Vault, track.URI, os.Stdout)
				}
			}
			return errors.New("no eligible track to probe")
		default:
			return fmt.Errorf("unknown command %q; use --help", args[0])
		}
	}
	model := ui.New(cfg, p, manager, spotify.New(manager, p), player)
	_, err = tea.NewProgram(model, tea.WithAltScreen()).Run()
	return err
}
