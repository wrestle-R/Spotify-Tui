# Future GitHub build and release workflow

This document plans automation and acceptance for a later complete app. The requested Linux download now has a separately labelled `v0.2.0-alpha.1` development prerelease with a verified archive, local installer and download script. Offline downloads remain unsupported, so it does not pass the complete-product gates below. No GitHub Actions workflow is installed.

## Pull requests and branch changes

Run Go formatting/vet/tests/race checks, Rust formatting/clippy/tests, and shared protocol fixture validation. Use synthetic test data. Live-account tests are manual and must not expose credentials to forked pull requests. Pin third-party GitHub Actions to reviewed commit SHAs when creating the workflow.

## Release candidate

Build matching Go and Rust versions from one commit. Record compiler versions, lockfiles, supported OS/architecture, dynamic audio dependencies, and provider requirements. Generate a license inventory. Exercise both executables together on Arch Linux before tagging an accepted release.

Archive only intended executables, licenses, user instructions, and version metadata. Include no `.env`, credentials, private library data, local notes, debug logs, or upstream binaries without permission. Generate checksums over the final archive, then inspect and test the same archive that will be uploaded.

## Gated release

Before tagging the working app, link dated evidence for real login, required playlist access, local full-track playback, offline cold start, and resource budgets. A passing CI build cannot approve offline capability. All required gates in ACCEPTANCE.md must pass.

Tag the verified commit, upload the verified assets, download them again through the same path as the installer, compare checksums, extract in a clean directory, and run version/doctor checks followed by the Arch smoke test. Document known limits. Never replace a binary in-place under an unchanged version after users have verified its checksum.

## Installer validation

Cover missing dependencies, private-repo authentication failure, no network, checksum mismatch, unsupported architecture, path spaces, existing installation, interrupted update, rollback, and uninstall. Do not destroy old runtime state during upgrade. Credential migration, if ever needed, requires an explicit tested migration step.

## Planning release

The current `planning-v0.1.0` release packages documents and a verified-download script only. Its validation is archive integrity and document completeness. It makes no working-player claim and uses no Spotify credentials.
