import Image from "next/image";
import Link from "next/link";
import { download, latest, repo } from "@/lib/releases";

// Build the decorative curve once during static export, with no client runtime.
const wave = Array.from({ length: 161 }, (_, i) => {
  const x = i * 7.5;
  const envelope = 5 + 17 * Math.sin((Math.PI * x) / 1200) ** 2;
  const y = 50 + envelope * Math.sin((Math.PI * 2 * x) / 150);
  return `${i === 0 ? "M" : "L"}${x},${y.toFixed(2)}`;
}).join(" ");

export function Footer() {
  return (
    <footer className="site-footer">
      <div className="footer-wave" aria-hidden="true">
        <div className="footer-wave-track">
          {[0, 1].map((copy) => (
            <svg key={copy} viewBox="0 0 1200 100" preserveAspectRatio="none">
              <path className="footer-wave-halo" d={wave} />
              <path className="footer-wave-line" d={wave} />
            </svg>
          ))}
        </div>
      </div>
      <div className="footer-inner wrap">
        <div className="footer-main">
          <div className="footer-identity">
            <Link
              href="/"
              className="footer-brand"
              aria-label="Spotify TUI home"
            >
              <span className="footer-logo">
                <Image src="/logo-glass.png" width={44} height={44} alt="" />
              </span>
              <span>
                spotify<span className="footer-brand-light">tui</span>
                <span className="footer-brand-dot">.</span>
              </span>
            </Link>
            <p>Made coz RAM matters.</p>
          </div>
          <div className="footer-right">
            <nav aria-label="Footer navigation">
              <Link href="/docs">Setup</Link>
              <Link href="/changelog">Changelog</Link>
              <a href={repo}>
                GitHub <span aria-hidden="true">↗</span>
              </a>
            </nav>
            <a className="footer-download" href={download()}>
              Get the player
              <span aria-hidden="true">↗</span>
            </a>
          </div>
        </div>
        <div className="footer-bottom">
          <span>
            © {new Date().getFullYear()} Spotify TUI{" "}
            <span className="footer-divider">/</span>{" "}
            <a href={`${repo}/releases/tag/${latest()}`}>{latest()}</a>
          </span>
          <span>Independent project. Not affiliated with Spotify.</span>
          <a href="#main" className="footer-top">
            Back to top <span aria-hidden="true">↑</span>
          </a>
        </div>
      </div>
    </footer>
  );
}
