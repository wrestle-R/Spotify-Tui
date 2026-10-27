# Spotify TUI

A simple music player for Arch Linux: playlists, liked songs, search, and standalone playback. Go handles the terminal UI; Rust handles the playback helper.

This directory currently contains a detailed plan, not a working application.

- [Full implementation plan](PLAN.md)
- [Acceptance checks and evidence requirements](docs/ACCEPTANCE.md)
- [Research sources and platform constraints](docs/SOURCES.md)
- [Terminal design reference](design/tui-concept.png)
- [Download the planning bundle](docs/DOWNLOADS.md)
- [Future GitHub release workflow](docs/RELEASE-PLAN.md)

Private GitHub repository: [wrestle-R/Spotify-Tui](https://github.com/wrestle-R/Spotify-Tui).

Download the plan with `bash scripts/download-plan.sh --destination "$PWD/Spotify-Tui-plan"` using an authenticated GitHub CLI account with repository access. The script verifies SHA-256 before extraction and refuses to overwrite an existing destination. See the download guide for obtaining the script without cloning.

**Offline playback is required, and a supported implementation route has not been established.** The first development gate resolves that requirement. An online-only player does not satisfy this project.

Prepared on 28 October 2026, Asia/Calcutta. No account login, software installation, application implementation, or playback testing was performed while preparing these documents.
