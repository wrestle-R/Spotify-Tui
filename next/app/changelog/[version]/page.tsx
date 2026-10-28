import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { Markdown } from "@/components/markdown";
import { versions, release, repo } from "@/lib/releases";
export function generateStaticParams() {
  return versions().map((version) => ({ version }));
}
export const dynamicParams = false;
export async function generateMetadata({
  params,
}: {
  params: Promise<{ version: string }>;
}): Promise<Metadata> {
  return { title: `Release ${(await params).version}` };
}
export default async function Release({
  params,
}: {
  params: Promise<{ version: string }>;
}) {
  const { version } = await params;
  const content = release(version);
  if (!content) notFound();
  return (
    <main id="main" className="release-page wrap">
      <Link href="/changelog" className="eyebrow">
        ← ALL RELEASES
      </Link>
      <Markdown base="docs/releases">{content}</Markdown>
      <a href={`${repo}/releases/tag/${version}`} className="button primary">
        Download {version} ↓
      </a>
    </main>
  );
}
