"use client";
import { useEffect, useRef, useState } from "react";

export function Code({ children, label = "Commands" }: { children: string; label?: string }) {
  const [state, setState] = useState("Copy");
  const code = useRef<HTMLElement>(null);
  const timer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  useEffect(() => () => clearTimeout(timer.current), []);
  async function copy() {
    clearTimeout(timer.current);
    try {
      await navigator.clipboard.writeText(children);
      setState("Copied");
      timer.current = setTimeout(() => setState("Copy"), 2200);
    } catch {
      // Clipboard permissions can be denied. Select the real text for manual copy.
      const selection = window.getSelection();
      if (code.current && selection) {
        const range = document.createRange();
        range.selectNodeContents(code.current);
        selection.removeAllRanges();
        selection.addRange(range);
        setState("Selected");
      } else setState("Select to copy");
    }
  }
  return <div className="code" data-copy-state={state.toLowerCase()}>
    <div className="code-header"><span>{label}</span>
      <button type="button" onClick={copy} aria-label={`Copy ${label.toLowerCase()}`}>
        <span aria-hidden="true">{state === "Copied" ? "✓ " : ""}</span>{state}
      </button>
    </div>
    <pre><code ref={code}>{children}</code></pre>
    <span className="sr-only" role="status">{state === "Copied" ? `${label} copied` : ""}</span>
    {(state === "Selected" || state === "Select to copy") && <p className="copy-help" role="status">
      {state === "Selected" ? "Text selected." : "Select the command text."} Press Ctrl+C or ⌘C to copy.
    </p>}
  </div>;
}
