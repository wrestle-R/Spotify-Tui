"use client";

import { useEffect, useState } from "react";
import { usePathname } from "next/navigation";

// One observer for the whole route. Content stays readable without JavaScript.
export function Motion() {
  const pathname = usePathname();
  const [paused, setPaused] = useState(false);

  useEffect(() => {
    document.documentElement.dataset.motionPaused = String(paused);
    return () => { delete document.documentElement.dataset.motionPaused; };
  }, [paused]);

  useEffect(() => {
    const root = document.documentElement;
    const preference = window.matchMedia("(prefers-reduced-motion: reduce)");
    let observer: IntersectionObserver | undefined;
    const revealAll = () => {
      document.querySelectorAll<HTMLElement>("[data-reveal]").forEach((el) => {
        el.dataset.visible = "true";
      });
    };
    const setup = () => {
      observer?.disconnect();
      if (preference.matches || !("IntersectionObserver" in window)) {
        revealAll();
        delete root.dataset.motionReady;
        return;
      }
      observer = new IntersectionObserver((entries) => {
        entries.forEach(({ target, isIntersecting }) => {
          if (isIntersecting) (target as HTMLElement).dataset.visible = "true";
          if (target.hasAttribute("data-loop")) {
            (target as HTMLElement).dataset.running = String(isIntersecting);
          } else if (isIntersecting) observer?.unobserve(target);
        });
      }, { threshold: 0.08 });
      document.querySelectorAll<HTMLElement>("[data-reveal], [data-loop]").forEach((el) => {
        // Don't hide content already in view or above a restored scroll position.
        if (el.getBoundingClientRect().top < window.innerHeight) el.dataset.visible = "true";
        observer?.observe(el);
      });
      root.dataset.motionReady = "true";
    };
    const visibility = () => { root.dataset.pageHidden = String(document.hidden); };
    const focus = (event: FocusEvent) => {
      if (event.target instanceof Element) {
        const section = event.target.closest<HTMLElement>("[data-reveal]");
        if (section) section.dataset.visible = "true";
      }
    };
    setup();
    visibility();
    preference.addEventListener("change", setup);
    document.addEventListener("visibilitychange", visibility);
    document.addEventListener("focusin", focus);
    return () => {
      observer?.disconnect();
      revealAll();
      delete root.dataset.motionReady;
      delete root.dataset.pageHidden;
      preference.removeEventListener("change", setup);
      document.removeEventListener("visibilitychange", visibility);
      document.removeEventListener("focusin", focus);
    };
  }, [pathname]);

  return (
    <button className="motion-control" type="button" aria-pressed={paused}
      onClick={() => setPaused(!paused)}>
      <span aria-hidden="true">{paused ? "▷" : "Ⅱ"}</span>
      {paused ? "Resume animations" : "Pause animations"}
    </button>
  );
}
