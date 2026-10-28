"use client";
import { useState } from "react";
export function Code({ children }: { children: string }) {
  const [state, setState] = useState("Copy");
  async function copy() {
    try {
      await navigator.clipboard.writeText(children);
      setState("Copied");
      setTimeout(() => setState("Copy"), 1800);
    } catch {
      setState("Select to copy");
    }
  }
  return (
    <div className="code">
      <pre>
        <code>{children}</code>
      </pre>
      <button type="button" onClick={copy} aria-label="Copy commands">
        {state}
      </button>
      <span className="sr-only" aria-live="polite">
        {state === "Copied" ? "Commands copied" : ""}
      </span>
    </div>
  );
}
