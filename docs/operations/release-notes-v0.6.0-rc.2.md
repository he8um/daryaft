# Daryaft v0.6.0-rc.2 Internal Release Candidate

Status: internal validation only

`v0.6.0-rc.2` is an internal release candidate for validating the current
pre-1.0 Daryaft foundation. It supersedes `v0.6.0-rc.1`. It is not a public
stable release and does not enable package-manager publishing or public
install-channel guarantees.

Public stable release remains planned for `v1.0.0`.

## Why rc.2

`v0.6.0-rc.2` supersedes `v0.6.0-rc.1` with the following improvements to the
CI security posture and release validation workflow:

- Restored blocking `govulncheck ./...` in the CI `security` job after
  upgrading CI to Go `1.26.4` or newer. The previous rc.1 advisory gap
  (GO-2026-5039 and GO-2026-5037 in Go `1.26.3` standard library) is fully
  resolved.
- `make rc-check` now includes blocking `govulncheck` and `gosec` checks; it
  no longer skips `govulncheck` during release-candidate validation.
- Real-terminal interactive TUI QA completed: download, inspect, checksum,
  batch, cancel, and config flows verified in a live terminal.
- All product behavior from rc.1 is carried forward unchanged.

## Highlights

- CLI single URL downloads over HTTP/HTTPS.
- TUI download flows for single URLs and `.txt` URL files.
- Sequential CLI batch downloads with continue-on-error reporting.
- Retry with exponential backoff and resume with `.part` files and
  `.part.daryaft.json` sidecar metadata.
- Safe CLI and TUI cancellation that preserves partial state for resume.
- CLI and TUI inspect flows for read-only URL metadata, with JSON output.
- CLI and TUI checksum verification for single URL downloads (`sha256` and
  `sha512`).
- Built-in default output directory to the user's `Downloads` directory.
- YAML configuration with `config path`, `config show`, `config set`,
  `config get`, `config reset`, and `config init` commands.
- `DARYAFT_*` environment variable overrides.
- Shell completion generation for bash, zsh, fish, and PowerShell.
- `doctor` diagnostics with human-readable, JSON, and strict modes.
- `version` and `version --json` with GoReleaser ldflags injection.
- GoReleaser v2 snapshot and release-check readiness for local validation.
- Local and CI quality gates: tests, builds, race tests, lint, govulncheck,
  gosec, GoReleaser config validation.

## Security and Quality Posture

- `govulncheck ./...` is blocking in CI (`security` job).
- `gosec ./...` is blocking in CI (`security` job).
- `make security` passes locally with Go `1.26.4` or newer; no vulnerabilities
  found.
- `make rc-check` includes blocking `govulncheck` and `gosec` checks.
- `make release-check` validates local GoReleaser snapshot artifacts without
  publishing.
- GitHub Actions for `v0.6.0-rc.2` passed (Go test/build matrix on Linux and
  macOS, goreleaser-check, lint, security).
- Real-terminal TUI QA completed and passed.

## Known Limitations

- This is not a public stable release.
- Windows is not officially tested or supported yet.
- Self-update is not implemented.
- Proxy, custom headers, and auth are not implemented.
- Concurrent and segmented downloads are not implemented.
- Batch checksum semantics are not implemented.
- Checksum file discovery and signed checksum verification are not implemented.
- Queue and history are not implemented.
- Package manager publishing is not enabled.
- Public stable release remains planned for `v1.0.0`.

## Validation Commands

Run the following from the repository root:

```bash
git fetch --tags
git tag --list "v0.6.0-rc.*"
git describe --tags --always
go run . version
go run . version --json
make rc-info
make rc-check
make security
goreleaser check
make release-check
find dist -maxdepth 2 -type f | sort
git status --short --ignored dist bin
```

To inspect the tag without switching branches:

```bash
git show --stat v0.6.0-rc.2
```

Use [Release-Candidate Validation](rc-validation.md) for the full validation
workflow and finding-record guidance.

## No-Publish Note

This RC is a validation marker only. It does not imply package-manager
publishing, public GitHub Release publishing, or stable install support.

If creating a GitHub release manually for internal sharing, mark it as
pre-release:

```bash
gh release create v0.6.0-rc.2 \
  --title "Daryaft v0.6.0-rc.2" \
  --notes-file docs/operations/release-notes-v0.6.0-rc.2.md \
  --prerelease \
  --verify-tag
```

Do not run this unless intentionally publishing a GitHub pre-release. Do not
enable package-manager artifact publishing.

If validation finds a blocker, fix the finding on `main` and create a later
internal RC tag after confirming quality gates and GitHub Actions are healthy.
Do not retag or overwrite `v0.6.0-rc.2`.
