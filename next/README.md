# spotify tui website

Landing page, setup guide and changelog. Next.js App Router with a static export. No accounts, secrets, tracking or backend services. Website sign-in instructions lead to the app's local browser flow; the website itself doesn't handle Spotify credentials.

From this directory:

```bash
npm ci
npm run dev
```

Open http://127.0.0.1:3000. Production verification: `npm run check && npm run build`. The static output is `out/`; deploy that directory to a static host. No deployment has been configured automatically.

Release pages are read from `../docs/releases/` and the changelog from `../CHANGELOG.md` at build time. Keep the repository root available when building. The newest ordinary version determines the download links. The hero uses separate emerald-glass headphone artwork. Logo assets come from the approved music-note concept; the screenshot is the generic demo, never account data.

The setup guide follows Install → Sign in → Pair → Play. Platform selection keeps
archive links and Soloist commands together. Each stage has a visible checkpoint;
copy buttons select the command text when clipboard permission is unavailable.
The website never receives account credentials. Reading sections stay still while
navigation, platform choices and command feedback have short transitions.

Build the production export, then preview it with:

```bash
python3 ../scripts/serve-website.py
```

Open http://127.0.0.1:4173. Browser regressions run on the production export:

```bash
npx playwright install chromium
npm run test:e2e -- --project=chromium
```

CI also runs these flows in WebKit on macOS 26, including platform-specific links,
keyboard navigation, copy failure recovery, 320px layout, reduced motion and
JavaScript-disabled reading. Chromium tests real clipboard read/write; WebKit runs
the denial/manual-selection path because its clipboard permission API differs.

Separately, each native Linux/macOS 15/macOS 26 job downloads the latest **published
stable** release, verifies the downloader checksum, uses its real download and
bundled install scripts in an isolated prefix, checks both installed versions and
runs the PTY walkthrough. This tests the public download route alongside the fresh
source build. It does not sign in to a Spotify account or prove audible Mac playback.
