#!/usr/bin/env python3
"""Linux PSS/RSS sampler. Private selectors in; anonymous aggregate statistics out."""
import argparse
from datetime import datetime, timezone
import json
import os
from pathlib import Path
import platform
import statistics
import time


def snapshot(selectors):
    processes = {}
    for path in Path('/proc').iterdir():
        if not path.name.isdigit():
            continue
        try:
            exe = os.readlink(path / 'exe')
            status = (path / 'status').read_text()
            parent = int(next(line.split()[1] for line in status.splitlines() if line.startswith('PPid:')))
            processes[int(path.name)] = (exe, parent)
        except (FileNotFoundError, PermissionError, ProcessLookupError):
            continue
    assigned = {}
    for name, selector in selectors.items():
        for pid, (exe, _) in processes.items():
            if exe in selector.get('executables', []) or pid in selector.get('roots', []):
                if pid in assigned:
                    raise ValueError('Overlapping selectors; each process must belong to one group')
                assigned[pid] = name
    # Include descendants, including renderers with different executable names.
    changed = True
    while changed:
        changed = False
        for pid, (_, parent) in processes.items():
            if pid not in assigned and parent in assigned:
                assigned[pid] = assigned[parent]
                changed = True
    groups = {name: {'pss_kib': 0, 'rss_kib': 0, 'swap_pss_kib': 0, 'processes': 0} for name in selectors}
    for pid, name in assigned.items():
        try:
            fields = {}
            for line in Path(f'/proc/{pid}/smaps_rollup').read_text().splitlines():
                key, _, value = line.partition(':')
                if key in ('Pss', 'Rss', 'SwapPss'):
                    fields[key] = int(value.split()[0])
        except (FileNotFoundError, ProcessLookupError):
            raise RuntimeError('Process changed during sample; rerun after startup settles') from None
        # Permission errors must fail, never silently undercount.
        groups[name]['pss_kib'] += fields['Pss']
        groups[name]['rss_kib'] += fields['Rss']
        groups[name]['swap_pss_kib'] += fields.get('SwapPss', 0)
        groups[name]['processes'] += 1
    if any(not group['processes'] for group in groups.values()):
        raise RuntimeError('A required process group is absent')
    return groups


def summary(values):
    return {key: round(fn(values), 2) for key, fn in (
        ('median', statistics.median), ('min', min), ('max', max))}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--selectors', type=Path, required=True, help='Private JSON: group -> executables and/or root PIDs')
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--warmup', type=float, default=60)
    parser.add_argument('--windows', type=int, default=3)
    parser.add_argument('--samples', type=int, default=20)
    parser.add_argument('--interval', type=float, default=2)
    args = parser.parse_args()
    if platform.system() != 'Linux':
        parser.error('Linux /proc/smaps_rollup required; this is not a macOS benchmark')
    if min(args.windows, args.samples, args.interval) <= 0 or args.warmup < 0:
        parser.error('Invalid sampling parameters')
    selectors = json.loads(args.selectors.read_text())
    time.sleep(args.warmup)
    windows = []
    raw = []
    for index in range(args.windows):
        samples = []
        for _ in range(args.samples):
            samples.append(snapshot(selectors))
            time.sleep(args.interval)
        raw.extend(samples)
        windows.append({name: summary([s[name]['pss_kib'] / 1024 for s in samples]) for name in selectors})
        print(f'Window {index + 1}/{args.windows} complete', flush=True)
    combined = [sum(s[name]['pss_kib'] for name in selectors if name != 'gui') / 1024 for s in raw]
    gui = [s['gui']['pss_kib'] / 1024 for s in raw] if 'gui' in selectors else []
    output = {
        'measured_at_utc': datetime.now(timezone.utc).isoformat(),
        'platform': {'os': platform.system(), 'kernel': platform.release(), 'architecture': platform.machine()},
        'method': {'metric': 'PSS', 'unit': 'MiB', 'warmup_seconds': args.warmup,
                   'windows': args.windows, 'samples_per_window': args.samples, 'interval_seconds': args.interval,
                   'simultaneous': True, 'independent_restarts': False},
        'groups': {name: {
            'pss_mib': summary([s[name]['pss_kib'] / 1024 for s in raw]),
            'rss_mib': summary([s[name]['rss_kib'] / 1024 for s in raw]),
            'swap_pss_mib': summary([s[name]['swap_pss_kib'] / 1024 for s in raw]),
            'process_count': summary([s[name]['processes'] for s in raw]),
        } for name in selectors},
        'windows_pss_mib': windows,
        'tui_total_pss_mib': summary(combined),
    }
    if gui:
        output['gui_total_pss_mib'] = summary(gui)
        output['reduction_percent'] = round((1 - statistics.median(combined) / statistics.median(gui)) * 100, 1)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(output, indent=2) + '\n')
    print('Saved aggregate statistics; no PIDs, paths, arguments or account data included')


if __name__ == '__main__':
    main()
