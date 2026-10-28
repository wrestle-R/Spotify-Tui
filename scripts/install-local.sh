#!/usr/bin/env bash
set -euo pipefail
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
if [[ -f "$script_dir/release.py" && -x "$script_dir/spotify-tui" ]]; then
  exec python3 "$script_dir/release.py" install --from "$script_dir" "$@"
else
  exec python3 "$script_dir/release.py" install --from "$script_dir/../dist" "$@"
fi
