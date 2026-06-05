# QA Results: 0.6.0-dev

Date: 2026-06-05

Version: `0.6.0-dev`

Verdict: PASS WITH NOTES

Scope: internal validation readiness

## Summary

The automated and terminal-driven QA pass completed successfully for internal
`0.6.0-dev` release-candidate validation readiness.

- Tests, build, and targeted race checks passed.
- `make lint` passed.
- `gosec ./...` passed.
- `goreleaser check` passed.
- `make rc-check` passed.
- `make release-check` passed.
- CLI inspect passed.
- CLI dry-run passed.
- CLI single download passed.
- CLI batch download passed.
- CLI checksum passed.
- Config and environment override checks passed.
- Doctor checks passed.
- Shell completion generation passed.
- Release tooling snapshot passed.

## Known Notes

- At the time of this QA pass (Go `1.26.3`), `govulncheck` reported known Go
  standard-library advisories `GO-2026-5039` and `GO-2026-5037`.
- These advisories are fixed in Go `1.26.4`. This was tracked as a toolchain
  patch gap, not a Daryaft source-code finding.
- That gap is now resolved: CI and local tooling use Go `1.26.4` or newer,
  and `govulncheck` is blocking in CI with no vulnerabilities reported.
- Windows is not officially tested or supported yet.
- Full interactive TUI QA was completed in a real terminal during `v0.6.0-rc.2`
  validation; all flows passed.

## Findings

- No blocker findings.
- No source-code security findings from `gosec`.
- No release publishing occurred.
- No tag was created.

## Next Decision

- The current internal RC is `v0.6.0-rc.2`. GitHub pre-release published. See
  [Release-Candidate Validation](rc-validation.md),
  [Daryaft v0.6.0-rc.2 Internal Release Candidate](release-notes-v0.6.0-rc.2.md),
  and [v0.6.0-rc.2 Release Status](release-status-v0.6.0-rc.2.md).
- `v0.6.0-rc.1` is superseded; see
  [Daryaft v0.6.0-rc.1 Internal Release Candidate](release-notes-v0.6.0-rc.1.md)
  for historical notes.
- Do not publish a public stable release yet.
- Public stable remains planned for `v1.0.0`.
