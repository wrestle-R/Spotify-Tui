# Feasibility result

Reviewed again on 28 October 2026 after the user requested implementation, browser sign-in, and proper offline tests.

Later on the same date, the user explicitly requested starting implementation. The independent Go/Rust components have now been built and real browser/library/local-playback checks have passed. See [IMPLEMENTATION.md](IMPLEMENTATION.md) for the current results. The older planning-stage results below are historical and must not be read as the latest code status. The full offline requirement remains unsatisfied.

## Decision

**The full requested product is currently blocked by the mandatory standalone Spotify offline requirement.** No suitable offline provider has been established. The original provider-capability finding remains; later implementation tests in [IMPLEMENTATION.md](IMPLEMENTATION.md) supersede the historical planning results below.

The user's latest message authorizes implementation and browser-based login. A subsequent clarification reaffirmed Premium and mandatory offline playback. It does not remove offline downloads from the requirements, authorize replacing Spotify audio with another service, or authorize keeping the official desktop player open. Premium is account eligibility; it does not itself supply a custom player's missing offline API.

## Evidence checked

| Candidate | Evidence | Consequence |
| --- | --- | --- |
| librespot | [Maintainer discussion #610](https://github.com/librespot-org/librespot/discussions/610), response dated 6 November 2025 | Its encrypted cache needs network access to decrypt; cache is not a standalone offline library |
| Spotify Soloist | [Features](https://developer.spotify.com/documentation/soloist/features) and [CLI](https://developer.spotify.com/documentation/soloist/reference/command-line), rechecked 28 October 2026 | Provides headless playback/control and cache settings, but reviewed documentation does not establish music-download management or offline cold-start support |

Browser authorization supplies account access; it cannot add an absent provider playback capability. A local cache or a buffered track continuing after connection loss does not meet the requirement.

This does not establish that no separately authorized Spotify integration could ever provide offline playback. It establishes that the evaluated routes are insufficient to promise and implement the requested complete player now.

## Historical planning-stage tests

| Item | Result |
| --- | --- |
| Project documents and relative file links | PASS |
| Shell syntax and download-script help | PASS |
| Private GitHub repository and main branch | PASS |
| Planning-release assets uploaded | PASS |
| Published script downloads the planning archive through GitHub | PASS |
| Downloaded archive SHA-256 verification | PASS |
| Extracted PLAN.md equals project source | PASS |
| Published script equals project source | PASS |
| Private notes and Git metadata excluded from archive | PASS |
| Existing destination rejected | PASS |
| Unknown option rejected | PASS |
| Deliberately incorrect checksum prevents extraction | PASS |
| Temporary directory cleaned after checksum failure | PASS |
| Secondary HTTP/1-only download experiment with a 60-second cap | TIMEOUT; not required by the successful normal script path |
| Browser Spotify authorization with a real application | NOT RUN |
| Account API/library access | NOT RUN |
| Go/Rust application build | NOT RUN at the planning stage |
| Standalone full-track playback | NOT RUN |
| Downloaded-track cold start without network | BLOCKED by unavailable provider route |
| Memory/CPU savings | NOT MEASURED |

The GitHub transfer was slow but the normal end-to-end script completed and verified successfully. Do not describe the separate timed-out experiment as evidence that the successful path failed, or claim all networks have been tested.

## Necessary product decision

The user has reaffirmed the original offline requirement. Keep the project blocked on a qualifying provider. A streaming-only TUI, an interface relying on official Spotify for offline, or a local-file player would each change the scope and has not been accepted.

Until that choice or a qualifying provider changes the conditions, do not publish a working-app release or add a Downloads control that cannot perform the promised operation.
