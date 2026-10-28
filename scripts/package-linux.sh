#!/usr/bin/env bash
set -euo pipefail
project_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
cd -- "$project_dir"
version=${1:-0.2.1}
[[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]] || { printf 'Invalid version\n' >&2; exit 1; }
[[ "$(uname -s)" == Linux && "$(uname -m)" == x86_64 ]] || { printf 'Package on Linux x86_64\n' >&2; exit 1; }
[[ "$(dist/spotify-tui version)" == "spotify-tui $version" && "$(dist/spotify-tui-playback --version)" == "spotify-tui-playback $version" ]] || { printf 'Build both matching binaries first\n' >&2; exit 1; }
name="spotify-tui_${version}_linux_amd64"
mkdir -p dist/release
stage=$(mktemp -d -- dist/.package.XXXXXX)
trap 'rm -rf -- "$stage"' EXIT
mkdir -- "$stage/$name"
install -m 755 dist/spotify-tui dist/spotify-tui-playback "$stage/$name/"
install -m 755 scripts/install-local.sh "$stage/$name/install.sh"
cp README.md "$stage/$name/README.md"
cp -r docs design "$stage/$name/"
cp PLAN.md CHANGELOG.md "$stage/$name/"
printf 'Linux x86_64 release. Tested on Arch Linux.\nRust helper requires glibc >= 2.39 and libgcc_s; separate Soloist requires PipeWire or PulseAudio.\nDesktop keyring recommended. Spotify Premium and personal developer/provider setup required.\nOffline downloads unsupported. Source commit before packaging: %s\n' "$(git rev-parse HEAD)" > "$stage/$name/RUNTIME.txt"
tar -czf "dist/release/$name.tar.gz" -C "$stage" "$name"
cp scripts/download-linux.sh dist/release/download-linux.sh
(cd dist/release && sha256sum "$name.tar.gz" download-linux.sh > SHA256SUMS)
printf 'Packaged %s\n' "dist/release/$name.tar.gz"
