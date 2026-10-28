#!/usr/bin/env bash
set -euo pipefail

project_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
from_dir="$project_dir/dist"
prefix="${HOME}/.local"
dry_run=false

fail() { printf 'Error: %s\n' "$*" >&2; exit 1; }
while [[ $# -gt 0 ]]; do
  case "$1" in
    --prefix) [[ $# -ge 2 && -n "$2" ]] || fail '--prefix needs a directory'; prefix=$2; shift 2 ;;
    --from) [[ $# -ge 2 && -n "$2" ]] || fail '--from needs a directory'; from_dir=$2; shift 2 ;;
    --dry-run) dry_run=true; shift ;;
    --help|-h)
      printf 'Install the development Go/Rust binaries locally. No sudo or Spotify login.\nUsage: bash scripts/install-local.sh [--prefix PATH] [--from DIST_PATH] [--dry-run]\nOffline Spotify downloads remain unsupported.\n'
      exit 0 ;;
    *) fail "Unknown option: $1" ;;
  esac
done

for binary in spotify-tui spotify-tui-playback; do
  [[ -f "$from_dir/$binary" && -x "$from_dir/$binary" ]] || fail "Missing executable $from_dir/$binary; run make build first"
done
if [[ "$dry_run" == true ]]; then
  printf 'Would install both binaries from %s to a version directory under %s/lib/spotify-tui and link them in %s/bin.\n' "$from_dir" "$prefix" "$prefix"
  exit 0
fi

mkdir -p -- "$prefix"
prefix=$(cd -- "$prefix" && pwd -P)
for binary in spotify-tui spotify-tui-playback; do
  target="$prefix/bin/$binary"
  if [[ -e "$target" || -L "$target" ]]; then
    [[ -L "$target" ]] || fail "Refusing to replace an existing unrelated file: $target"
    case "$(readlink -- "$target")" in
      "$prefix"/lib/spotify-tui/*/"$binary") ;;
      *) fail "Refusing to replace an unrelated symlink: $target" ;;
    esac
  fi
done

mkdir -p -- "$prefix/lib/spotify-tui" "$prefix/bin"
stage=$(mktemp -d -- "$prefix/lib/spotify-tui/build.XXXXXX")
keep=false
cleanup() {
  rm -f -- "$prefix/bin/.spotify-tui-link.$$" "$prefix/bin/.spotify-tui-playback-link.$$"
  if [[ "$keep" != true ]]; then rm -rf -- "$stage"; fi
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
install -m 0755 -- "$from_dir/spotify-tui" "$stage/spotify-tui"
install -m 0755 -- "$from_dir/spotify-tui-playback" "$stage/spotify-tui-playback"
for binary in spotify-tui spotify-tui-playback; do
  ln -s -- "$stage/$binary" "$prefix/bin/.$binary-link.$$"
done

# Keep a complete pair if interruption occurs during activation.
keep=true
for binary in spotify-tui spotify-tui-playback; do
  mv -Tf -- "$prefix/bin/.$binary-link.$$" "$prefix/bin/$binary"
done
printf 'Installed development build in %s\nRun: %s/bin/spotify-tui\n' "$stage" "$prefix"
printf 'Existing build directories are retained for rollback. Offline downloads are unsupported.\n'
