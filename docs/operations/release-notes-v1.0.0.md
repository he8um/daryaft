# Daryaft v1.0.0

> **Status: RELEASED** — Stable release. Tag: `v1.0.0`.

## Summary

Daryaft v1.0.0 is the first stable baseline release.

This release ships the CLI and TUI foundation that has been validated through
the `v0.6.0-rc.2` internal release candidate. It is a stable baseline, not a
feature-complete terminal downloader. Known limitations are documented below
and post-1.0 features are tracked separately.

## Highlights

- **Single URL HTTP/HTTPS download** with safe `.part` write and final rename.
- **Resume support** via HTTP Range requests and `.part.daryaft.json` metadata
  sidecars. Resumes correctly on server support; restarts from byte 0 on servers
  that do not support Range.
- **Sequential batch downloads** from multiple URL arguments or a `.txt` URL
  file. Continue-on-error with final summary.
- **Retry with exponential backoff** for transient network and server failures.
  Configurable up to 20 retries.
- **CLI checksum verification** (`--checksum sha256:<hex>` / `sha512:<hex>`)
  for single URL downloads. Verifies after final file rename.
- **Dry-run planning** (`--dry-run`) validates inputs and prints the download
  plan without writing files.
- **`inspect` command** shows HTTP metadata (URL, filename, content-length,
  content-type, resume support, ETag, Last-Modified) without downloading.
- **Interactive TUI** (`daryaft` with no arguments): URL and `.txt` file input,
  single URL with optional custom filename and optional checksum, batch download,
  Inspect URL flow, dry-run plan screen, and download execution with live
  progress. Resizes to terminal window.
- **YAML configuration** in `<UserConfigDir>/daryaft/config.yaml` with
  `DARYAFT_*` environment variable overrides. Keys: `download_dir`, `retries`,
  `resume`, `no_color`, `no_tui`, `theme`.
- **`doctor` diagnostics** with human output, `--json` for automation, and
  `--strict` for CI (warnings become non-zero exit).
- **Shell completions** for bash, zsh, fish, and PowerShell.
- **`--verbose` flag** for download diagnostics including URL (redacted),
  output directory, filename, HTTP status, resume/retry decisions, and duration.
- **CI quality gates**: Go test/build matrix on Linux and macOS, goreleaser-check,
  golangci-lint, govulncheck (blocking, no vulnerabilities), gosec (blocking,
  Issues: 0).
- **GoReleaser v2** build for linux/amd64, linux/arm64, darwin/amd64,
  darwin/arm64 binary archives.

## Install

Binary archives are attached to this GitHub release:

```
daryaft_linux_amd64.tar.gz
daryaft_linux_arm64.tar.gz
daryaft_darwin_amd64.tar.gz
daryaft_darwin_arm64.tar.gz
checksums.txt
```

Download the archive for your platform, verify the checksum, extract, and place
the `daryaft` binary on your PATH:

```bash
# Example: macOS Apple Silicon
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.0.0/daryaft_darwin_arm64.tar.gz
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.0.0/checksums.txt
shasum -a 256 --check checksums.txt
tar -xzf daryaft_darwin_arm64.tar.gz
./daryaft version
```

Expected: `version: 1.0.0`, `built_by: goreleaser`.

Or build from source with Go 1.26.4 or newer:

```bash
git clone https://github.com/he8um/daryaft.git
cd daryaft
git checkout v1.0.0
go build -o daryaft .
./daryaft version
```

Package manager channels (Homebrew, deb, rpm, Arch) are post-1.0 work and are
not available at v1.0.0.

## Validation Summary

The `v0.6.0-rc.2` release candidate completed the following validation passes
before v1.0.0 was tagged:

- GitHub Actions green on the release commit: Go test/build (Linux + macOS),
  goreleaser-check, lint, security — all PASS.
- `govulncheck ./...` — no vulnerabilities found.
- `gosec ./...` — Issues: 0.
- `make rc-check` — PASS with blocking security checks.
- `make release-check` (local GoReleaser snapshot) — PASS.
- Real-terminal interactive TUI QA — PASS.
- Clean-directory install-and-use validation from source and GoReleaser snapshot
  artifacts on the `v0.6.0-rc.2` tag — PASS WITH NOTES (see
  [Clean Install Validation: v0.6.0-rc.2](clean-install-validation-v0.6.0-rc.2.md)).

## Known Limitations

The following limitations are present at v1.0.0 and are documented, not hidden.
They are not regressions; they are features deferred to post-1.0.

- **Windows**: Not officially supported or tested. Binaries are not provided for
  Windows at v1.0.0.
- **Concurrent and segmented downloads**: Not implemented. All downloads are
  single-connection. High-throughput parallel downloads are post-1.0.
- **Batch checksum verification**: `--checksum` is single URL only. Checksum
  file auto-discovery (`.sha256`, `.sha512`, `.sha256sum`) and signed checksum
  verification are not implemented.
- **Proxy, custom headers, and authentication**: Not implemented. Daryaft
  follows standard HTTP, including environment proxy variables, but has no
  explicit proxy configuration, per-request custom headers, or auth flows.
- **Self-update**: Not implemented. Update by downloading a new binary or
  rebuilding from source.
- **Queue and history**: Not implemented. Download queue management and
  persistent history are post-1.0.
- **Package manager publishing**: Homebrew, deb, rpm, Arch, and Scoop are
  post-1.0. No package channel is available at v1.0.0.
- **Checksum mismatch handling**: On a checksum mismatch, the completed
  download file is preserved (not deleted). The user must decide whether to
  keep or remove it.

## Upgrade Notes

There is no prior stable release. This is the first public stable release.

If you have been using the `0.6.0-dev` source builds:
- Run `./daryaft version` to confirm version `1.0.0` after upgrade.
- The config file format and location are unchanged.
- Existing `.part` files and `.part.daryaft.json` metadata sidecars are
  compatible.
- No migration steps are required.

## Post-1.0 Roadmap

Features deferred until after v1.0.0:

- Package manager publishing (Homebrew, deb, rpm, Arch, Scoop).
- Windows official support and CI.
- Self-update mechanism.
- Proxy, custom headers, and authentication.
- Concurrent and segmented downloads.
- Queue and history.
- Checksum file auto-discovery and signed checksum verification.
- Batch checksum semantics.
- Automated release pipeline.

See [Post-1.0 Feature Packs](../roadmap/post-1-feature-packs.md) for details.

## References

- [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md)
- [v1.0.0 Release Assets](release-assets.md)
- [v0.6.0-rc.2 Release Status](release-status-v0.6.0-rc.2.md)
- [Clean Install Validation: v0.6.0-rc.2](clean-install-validation-v0.6.0-rc.2.md)
- [Release Process](release-process.md)
