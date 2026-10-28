import type { CSSProperties } from "react";
import Link from "next/link";
import fs from "node:fs";
import path from "node:path";

type Range = { median: number; min: number; max: number };
type Result = {
  tui_total_pss_mib: Range;
  gui_total_pss_mib: Range;
  reduction_percent: number;
  groups: Record<string, { swap_pss_mib: Range }>;
};

export function MemoryStats() {
  const data: Result = JSON.parse(fs.readFileSync(path.resolve(process.cwd(),
    "../docs/benchmarks/linux-paused-2026-10-29.json"), "utf8"));
  const tui = data.tui_total_pss_mib;
  const gui = data.gui_total_pss_mib;
  const tuiSwap = Object.entries(data.groups).filter(([name]) => name !== "gui")
    .reduce((sum, [, group]) => sum + group.swap_pss_mib.median, 0);
  return <section className="memory wrap" id="memory" aria-labelledby="memory-title">
    <div className="memory-intro" data-reveal>
      <span className="eyebrow">THE REASON THIS EXISTS</span>
      <h2 id="memory-title">Okay, but how<br /><em>much RAM?</em></h2>
      <p>Same laptop, same account, both paused with their libraries loaded.
        We counted the whole TUI setup, including Soloist, the tray and its terminal.</p>
      <Link className="text-link" href="/memory">Read the test conditions ↗</Link>
    </div>
    <div className="memory-chart" data-reveal>
      <div className="memory-headline"><strong>{data.reduction_percent}%</strong>
        <span>less resident RAM<br />in this Linux test</span></div>
      <div className="memory-row">
        <div><span>Spotify desktop</span><strong>{gui.median.toFixed(1)} <small>MiB</small></strong></div>
        <div className="memory-track"><span className="memory-bar gui-bar" /></div>
        <p>{gui.min.toFixed(1)}–{gui.max.toFixed(1)} MiB across samples</p>
      </div>
      <div className="memory-row">
        <div><span>Spotify TUI + audio + terminal</span><strong>{tui.median.toFixed(1)} <small>MiB</small></strong></div>
        <div className="memory-track"><span className="memory-bar tui-bar"
          style={{ "--bar-scale": Math.min(tui.median / gui.median, 1) } as CSSProperties} /></div>
        <p>{tui.min.toFixed(1)}–{tui.max.toFixed(1)} MiB across samples</p>
      </div>
      <p className="memory-note">PSS medians · 60 samples · 29 Oct 2026 · Arch x86_64<br />
        Swap PSS: GUI {data.groups.gui.swap_pss_mib.median.toFixed(1)} MiB, TUI groups {tuiSwap.toFixed(1)} MiB.
        One short test on a busy laptop, not a universal saving or a Mac benchmark.</p>
    </div>
  </section>;
}
