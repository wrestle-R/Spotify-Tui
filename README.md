# spotify tui

So Spotify was taking up way too much RAM and I just wanted to play my playlists. I made a terminal player coz I don't need a whole desktop app sitting there while I'm working.

Go for the TUI, Rust for the local playback helper. **Linux only for now.** The download is for x86_64, tested on Arch Linux.

![Spotify TUI running in a terminal on Arch Linux](design/tui-screenshot.png)

Actual terminal capture of `--demo`, using generic sample playlists and songs. The demo produces no audio and reads no account data.

[What's new](CHANGELOG.md) · [Full release details](docs/releases/v0.2.1.md)

## Download

**[Download for Linux x86_64](https://github.com/wrestle-R/Spotify-Tui/releases/download/v0.2.1/spotify-tui_0.2.1_linux_amd64.tar.gz)** · [Release notes and checksums](https://github.com/wrestle-R/Spotify-Tui/releases/tag/v0.2.1)

This is the first packaged release. The archive includes the TUI, Rust helper and local installer. It needs glibc 2.39 or newer, a Spotify Premium account and separate [Soloist setup](docs/SETUP.md#local-playback-through-official-soloist). Soloist comes from Spotify and isn't bundled here.

The repo is private, so downloads require GitHub access. Using an authenticated GitHub CLI:

```bash
gh release download v0.2.1 --repo wrestle-R/Spotify-Tui --pattern download-linux.sh
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
