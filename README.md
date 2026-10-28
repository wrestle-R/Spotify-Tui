# spotify tui

So Spotify was taking up way too much RAM and I just wanted to play my playlists. I made a terminal player coz I don't need a whole desktop app sitting there while I'm working.

Go for the TUI, Rust for the local playback helper. **Linux only for now.** The download is for x86_64, tested on Arch Linux.

![Spotify TUI running in a terminal on Arch Linux](design/tui-screenshot.png)

Actual terminal capture of `--demo`, using generic sample playlists and songs. The demo produces no audio and reads no account data.

## Download

**[Download for Linux x86_64](https://github.com/wrestle-R/Spotify-Tui/releases/download/v0.2.0-alpha.1/spotify-tui_0.2.0-alpha.1_linux_amd64.tar.gz)** · [Release notes and checksums](https://github.com/wrestle-R/Spotify-Tui/releases/tag/v0.2.0-alpha.1)

This is an alpha. The archive includes the TUI, Rust helper and local installer. It needs glibc 2.39 or newer, a Spotify Premium account and separate [Soloist setup](docs/SETUP.md#local-playback-through-official-soloist). Soloist comes from Spotify and isn't bundled here.

The repo is private, so downloads require GitHub access. Using an authenticated GitHub CLI:

```bash
gh release download v0.2.0-alpha.1 --repo wrestle-R/Spotify-Tui --pattern download-linux.sh
less download-linux.sh
bash download-linux.sh --destination ./spotify-tui-linux
bash ./spotify-tui-linux/install.sh
~/.local/bin/spotify-tui
```

The downloader verifies SHA-256 before extraction and prints where it saved the files. Installation goes under `~/.local`; it doesn't need sudo. It keeps previous build directories for rollback. See [downloads](docs/DOWNLOADS.md) for manual verification and a different install prefix.

Try the keyboard flow without an account:

```bash
~/.local/bin/spotify-tui --demo
```

## First run

Use your own Spotify Developer application's public client ID, register `http://127.0.0.1:43821/callback`, then:

```bash
~/.local/bin/spotify-tui setup --client-id YOUR_PUBLIC_CLIENT_ID --soloist-path /absolute/path/to/soloist
~/.local/bin/spotify-tui login
~/.local/bin/spotify-tui provider-key
~/.local/bin/spotify-tui
```

Sign-in happens in your normal browser. `provider-key` hides the key while you type. Select **Spotify TUI Local** in Spotify's device picker for the initial pairing; you can close the official player afterward. Full [setup instructions](docs/SETUP.md) cover verification, the keyring and pairing.

## Moving around

The highlighted row is your selection. The dot beside **Library** or the collection title shows which pane has keyboard focus.

| Key | What it does |
| --- | --- |
| ↑ / ↓ or `j` / `k` | Move within the focused pane |
| → / Enter | Open the selected collection; Enter on a song plays it |
| ← / Esc | Go back to the library |
| Tab / Shift-Tab | Switch panes |
| Page Up / Down, Home / End | Page or jump within the focused pane |
| `/` or Ctrl-F | Search songs and artists; Enter searches, Esc cancels |
| Space | Pause or resume |
| `n` / `p` | Next or previous song |
| `[` / `]` | Seek backward or forward five seconds |
| `+` / `-` | Change volume |
| `s` / `r` | Shuffle / cycle repeat |
| `a` / `P` | Queue the selected song / play the open playlist from its start |
| `R` / `c` / `L` | Refresh the focused pane / reconnect / browser sign-in |
| `?` | Keyboard help |
| `q` / `Q` | Quit and keep audio playing / stop audio and quit |

Enter starts the local player automatically when needed. You can keep browsing while a request is running. On a narrow terminal, the focused pane takes the full width. Left/Esc brings the library back.

## Where it's at

Playlists, liked songs, search, local audio and playback controls have been tested with a real account on Arch. The keyboard flow is covered by regression checks and a real terminal walkthrough, including resize, search editing, automatic startup and quitting cleanly. [Test details](docs/IMPLEMENTATION.md).

**Offline downloads still don't work.** The current provider has no supported download operation, and its isolated offline cold-start test failed. Offline remains a requirement for the full app.

Soloist can occasionally stall while restoring its session; `c` restarts it. A selected liked song plays individually; a full liked-song sequence isn't wired up yet. `P` plays a playlist context. Search currently shows the first ten results. Spotify can restrict listing songs in followed or curated playlists. Long-session RAM and stability testing are still pending, so the RAM motivation above isn't a measured savings claim.

## Build from source

```bash
make build
./dist/spotify-tui
```

Requires Go and Rust/Cargo. `make test` runs the Go race checks and Rust tests; `make check` runs formatting, vet and clippy. The original [plan](PLAN.md) and [acceptance checks](docs/ACCEPTANCE.md) track the remaining work.

## Code layout

```text
go/        Go module: terminal UI, account sign-in and Spotify library
rust/      Rust crate: local playback helper and provider control
protocol/  Shared Go/Rust message fixtures
scripts/   Build packaging, downloads and local installation
docs/      Setup, verification and the remaining work
design/    Screenshot and original terminal concept
```

Run `make build`, `make test` and `make check` from the repository root. `make package-check` builds the Linux archive and checks its downloader and installer.
