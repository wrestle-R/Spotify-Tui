# spotify tui

So Spotify was taking up way too much RAM and I just wanted to play my playlists. I made a terminal player coz I don't need a whole desktop app sitting there while I'm working.

Go for the TUI, Rust for the local playback helper. **Linux only for now.** The download is for x86_64, tested on Arch Linux.

![Spotify TUI running in a terminal on Arch Linux](design/tui-screenshot.png)

Actual terminal capture of `--demo`, using generic sample playlists and songs. The demo produces no audio and reads no account data.

[What's new](CHANGELOG.md) · [Full release details](docs/releases/v0.3.0.md)

## Download

**[Download for Linux x86_64](https://github.com/wrestle-R/Spotify-Tui/releases/download/v0.3.0/spotify-tui_0.3.0_linux_amd64.tar.gz)** · [Release notes and checksums](https://github.com/wrestle-R/Spotify-Tui/releases/tag/v0.3.0)

 The archive includes the TUI, Rust helper and local installer. It needs glibc 2.39 or newer, a Spotify Premium account and separate [Soloist setup](docs/SETUP.md#local-playback-through-official-soloist). Soloist comes from Spotify and isn't bundled here.

The repo and downloads are public:

```bash
curl -fL https://github.com/wrestle-R/Spotify-Tui/releases/download/v0.3.0/download-linux.sh -o download-linux.sh
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

The Linux tray starts with the player when your panel supports it. `q` leaves music playing; `Q` stops music and exits the tray. The app checks stable updates on startup; **Ctrl-U** installs a verified update and restarts. [Tray and update help](docs/SETUP.md#tray-and-updates).
