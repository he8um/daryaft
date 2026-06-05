# Daryaft v0.6.0-rc.1 Internal Release Candidate

Status: superseded — use `v0.6.0-rc.2` for current internal validation

> **Superseded.** `v0.6.0-rc.1` has been superseded by
> [`v0.6.0-rc.2`](release-notes-v0.6.0-rc.2.md). Use `v0.6.0-rc.2` for
> current internal validation. Historical content below is preserved for
> reference.

`v0.6.0-rc.1` is an internal release candidate for validating the current
pre-1.0 Daryaft foundation. It is not a public stable release and does not
enable package-manager publishing or public install-channel guarantees.

Public stable release remains planned for `v1.0.0`.

## Highlights

- CLI single URL downloads.
- Sequential CLI batch downloads.
- TUI download flows for single URLs and `.txt` URL files.
- CLI and TUI inspect flows for read-only URL metadata.
- Retry and resume with `.part` files and `.part.daryaft.json` sidecar
  metadata.
- Safe CLI and TUI cancellation that preserves partial state for resume.
- CLI and TUI checksum verification for single URL downloads.
- Built-in default output directory to the user's `Downloads` directory.
- YAML configuration and `DARYAFT_*` environment overrides.
- Shell completion generation for bash, zsh, fish, and PowerShell.
- Doctor diagnostics with human, JSON, and strict modes.
- GoReleaser snapshot readiness for local validation.
- Local and CI quality gates for tests, builds, lint, GoReleaser config, and
  security scanning.

## Known Limitations

- This is not a public stable release.
- Windows is not officially tested or supported yet.
- Self-update is not implemented.
- Proxy, custom headers, and auth are not implemented.
- Concurrent and segmented downloads are not implemented.
- Batch checksum semantics are not implemented.
- Checksum file discovery and signed checksum verification are not implemented.
- At the time of the initial RC tag, Go `1.26.3` caused `govulncheck` to report
  standard-library advisories `GO-2026-5039` and `GO-2026-5037`; both are fixed
  in Go `1.26.4`. This gap is now resolved: CI uses Go `1.26.4` or newer and
  `govulncheck` is blocking with no findings.
- Full interactive TUI QA should still be completed in a real terminal before
  wider public release.

## Validation Commands

```bash
git fetch --tags
git tag --list "v0.6.0-rc.*"
git describe --tags --always
go run . version
go run . version --json
goreleaser check
make rc-check
make release-check
find dist -maxdepth 2 -type f | sort
git status --short --ignored dist bin
```

Use [Release-Candidate Validation](rc-validation.md) for the full validation
workflow and finding-record guidance.

## Rollback And No-Publish Note

This RC is a validation marker only. It does not imply package-manager
publishing, public GitHub Release publishing, or stable install support.

If validation finds a blocker, fix the finding on `main` and create a later
internal RC tag after confirming the quality gates and GitHub Actions are
healthy. Do not retag or overwrite `v0.6.0-rc.1`.
