"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
export function NavLinks({ repo }: { repo: string }) {
  const pathname = usePathname();
  return <nav aria-label="Main navigation">
    {[['/docs','Docs'],['/changelog','Changelog']].map(([href,label]) =>
      <Link key={href} href={href} aria-current={pathname === href || pathname.startsWith(href + '/') ? 'page' : undefined}>{label}</Link>)}
    <a href={repo} className="github">GitHub <span aria-hidden="true">↗</span></a>
  </nav>;
}
