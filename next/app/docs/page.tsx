import type { Metadata } from "next";
import Link from "next/link";
import { Code } from "@/components/code";
import { download, latest, repo } from "@/lib/releases";
export const metadata: Metadata = { title: "Get started" };
const sections = [
  ["before", "Before you start"],
  ["install", "Download & install"],
  ["sign-in", "Browser sign-in"],
  ["audio", "Connect local audio"],
  ["keyboard", "Keyboard & tray"],
  ["updates", "Updates"],
  ["troubleshooting", "Troubleshooting"],
  ["limits", "Current limits"],
];
export default function Docs() {
  return (
    <main id="main" className="doc-layout wrap">
      <aside className="doc-nav">
        <span className="eyebrow">THE MANUAL</span>
        {sections.map(([id, title]) => (
          <a href={`#${id}`} key={id}>
            {title}
          </a>
        ))}
        <Link href="/changelog">Changelog ↗</Link>
      </aside>
      <article className="doc-content">
        <div className="eyebrow">GET STARTED / {latest()}</div>
        <h1>
          A few steps.
          <br />
          <em>Then hit play.</em>
        </h1>
        <p className="lede">
          Two connections, one player. Browser sign-in loads your library.
          Spotify Soloist handles the actual audio.
        </p>
        <section id="before">
          <span className="step">00 / THE CHECKLIST</span>
          <h2>Before you start</h2>
          <ul>
            <li>
              Linux x86_64 with glibc 2.39+, or Apple Silicon Mac with macOS 15+. Python 3.9+ for installation.
            </li>
            <li>
              Spotify Premium to generate your own personal Soloist key. Don’t
              use someone else’s key.
            </li>
            <li>
              A Spotify Developer app with Web API access, a normal browser, and
              a working desktop keyring.
            </li>
            <li>
              Linux uses PipeWire or PulseAudio. macOS uses the official Soloist Mac build. The optional StatusNotifier tray is Linux-only.
            </li>
          </ul>
          <p>
            This is an independent player. The website never collects Spotify
            passwords, tokens or provider keys.
          </p>
        </section>
        <section id="install">
          <span className="step">01 / INSTALL</span>
          <h2>Download your platform release</h2>
          <p>
            <a href={download()}>Linux x86_64 ↓</a> or <a href={download("darwin_arm64")}>Apple Silicon Mac ↓</a> ({latest()}). You can also use the commands below. The downloader verifies SHA-256 before
            extracting. Read the script before running it.
          </p>
          <Code>{`curl -fL https://github.com/wrestle-R/Spotify-Tui/releases/download/${latest()}/download.py -o download.py
less download.py
python3 download.py download --destination ./spotify-tui-download
bash ./spotify-tui-download/install.sh
~/.local/bin/spotify-tui --demo`}</Code>
          <p>
            Mac binaries are ad-hoc signed, not Apple-notarized. Intel Macs are not supported in this release. No sudo. The installer puts both binaries under{" "}
            <code>~/.local</code>. The demo uses sample music and makes no
            sound. Add <code>~/.local/bin</code> to your PATH if you want to run{" "}
            <code>spotify-tui</code> directly.
          </p>
        </section>
        <section id="sign-in">
          <span className="step">02 / YOUR LIBRARY</span>
          <h2>Sign in through your browser</h2>
          <ol>
            <li>
              Open the{" "}
              <a href="https://developer.spotify.com/dashboard">
                Spotify Developer dashboard ↗
              </a>{" "}
              in your usual browser. Create or select your own app with Web API
              access.
            </li>
            <li>
              Add this exact Redirect URI:{" "}
              <code>http://127.0.0.1:43821/callback</code>. Allowlist your
              account if the dashboard requires it.
            </li>
            <li>
              Copy the public <strong>client ID</strong>. You don’t need the
              client secret.
            </li>
            <li>Run setup, then login:</li>
          </ol>
          <Code>{`~/.local/bin/spotify-tui setup --client-id YOUR_PUBLIC_CLIENT_ID
~/.local/bin/spotify-tui login`}</Code>
          <p>
            Keep the command running. Your normal browser opens Spotify’s
            consent screen and returns to the local callback on this computer.
            The terminal confirms when sign-in is saved in your keyring. Don’t
            sign in through an embedded browser if Spotify verification fails
            there.
          </p>
          <details>
            <summary>Browser didn’t open?</summary>
            <p>
              Run <code>~/.local/bin/spotify-tui login --no-browser</code> and
              open the printed URL in your normal browser on the same computer.
              The callback expires after five minutes. Finish Spotify’s
              verification flow there, then retry login if needed.
            </p>
          </details>
        </section>
        <section id="audio">
          <span className="step">03 / THE AUDIO</span>
          <h2>Connect the local player</h2>
          <p>
            Library sign-in doesn’t start audio. Install the separate official
            Soloist binary from{" "}
            <a href="https://developer.spotify.com/documentation/soloist/reference/downloads-and-updates">
              Spotify’s downloads page ↗
            </a>
            . It isn’t bundled in this release.
          </p>
          <ol>
            <li>
              Put <code>soloist</code> on your PATH, or supply its full path
              below.
            </li>
            <li>
              Generate your own{" "}
              <a href="https://developer.spotify.com/dashboard/soloist">
                personal Soloist key ↗
              </a>
              . The key owner needs Premium.
            </li>
            <li>
              Save the key using the hidden terminal prompt. Never put it in a
              command, screenshot, issue or this website.
            </li>
          </ol>
          <Code>{`~/.local/bin/spotify-tui setup --client-id YOUR_PUBLIC_CLIENT_ID --soloist-path /absolute/path/to/soloist
~/.local/bin/spotify-tui provider-key
~/.local/bin/spotify-tui connect`}</Code>
          <p>
            For first-time pairing, open Spotify on the same local network and
            select <strong>Spotify TUI Local</strong> in the device picker. Then
            close the official app and launch:
          </p>
          <Code>{`~/.local/bin/spotify-tui`}</Code>
          <p>
            Select a collection, press Enter to open it, select a song and press
            Enter to play. Keep the app’s data directory so Soloist can restore
            its session later. Restoration can sometimes stall; <kbd>c</kbd>{" "}
            restarts the helper.
          </p>
        </section>
        <section id="keyboard">
          <span className="step">04 / MAKE YOURSELF AT HOME</span>
          <h2>Keys you’ll actually remember</h2>
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>Key</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                {[
                  ["↑ ↓ / j k", "Move in the focused pane"],
                  ["Enter / →", "Open collection; Enter on a song plays"],
                  ["Esc / ←", "Back to library"],
                  ["Tab / Shift-Tab", "Switch panes"],
                  ["/", "Search; Enter submits, Esc cancels"],
                  ["Space", "Pause / resume"],
                  ["n / p", "Next / previous song"],
                  ["[ / ]", "Seek five seconds"],
                  ["+ / −", "Change volume"],
                  ["c / L", "Reconnect / browser sign-in"],
                  ["Ctrl-U", "Check, install stable update and restart"],
                  ["q / Q", "Quit, keep music / stop music and quit"],
                  ["?", "Full keyboard help"],
                ].map(([key, action]) => (
                  <tr key={key}>
                    <td>
                      <kbd>{key}</kbd>
                    </td>
                    <td>{action}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <h3>The tray</h3>
          <p>
            The player starts a small tray companion when your panel has a
            StatusNotifier host. Left-click opens the TUI; the menu offers
            pause/resume, next, previous, and stop. <kbd>q</kbd> leaves music
            and tray running; <kbd>Q</kbd> stops both.
          </p>
          <p>
            On Hyprland, enable the tray component in Quickshell or your panel.
            GNOME may need an AppIndicator extension. The app doesn’t install or
            modify your panel configuration. Run <code>spotify-tui doctor</code>{" "}
            to check the host, or <code>spotify-tui --no-tray</code> to skip the
            tray. Opening from the tray requires Kitty, Foot, Alacritty,
            Konsole, GNOME Terminal or xterm.
          </p>
        </section>
        <section id="updates">
          <span className="step">05 / KEEP IT CURRENT</span>
          <h2>Stay in the terminal</h2>
          <p>
            The app checks the latest stable release in the background when
            opened. When an update is available, press <kbd>Ctrl-U</kbd>. It
            downloads the archive, verifies SHA-256 and both executable
            versions, switches the installed pair together, then restarts.
            Existing build directories stay available for rollback.
          </p>
          <Code>{`~/.local/bin/spotify-tui check-updates
~/.local/bin/spotify-tui update`}</Code>
          <p>
            In-app updates need the managed installation from{" "}
            <code>install.sh</code>. They don’t overwrite source builds. A
            failed download or verification leaves the installed version intact.
            No automatic installation happens without Ctrl-U or the update
            command.
          </p>
        </section>
        <section id="troubleshooting">
          <h2>If something gets stuck</h2>
          <details>
            <summary>The device isn’t in Spotify’s picker</summary>
            <p>
              Run <code>spotify-tui connect</code>, then{" "}
              <code>spotify-tui doctor</code>. Both devices need local-network
              reachability. Discovery across phone USB tethering can be
              unreliable; try Spotify on this computer for the initial pairing.
              Don’t delete credentials or provider data first.
            </p>
          </details>
          <details>
            <summary>No library or search results</summary>
            <p>
              Press L or run <code>spotify-tui login</code> again. Check your
              client ID, redirect URI and allowed users in the developer
              dashboard. Spotify can restrict reading followed or curated
              playlists.
            </p>
          </details>
          <details>
            <summary>No keyring available</summary>
            <p>
              Start a desktop secret service such as GNOME Keyring or KWallet
              and unlock it. There is an explicit plaintext fallback described
              in the repository setup guide; it is not enabled by default.
            </p>
          </details>
          <details>
            <summary>Need diagnostic commands?</summary>
            <Code>{`~/.local/bin/spotify-tui doctor
~/.local/bin/spotify-tui stop
~/.local/bin/spotify-tui tray-stop`}</Code>
            <p>
              Doctor reports status without printing your keys. Read{" "}
              <a href={`${repo}/blob/main/docs/SETUP.md`}>
                the full setup guide ↗
              </a>{" "}
              for build instructions and integration caveats.
            </p>
          </details>
        </section>
        <section id="limits">
          <h2>What’s still missing</h2>
          <p>
            <strong>Genuine offline downloads are unsupported.</strong> Premium
            doesn’t provide this integration with a download API. The isolated
            offline cold-start check failed.
          </p>
          <p>
            Search currently returns ten results. A liked song plays
            individually; full liked-song sequencing and playlist offsets are
            unfinished. MPRIS/media keys and long-session memory testing are
            still pending. Connectivity and Soloist session restoration can
            occasionally stall. Press c to reconnect.
          </p>
          <p>
            This project does not promise ad-free playback for Free accounts.
            Each installation needs its own valid provider setup. See{" "}
            <a href={`${repo}/blob/main/docs/IMPLEMENTATION.md`}>
              verification details ↗
            </a>
            .
          </p>
        </section>
      </article>
    </main>
  );
}
