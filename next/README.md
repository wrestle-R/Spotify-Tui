# spotify tui website

Landing page, setup guide and changelog. Next.js App Router with a static export. No accounts, secrets, tracking or backend services. Website sign-in instructions lead to the app's local browser flow; the website itself doesn't handle Spotify credentials.

From this directory:

```bash
npm ci
npm run dev
```

Open http://127.0.0.1:3000. Production verification: `npm run check && npm run build`. The static output is `out/`; deploy that directory to a static host. No deployment has been configured automatically.

Release pages are read from `../docs/releases/` and the changelog from `../CHANGELOG.md` at build time. Keep the repository root available when building. The newest ordinary version determines the download links. The hero uses separate emerald-glass headphone artwork. Logo assets come from the approved music-note concept; the screenshot is the generic demo, never account data.
