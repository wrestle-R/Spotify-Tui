# Downloads and release plan

## Linux alpha download

[Linux x86_64 archive](https://github.com/wrestle-R/Spotify-Tui/releases/download/v0.2.0-alpha.1/spotify-tui_0.2.0-alpha.1_linux_amd64.tar.gz) · [release files](https://github.com/wrestle-R/Spotify-Tui/releases/tag/v0.2.0-alpha.1). This is a development alpha, not the completed offline player. Tested on Arch; glibc 2.39 or newer and libgcc_s are required. Soloist is obtained separately from Spotify.

The private repository requires authenticated access. Download the script, inspect it and run it:

```bash
gh release download v0.2.0-alpha.1 --repo wrestle-R/Spotify-Tui --pattern download-linux.sh
less download-linux.sh
bash download-linux.sh --destination ./spotify-tui-linux
bash ./spotify-tui-linux/install.sh --dry-run
bash ./spotify-tui-linux/install.sh
~/.local/bin/spotify-tui --demo
```

The downloader validates the release archive checksum, paths and entry types before extraction; it never executes the downloaded app. The installer activates both binaries under a version directory and refuses unrelated files. To change the destination, pass `--prefix /your/prefix` to `install.sh`. Previous build directories are retained. Spotify credentials and provider data stay outside the install directory.

For a manual download, retrieve the archive and `SHA256SUMS` from that release, verify with `sha256sum --check --strict --ignore-missing SHA256SUMS`, then extract into an empty directory. SHA-256 from the same release detects corruption; it is not an independent signature. The archive contains both binaries, `install.sh`, documentation and the generic demo screenshot. No account data or provider binary is bundled.

## Development source

Main contains the initial Go/Rust implementation, with the limitations in [IMPLEMENTATION.md](IMPLEMENTATION.md). This is not an accepted offline-player release. To retrieve source without installing or executing it, inspect `scripts/download-source.sh` and run:

```bash
bash scripts/download-source.sh --destination "$PWD/Spotify-Tui-source"
```

The script uses authenticated GitHub CLI access, verifies Git objects, exports only the requested commit's tracked files, and prints the resolved commit. Use `--ref FULL_COMMIT_SHA` for a specific version. It refuses existing destinations and excludes Git metadata, private notes, credentials, generated binaries and the vendor provider. It requires `gh`, `git`, and `tar`. The default `main` can change; Git object integrity is not an independent publisher signature.

If you do not have the source yet, obtain it with `gh repo clone wrestle-R/Spotify-Tui`, inspect the scripts, then follow [SETUP.md](SETUP.md). The local installer in `scripts/install-local.sh` operates on binaries you build yourself; it does not download a released player or Soloist.

## Available now: planning documents

Repository: [wrestle-R/Spotify-Tui](https://github.com/wrestle-R/Spotify-Tui), private by default. Only accounts with access can download its files. An anonymous public link would require a separate visibility decision; this planning task does not make the repository public.

Planning release: [planning-v0.1.0](https://github.com/wrestle-R/Spotify-Tui/releases/tag/planning-v0.1.0).

Release assets:

- [Planning archive](https://github.com/wrestle-R/Spotify-Tui/releases/download/planning-v0.1.0/spotify-tui-plan-v0.1.0.tar.gz)
- [Download script](https://github.com/wrestle-R/Spotify-Tui/releases/download/planning-v0.1.0/download-plan.sh)
- [SHA-256 checksum](https://github.com/wrestle-R/Spotify-Tui/releases/download/planning-v0.1.0/SHA256SUMS)

For private repositories use authenticated `gh release download`; direct asset URLs can require an authenticated browser and may not work with plain anonymous curl. The verification below confirms download integrity, not the functionality of an application.

From the project directory, inspect and run the included script:

```bash
bash scripts/download-plan.sh --help
bash scripts/download-plan.sh --destination "$PWD/Spotify-Tui-plan"
```

To get the script without cloning, first sign into GitHub CLI with an account that can access the repository. Run the download in an empty directory so an existing script is not overwritten:

```bash
gh release download planning-v0.1.0 --repo wrestle-R/Spotify-Tui --pattern download-plan.sh
less download-plan.sh
bash download-plan.sh --destination "$PWD/Spotify-Tui-plan"
```

The script downloads the archive and checksum, verifies SHA-256, checks archive paths/types, and extracts to a new directory. It refuses to overwrite an existing destination. It does not use sudo, install packages, authenticate Spotify, or change audio settings. Requires the GNU tools normally available on Arch Linux and GitHub CLI.

The archive contains README, PLAN, source ledger, acceptance specification, download guide, download script, GitHub release workflow plan, and the terminal mockup. It excludes `.git`, private `local-docs`, tokens, and environment files.

Checksums from the same release detect corruption and mismatches. They are not independent proof against compromise of the release account. Do not describe SHA-256 alone as a signature.

## Future: working player release

The user also wants a simple app download/install route. Implement it only after PLAN.md's offline and playback gates pass. No app binary or app installer is supplied by this planning release.

Proposed release naming:

- Tag `v0.1.0` for the first actually accepted player, distinct from `planning-v0.1.0`.
- Asset `spotify-tui_v0.1.0_linux_amd64.tar.gz` for the initial x86_64 Arch target.
- Archive includes matching `spotify-tui` (Go) and `spotify-tui-playback` (Rust) executables.
- `SHA256SUMS` plus signed checksums or build attestations when configured.
- `install.sh` with documented flags and an uninstall path.

The future installer must detect Linux/CPU architecture, select an explicit version, download to a temporary directory, verify integrity before extraction, validate archive paths, install both binaries atomically to a user prefix such as `~/.local/bin`, and keep rollback information. It must support `--version`, `--prefix`, `--dry-run`, and `--help`; document authentication for private assets. It must not execute downloaded content before verification, request Spotify credentials, mutate shell configuration silently, or use sudo.

Do not provide a fake `curl | sh` command before an actual release and installer exist. The initial documented flow should download the script, allow inspection, and run it explicitly. If a future public single-command bootstrap is offered, it must point at a real versioned installer and explain its trust boundary.

If Soloist is the provider, obtain it separately from its official distribution unless Spotify changes the redistribution permission. A user-generated API key remains local; it is never embedded in our installer or release.

## Future GitHub workflow

See [RELEASE-PLAN.md](RELEASE-PLAN.md) for the proposed build/release checks. This is a workflow specification, not an installed GitHub Actions workflow.
