#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
if [[ -x "$script_dir/spotify-tui" && -x "$script_dir/spotify-tui-playback" ]]; then
  from_dir="$script_dir"
else
  project_dir=$(cd -- "$script_dir/.." && pwd)
  from_dir="$project_dir/dist"
fi
prefix="${HOME}/.local"
dry_run=false

fail() { printf 'Error: %s\n' "$*" >&2; exit 1; }
while [[ $# -gt 0 ]]; do
  case "$1" in
    --prefix) [[ $# -ge 2 && -n "$2" ]] || fail '--prefix needs a directory'; prefix=$2; shift 2 ;;
    --from) [[ $# -ge 2 && -n "$2" ]] || fail '--from needs a directory'; from_dir=$2; shift 2 ;;
    --dry-run) dry_run=true; shift ;;
    --help|-h)
      printf 'Install the Go/Rust release locally. No sudo or Spotify login.\nUsage: bash scripts/install-local.sh [--prefix PATH] [--from DIST_PATH] [--dry-run]\nOffline Spotify downloads remain unsupported.\n'
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
command -v flock >/dev/null || fail 'flock is required'
exec 9>"$prefix/lib/spotify-tui/.update.lock"
flock -n 9 || fail 'Another installation or update is already running'
current="$prefix/lib/spotify-tui/current"
if [[ -e "$current" || -L "$current" ]]; then
 [[ -L "$current" ]] || fail 'Current build pointer must be a symlink'
 case "$(readlink -f -- "$current")" in "$prefix"/lib/spotify-tui/build.*) ;; *) fail 'Invalid current build pointer';; esac
elif [[ -L "$prefix/bin/spotify-tui" ]]; then
 old_go=$(readlink -f -- "$prefix/bin/spotify-tui")
 old_rust=$(readlink -f -- "$prefix/bin/spotify-tui-playback")
 [[ "$(dirname -- "$old_go")" == "$(dirname -- "$old_rust")" ]] || fail 'Existing binary pair is inconsistent'
 ln -s -- "$(dirname -- "$old_go")" "$current"
fi
stage=$(mktemp -d -- "$prefix/lib/spotify-tui/build.XXXXXX")
keep=false
cleanup() {
  rm -f -- "$prefix/bin/.spotify-tui-link.$$" "$prefix/bin/.spotify-tui-playback-link.$$" "$prefix/lib/spotify-tui/.current.$$"
  if [[ "$keep" != true ]]; then rm -rf -- "$stage"; fi
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
install -m 0755 -- "$from_dir/spotify-tui" "$stage/spotify-tui"
install -m 0755 -- "$from_dir/spotify-tui-playback" "$stage/spotify-tui-playback"
for binary in spotify-tui spotify-tui-playback; do
  ln -s -- "$current/$binary" "$prefix/bin/.$binary-link.$$"
done

# Bootstrap current before changing command links. Upgrades switch one pointer.
if [[ ! -L "$current" ]]; then ln -s -- "$stage" "$current"; keep=true; fi
for binary in spotify-tui spotify-tui-playback; do
  mv -Tf -- "$prefix/bin/.$binary-link.$$" "$prefix/bin/$binary"
done
ln -s -- "$stage" "$prefix/lib/spotify-tui/.current.$$"
mv -Tf -- "$prefix/lib/spotify-tui/.current.$$" "$current"
keep=true
printf 'Installed release in %s\nRun: %s/bin/spotify-tui\n' "$stage" "$prefix"
printf 'Existing build directories are retained for rollback. Offline downloads are unsupported.\n'
