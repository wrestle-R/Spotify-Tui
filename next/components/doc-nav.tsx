"use client";
import Link from "next/link";
import { useEffect, useState } from "react";

export function DocNav({ sections }: { sections: string[][] }) {
  const [active, setActive] = useState("before");
  useEffect(() => {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => { if (entry.isIntersecting) setActive(entry.target.id); });
    }, { rootMargin: "-10% 0px -65% 0px", threshold: 0 });
    sections.forEach(([id]) => { const el = document.getElementById(id); if (el) observer.observe(el); });
    const hash = () => {
      const id = window.location.hash.slice(1);
      if (sections.some(([section]) => section === id)) setActive(id);
    };
    hash();
    window.addEventListener("hashchange", hash);
    return () => { observer.disconnect(); window.removeEventListener("hashchange", hash); };
  }, [sections]);
  return <aside className="doc-nav">
    <span className="eyebrow">YOUR FIRST SONG</span>
    <nav aria-label="Setup sections">
      {sections.map(([id, title]) => <a href={`#${id}`} key={id}
        aria-current={active === id ? "location" : undefined}
        onClick={() => setActive(id)}>{title}</a>)}
    </nav>
    <Link href="/changelog">Changelog ↗</Link>
  </aside>;
}
