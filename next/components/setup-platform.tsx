"use client";

import { createContext, useContext, useState } from "react";
import { Code } from "@/components/code";

type Platform = "linux" | "mac";
const PlatformContext = createContext<{ platform: Platform; select: (value: Platform) => void }>({
  platform: "linux", select: () => {},
});

export function SetupPlatform({ children }: { children: React.ReactNode }) {
  const [platform, select] = useState<Platform>("linux");
  return <PlatformContext.Provider value={{ platform, select }}>{children}</PlatformContext.Provider>;
}

export function InstallPlayer({ version, linux, mac }: { version: string; linux: string; mac: string }) {
  const { platform, select } = useContext(PlatformContext);
  const isMac = platform === "mac";
  const folder = `spotify-tui-${version.replace(/^v/, "")}`;
  return <>
    <div className="platform-picker" role="group" aria-label="Choose your operating system">
      <button type="button" aria-pressed={!isMac} onClick={() => select("linux")}>Linux <span>x86_64</span></button>
      <button type="button" aria-pressed={isMac} onClick={() => select("mac")}>Mac <span>Apple Silicon</span></button>
    </div>
    <div className="platform-description" aria-live="polite" aria-atomic="true">
      {isMac ? <p>Apple Silicon (M1 or later), macOS 15+. Intel Macs aren’t supported. These binaries are ad-hoc signed, not Apple-notarized; see the <a href="https://github.com/wrestle-R/Spotify-Tui/blob/main/docs/MACOS.md#signing-and-gatekeeper">Mac signing notes</a>.</p>
        : <p>Linux x86_64 with glibc 2.39+. Audio needs PipeWire or PulseAudio. ARM Linux isn’t included in this release.</p>}
    </div>
    <p>Both platforms need <strong>Python 3.9+</strong>, <code>curl</code> and a terminal. The downloaded app itself doesn’t need Go, Rust or Python. No sudo.</p>
    <h3>1. Download the installer</h3>
    <Code label="Download script">{`curl -fL --retry 2 https://github.com/wrestle-R/Spotify-Tui/releases/download/${version}/download.py -o download.py`}</Code>
    <p>Open <code>download.py</code> in your editor to review it, or use the command below. Press <kbd>q</kbd> to leave the reader before continuing.</p>
    <Code label="Review script">less download.py</Code>
    <h3>2. Download and install the matching pair</h3>
    <p>The downloader selects this computer’s OS and architecture, verifies the archive checksum, then extracts it. The installer runs only if that succeeds.</p>
    <Code label="Download and install">{String.raw`python3 download.py download --version ${version.replace(/^v/, "")} \
  --destination ./${folder} && \
  bash ./${folder}/install.sh`}</Code>
    <p>This installs the player and its Rust helper under <code>~/.local</code>, keeping older builds for rollback. You can use the full commands below without changing your PATH.</p>
    <details><summary>Already downloaded it, or prefer a manual archive?</summary>
      <p>If <code>{folder}</code> already exists from a successful download, run its <code>install.sh</code> again. If it was an incomplete attempt, choose a new destination and use that same directory in the installer command. The downloader never overwrites an existing directory.</p>
      <p><a href={isMac ? mac : linux}>Download {isMac ? "Apple Silicon" : "Linux x86_64"} archive</a>. Extract the archive, open a terminal inside its folder and run <code>bash ./install.sh</code>. <a href={`https://github.com/wrestle-R/Spotify-Tui/releases/tag/${version}`}>Checksums and all release files ↗</a>.</p>
    </details>
    <h3>3. Check the installation</h3>
    <Code label="Verify installed pair">{`~/.local/bin/spotify-tui version
~/.local/bin/spotify-tui-playback --version
~/.local/bin/spotify-tui --demo`}</Code>
    <p>The two version commands should both report <strong>{version.replace(/^v/, "")}</strong>. The demo should open with sample music. It has no audio; press <kbd>q</kbd> to close it.</p>
    <noscript><p>Platform buttons need JavaScript. Manual archives: <a href={linux}>Linux x86_64</a> / <a href={mac}>Apple Silicon Mac</a>. Follow the <a href="https://github.com/wrestle-R/Spotify-Tui/blob/main/docs/MACOS.md">Mac setup guide</a> for Mac-specific instructions.</p></noscript>
  </>;
}

export function InstallSoloist() {
  const { platform } = useContext(PlatformContext);
  const archive = platform === "mac" ? "soloist_release_osx_arm64.tar.gz" : "soloist_release_x86_64.tar.gz";
  return <>
    <p>Open <a href="https://developer.spotify.com/documentation/soloist/reference/downloads-and-updates">Spotify’s official Soloist downloads ↗</a> and download <strong>{platform === "mac" ? "Apple Silicon Mac (arm64)" : "x86_64 Linux"}</strong> to your <code>Downloads</code> folder. Soloist is Spotify’s audio engine and is downloaded separately.</p>
    <p>The expected filename is <code>{archive}</code>. These commands extract it into a stable location, then check it starts:</p>
    <Code label="Install official audio engine">{String.raw`mkdir -p "$HOME/.local/lib/spotify-soloist" && \
  tar -xzf "$HOME/Downloads/${archive}" \
    -C "$HOME/.local/lib/spotify-soloist" && \
  "$HOME/.local/lib/spotify-soloist/soloist" --version`}</Code>
    <p>If your browser saved the archive elsewhere, adjust the path after <code>-xzf</code>. Continue when Soloist prints its version. Its location must stay the same after setup.</p>
  </>;
}
