# RAM, measured on one Linux laptop

This is a short, paused-library comparison, not a promise for every machine or a
long-session stability test. Both apps were signed into the same account, open at
the same time, with their libraries loaded and playback paused. The user confirmed
the GUI was signed in; its MPRIS playback state and the TUI helper state were checked
before and after sampling. Neither player was a demo or login screen.

## Result

**168.8 MiB for the whole TUI setup versus 814.1 MiB for the official GUI**, a
**79.3% lower median resident PSS** in this particular test. TUI samples ranged
from 162.7–169.7 MiB; GUI samples from 806.9–863.9 MiB. GUI swap PSS median was
3.4 MiB; every TUI group's swap PSS was zero across this run.

| Counted group | Median PSS (MiB) | Processes |
| --- | ---: | ---: |
| Go frontend and Linux tray | 31.2 | 2 |
| Rust helper | 1.8 | 1 |
| Soloist and crash handler | 46.4 | 2 |
| Dedicated Kitty terminal and descendants | 89.4 | 3 |
| Whole TUI setup | 168.8 | 8 |
| Official Spotify GUI and descendants | 814.1 | 9 |

Component medians need not sum exactly to the median of each sample's total.
The terminal accounts for more than half the TUI stack's resident memory here.

## What is counted

The TUI total includes **both Go processes (interface and tray), the Rust helper,
both official Soloist processes, and its dedicated Kitty terminal**. The official
GUI total includes its main process and all discovered descendants, including
renderer, GPU and utility processes. Processes are deduplicated. The desktop panel,
compositor, audio server and browser are outside both totals.

PSS (proportional set size) comes from each process's Linux `smaps_rollup` and
apportions shared pages, unlike adding RSS which counts shared pages repeatedly.
Values are MiB (1,048,576 bytes). RSS is included in the raw aggregate data for
reference, not used for the savings percentage. Swap PSS is reported separately;
this host was under memory pressure, so resident RAM alone is not the entire
memory footprint. No caches were flushed, swap disabled or desktop settings changed.

## Procedure

- Arch Linux x86_64, Hyprland Wayland, Intel Core i7-1360P, approximately 15.3 GiB RAM.
- Official Linux GUI `1.2.95.453.g0eeebbed`, downloaded from Spotify's official Debian
  repository and checked against the vendor package index's SHA-256. It was run
  from an extracted package with its normal sandbox, not installed system-wide.
- TUI `0.4.0`, Linux native artifact downloaded from
  [GitHub Actions run 40746914309](https://github.com/wrestle-R/Spotify-Tui/actions/runs/40746914309),
  source `1a885b86944ceb15788693b35940ec811c447756`. Its release archive checksums
  were verified and the real downloaded pair installed into a private managed prefix.
- The real account/provider session and normal cached library were used. The TUI
  was started in a dedicated Kitty window, with its normal Linux tray enabled.
- 60-second warm-up, then three consecutive windows of 20 samples each, with a
  two-second wait after every sample (about 40 seconds plus sampling overhead per window). These are **not three independent app restarts**.
- Both applications remained open simultaneously. PSS accounting therefore shares
  common pages with each other and with the rest of this running desktop.
- Headline values are medians across all 60 samples; the percentage is
  `(1 − median(TUI total PSS) / median(GUI total PSS)) × 100`.

The machine had other applications and development processes open, including a website build during sampling. Allocation, page reclaim and swapping
change with workload. Different libraries, terminal emulators, plugins, GPU drivers
and playback states will produce different results. This does not measure active
playback, cold launch peaks, CPU, battery life, or macOS memory use. An eight-hour
stability run is still outstanding.

## Reproduce it

Use Python 3 and a Linux account permitted to read its own processes' smaps files.
Prepare a **private** selectors JSON with groups named `gui`, `frontend_and_tray`,
`rust_helper`, `soloist`, and `terminal`. Each group accepts `executables` (absolute,
resolved `/proc/PID/exe` paths) and/or `roots` (PIDs). Include all detached providers
and tray processes explicitly; parent-only selection will miss them. The sampler
also follows descendants and refuses overlapping initial selectors or missing groups.
Use a dedicated terminal PID, not every terminal on the desktop.

```bash
python3 scripts/benchmark-memory.py \
  --selectors /private/path/selectors.json \
  --output /private/path/result.json \
  --warmup 60 --windows 3 --samples 20 --interval 2
```

Check both playback states before and after; discard a run if playback or process
membership changes unexpectedly. Permission failures abort instead of undercounting.
Keep selector paths and account logs private. The exported aggregate contains no
PIDs, process arguments, track titles, account identifiers, tokens or filesystem paths.

[Aggregate measurements](benchmarks/linux-paused-2026-10-29.json) include per-group
PSS/RSS/swap medians and ranges, process counts and all three window medians.
[Sampler source](../scripts/benchmark-memory.py) is part of the repository.
