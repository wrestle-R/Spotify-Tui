import type { Metadata } from "next";
import { Markdown } from "@/components/markdown";
import { changelog, repo } from "@/lib/releases";
export const metadata: Metadata = { title: "Changelog" };
export default function Changelog() {
  return (
    <main id="main" className="release-page wrap">
      <div className="eyebrow">THE RECORD / STABLE RELEASES</div>
      <h1>
        A little better.
        <br />
        <em>Release by release.</em>
      </h1>
      <p className="lede">
        What changed, what was checked, and what still needs work.
      </p>
      <Markdown>{changelog().replace(/^# What's new\s*/, "")}</Markdown>
      <a href={`${repo}/releases`} className="text-link">
        All GitHub releases ↗
      </a>
    </main>
  );
}
