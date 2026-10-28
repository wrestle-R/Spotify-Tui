import Link from "next/link";
export default function NotFound() {
  return (
    <main id="main" className="release-page wrap">
      <div className="eyebrow">404 / NOTHING PLAYING HERE</div>
      <h1>Lost the track.</h1>
      <p className="lede">
        This page doesn’t exist. Your music is probably fine.
      </p>
      <Link href="/" className="button primary">
        Back home ↗
      </Link>
    </main>
  );
}
