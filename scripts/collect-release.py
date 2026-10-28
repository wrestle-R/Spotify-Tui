#!/usr/bin/env python3
"""Reject incomplete/mismatched matrix outputs before publishing any release."""
import hashlib
from pathlib import Path
import re
import shutil
import sys

source, destination, tag = Path(sys.argv[1]), Path(sys.argv[2]), sys.argv[3]
if not re.fullmatch(r'v[0-9]+\.[0-9]+\.[0-9]+', tag):
    raise SystemExit('Stable version tag required')
version = tag[1:]
files = {}
for directory in sorted(source.iterdir()):
    sums = directory / 'SHA256SUMS'
    if not sums.is_file():
        raise SystemExit('Missing native job checksum manifest')
    for line in sums.read_text().splitlines():
        parts = line.split()
        if len(parts) != 2 or not re.fullmatch(r'[a-f0-9]{64}', parts[0]) or Path(parts[1]).name != parts[1]:
            raise SystemExit('Invalid artifact checksum record')
        name = parts[1]
        path = directory / name
        if not path.is_file() or path.is_symlink():
            raise SystemExit('Missing regular artifact: '+name)
        if hashlib.sha256(path.read_bytes()).hexdigest() != parts[0]:
            raise SystemExit('Native artifact checksum mismatch')
        if name in files and files[name].read_bytes() != path.read_bytes():
            raise SystemExit('Native jobs produced conflicting common assets')
        files[name] = path
required = {'spotify-tui_'+version+'_linux_amd64.tar.gz',
            'spotify-tui_'+version+'_darwin_arm64.tar.gz', 'download.py', 'download-linux.sh'}
if set(files) != required:
    raise SystemExit('Incomplete or unexpected release asset set: '+str(set(files)))
destination.mkdir()
for name,path in sorted(files.items()):
    shutil.copy2(path,destination/name)
(destination/'SHA256SUMS').write_text(''.join(hashlib.sha256((destination/name).read_bytes()).hexdigest()+'  '+name+'\n' for name in sorted(files)))
print('Both platform archives and matching download scripts verified for',tag)
