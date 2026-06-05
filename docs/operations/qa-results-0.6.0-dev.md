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

- On Go `1.26.3`, `govulncheck` reports known Go standard-library advisories:
  - `GO-2026-5039`
  - `GO-2026-5037`
- These advisories are fixed in Go `1.26.4`.
- This is tracked as a toolchain patch gap, not a Daryaft source-code finding.
- CI currently treats `govulncheck` as advisory while keeping `gosec`
  blocking.
- Local `make security` remains strict.
- Windows is not officially tested or supported yet.
- Full interactive TUI QA should still be completed in a real terminal before a
  wider public release, although automated TUI tests and race tests passed.

## Findings

- No blocker findings.
- No source-code security findings from `gosec`.
- No release publishing occurred.
- No tag was created.

## Next Decision

- Continue internal validation without a tag.
- Create an internal release-candidate tag such as `v0.6.0-rc.1` after
  confirming GitHub Actions remains green.
- Do not publish a public stable release yet.
- Public stable remains planned for `v1.0.0`.
