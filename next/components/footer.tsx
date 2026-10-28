import Image from "next/image";
import Link from "next/link";
import { download, latest, repo } from "@/lib/releases";

export function Footer() {
  const version = latest();
  return (
    <footer className="site-footer">
      <div className="footer-inner wrap">
        <div className="footer-sleeve">
          <div className="footer-editorial">
            <span className="footer-kicker">
              <span aria-hidden="true">↳</span> leave a little room for music
            </span>
            <h2>
              One more
              <br />
              <span>song.</span>
            </h2>
            <div className="footer-action-row">
              <a className="footer-download" href={download()}>
                Get the player
                <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path d="M12 4v15m-6-6 6 6 6-6M5 21h14" />
                </svg>
              </a>
              <span className="footer-platform">
                Linux / x86_64
                <br />
                {version} stable
              </span>
            </div>
          </div>
          <div className="footer-record-scene" aria-hidden="true">
            <div className="footer-record-shadow" />
            <div className="footer-record">
              <div className="footer-record-label">
                <span className="footer-record-type">SPOTIFY TUI · SIDE A</span>
                <Image src="/logo-flat.png" width={76} height={76} alt="" />
                <span className="footer-record-hole" />
                <span className="footer-record-type">GO + RUST / 33⅓ RPM</span>
              </div>
            </div>
          </div>
          <span className="footer-sleeve-index" aria-hidden="true">
            ST—001 / MADE COZ RAM MATTERS
          </span>
        </div>
        <div className="footer-directory">
          <Link
            href="/"
            className="footer-small-brand"
            aria-label="Spotify TUI home"
          >
            <Image src="/logo-glass.png" width={24} height={24} alt="" />
            spotify tui
          </Link>
          <nav aria-label="Footer navigation">
            <Link href="/docs">Setup</Link>
            <Link href="/changelog">Changelog</Link>
            <a href={`${repo}/releases`}>Releases</a>
            <a href={repo}>
              GitHub <span aria-hidden="true">↗</span>
            </a>
          </nav>
          <a href="#main" className="footer-top" aria-label="Back to top">
            ↑
          </a>
        </div>
        <div className="footer-bottom">
          <span>© {new Date().getFullYear()} Spotify TUI</span>
          <span>
            Independent project. Not affiliated with or endorsed by Spotify.
          </span>
          <span className="footer-made">
            a little less app. a little more music.
          </span>
        </div>
      </div>
    </footer>
  );
}
