# Spotify TUI

A simple music player for Arch Linux: playlists, liked songs, search, and standalone playback. Go handles the terminal UI; Rust handles the playback helper.

The initial Go/Rust development build is implemented. Real browser sign-in, library/search, the keyboard UI, local audio with the official clients closed, and a clean provider restart have passed on the development machine. The network-isolated cold start failed. Offline downloads remain blocked by the provider capability limitation.

```bash
make build
./dist/spotify-tui
```

See [setup and normal-browser login](docs/SETUP.md) for the client ID, callback, and separate Soloist pairing. `make test` and `make check` cover the automated checks. [Implementation status](docs/IMPLEMENTATION.md) separates verified behavior from remaining work.

- [Full implementation plan](PLAN.md)
- [Acceptance checks and evidence requirements](docs/ACCEPTANCE.md)
- [Research sources and platform constraints](docs/SOURCES.md)
- [Current feasibility decision and verified results](docs/FEASIBILITY.md)
- [Terminal design reference](design/tui-concept.png)
- [Download the planning bundle](docs/DOWNLOADS.md)
- [Download development source and build locally](docs/DOWNLOADS.md#development-source)
- [Future GitHub release workflow](docs/RELEASE-PLAN.md)

Private GitHub repository: [wrestle-R/Spotify-Tui](https://github.com/wrestle-R/Spotify-Tui).

Download the plan with `bash scripts/download-plan.sh --destination "$PWD/Spotify-Tui-plan"` using an authenticated GitHub CLI account with repository access. The script verifies SHA-256 before extraction and refuses to overwrite an existing destination. See the download guide for obtaining the script without cloning.

**Offline playback is required, and a supported implementation route has not been established.** The first development gate resolves that requirement. An online-only player does not satisfy this project.

Prepared on 28 October 2026, Asia/Calcutta. The original planning archive remains documentation-only. Main now also contains the initial implementation. No working offline-player claim is made.
