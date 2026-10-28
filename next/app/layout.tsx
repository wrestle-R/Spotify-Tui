import type { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";
import { repo } from "@/lib/releases";
import { Footer } from "@/components/footer";
import "./globals.css";
export const metadata: Metadata = {
  title: {
    default: "Spotify TUI · music, minus the whole desktop app",
    template: "%s · Spotify TUI",
  },
  description:
    "A Spotify terminal player for Linux, built with Go and Rust. Playlists, liked songs and search, with a small tray companion. Made because RAM matters.",
  icons: { icon: "/logo-flat.png" },
};
export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" data-scroll-behavior="smooth">
      <body>
        <a className="skip" href="#main">
          Skip to content
        </a>
        <header className="nav">
          <Link href="/" className="brand">
            <Image src="/logo-glass.png" width={32} height={32} alt="" />
            <span>
              spotify<span className="brand-light">tui</span>
              <span className="brand-dot">.</span>
            </span>
          </Link>
          <nav aria-label="Main navigation">
            <Link href="/docs">Docs</Link>
            <Link href="/changelog">Changelog</Link>
            <a href={repo} className="github">
              GitHub <span aria-hidden="true">↗</span>
            </a>
          </nav>
        </header>
        {children}
        <Footer />
      </body>
    </html>
  );
}
