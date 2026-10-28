#!/usr/bin/env python3
"""Preview the static Next export locally, including its extensionless routes."""
import argparse
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import urlsplit

root = Path(__file__).resolve().parents[1] / 'next/out'


class Handler(SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=str(root), **kwargs)

    def do_GET(self):
        route = urlsplit(self.path).path
        candidate = (root / (route.strip('/') + '.html')).resolve()
        if route != '/' and candidate.is_relative_to(root) and candidate.is_file():
            self.path = route.rstrip('/') + '.html'
        super().do_GET()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--port', type=int, default=4173)
    args = parser.parse_args()
    if not (root / 'index.html').is_file():
        parser.error('Build first: npm run build --prefix next')
    print(f'Production preview: http://127.0.0.1:{args.port}', flush=True)
    ThreadingHTTPServer(('127.0.0.1', args.port), Handler).serve_forever()
