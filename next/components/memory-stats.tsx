import type { CSSProperties } from "react";
import Link from "next/link";
import fs from "node:fs";
import path from "node:path";

type Range = { median: number; min: number; max: number };
type Result = {
  tui_total_pss_mib: Range;
  gui_total_pss_mib: Range;
};

export function MemoryStats() {
  const data: Result = JSON.parse(fs.readFileSync(path.resolve(process.cwd(),
    "../docs/benchmarks/linux-paused-2026-10-29.json"), "utf8"));
  // Keep the measured source intact; the landing page uses a rounded illustration.
  const tui = Math.round(data.tui_total_pss_mib.median / 10) * 10;
  const gui = Math.round(data.gui_total_pss_mib.median / 100) * 100;
  const reduction = Math.round((1 - tui / gui) * 10) * 10;
  return <section className="memory wrap" id="memory" aria-labelledby="memory-title">
    <div className="memory-intro" data-reveal>
      <span className="eyebrow">A LITTLE LESS RAM, MAYBE A LOT</span>
      <h2 id="memory-title">Okay, but how<br /><em>much RAM?</em></h2>
      <p>A smaller player could leave more RAM for whatever you’re working on.
        Here’s a rounded example of what that might look like. Your setup could be different.</p>
      <Link className="text-link" href="/memory">See the measurements behind this example ↗</Link>
    </div>
    <div className="memory-chart" data-reveal>
      <div className="memory-headline"><strong>~{reduction}%</strong>
        <span>less resident RAM<br />in this example</span></div>
      <div className="memory-row">
        <div><span>Spotify desktop</span><strong>~{gui} <small>MiB</small></strong></div>
        <div className="memory-track"><span className="memory-bar gui-bar" /></div>
        <p>Illustrative desktop footprint</p>
      </div>
      <div className="memory-row">
        <div><span>Spotify TUI + audio + terminal</span><strong>~{tui} <small>MiB</small></strong></div>
        <div className="memory-track"><span className="memory-bar tui-bar"
          style={{ "--bar-scale": Math.min(tui / gui, 1) } as CSSProperties} /></div>
        <p>Illustrative full setup, including audio, tray and terminal</p>
      </div>
      <p className="memory-note">Illustrative figures, rounded from one short Linux paused-library test.<br />
        Actual RAM use and savings vary with your machine, library and workload.
        These aren’t typical-use estimates or Mac measurements.</p>
    </div>
  </section>;
}
