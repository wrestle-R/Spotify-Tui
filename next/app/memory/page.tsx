import type { Metadata } from "next";
import fs from "node:fs";
import path from "node:path";
import Link from "next/link";
import { Markdown } from "@/components/markdown";
export const metadata: Metadata = { title: "The RAM test" };
export default function Memory() {
  const content = fs.readFileSync(path.resolve(process.cwd(), "../docs/MEMORY.md"), "utf8");
  return <main id="main" className="release-page wrap">
    <Link href="/#memory" className="eyebrow">← BACK TO THE NUMBERS</Link>
    <Markdown base="docs">{content}</Markdown>
  </main>;
}
