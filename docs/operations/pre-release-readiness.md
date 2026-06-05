# Pre-Release Readiness Review

## Status

- Daryaft current development version: `0.6.0-dev`.
- This is not a public stable release.
- Public stable release remains planned for `v1.0.0`.
- The project is suitable for continued local/internal pre-release validation.

## Readiness Verdict

Verdict: Ready for internal/manual `0.6.0-dev` validation. v1.0.0 requires
completing RC validation, clean install/use verification, and release notes —
not additional features.

Core CLI and TUI functionality exists, stabilization work is complete, quality
gates exist, and release tooling is validated locally. v1.0.0 is a stable
baseline release of the current feature set. Post-1.0 features (Windows,
self-update, proxy/auth, concurrency, package managers, etc.) are not blockers.

The completed QA pass is recorded in
[QA Results: 0.6.0-dev](qa-results-0.6.0-dev.md). Its verdict is PASS WITH
NOTES for internal validation readiness.

The current internal RC is `v0.6.0-rc.2`. GitHub pre-release published.
See [Release-Candidate Validation](rc-validation.md),
[Daryaft v0.6.0-rc.2 Internal Release Candidate](release-notes-v0.6.0-rc.2.md),
and [v0.6.0-rc.2 Release Status](release-status-v0.6.0-rc.2.md).
`v0.6.0-rc.1` is superseded; its notes remain at
[Daryaft v0.6.0-rc.1 Internal Release Candidate](release-notes-v0.6.0-rc.1.md).

## Implemented Capabilities

- CLI foundation.
- Single URL download.
- Manual CLI and TUI checksum verification for single URL downloads.
- Sequential batch download.
- Dry-run planning.
- Retry.
- Resume with `.part` and sidecar metadata.
- Safe CLI and TUI cancellation.
- Progress events.
- TUI download flows.
- CLI inspect.
- TUI inspect.
- YAML config.
- Environment overrides.
- Shell completion.
- `doctor`, `doctor --json`, and `doctor --strict`.
- `version` and `version --json`.
- GoReleaser v2 config.
- `make release-check`.
- Lint/security quality gates.
- Manual QA checklist.
- GitHub issue/PR templates.

## Quality Gates

Local checks:

```bash
go test ./...
go build ./...
go test -race ./internal/downloader
go test -race ./internal/tui
make rc-check
make lint
make security
goreleaser check
git diff --check
sh -n scripts/manual-qa-server.sh
```

CI checks:

- `Go test/build (ubuntu-latest)`.
- `Go test/build (macos-latest)`.
- `goreleaser-check`.
- `lint`.
- `security`.

`make release-check` remains local/manual and does not publish. `govulncheck`
and `gosec` are both blocking in CI. `make rc-check` includes blocking
`govulncheck` and `gosec` checks. `make security` is strict locally and in CI.

## Toolchain/Security Note

The previous Go `1.26.3` standard-library advisory gap (GO-2026-5039 and
GO-2026-5037) is resolved by using Go `1.26.4` or newer. CI and local
`make security` now use Go `1.26.4` or newer, and `govulncheck` reports no
vulnerabilities. This is not a Daryaft source-code finding and no suppression
is needed.

## Manual QA

Use the [Manual QA Checklist](manual-qa.md) for local pre-release validation.
The latest completed pass is documented in
[QA Results: 0.6.0-dev](qa-results-0.6.0-dev.md).
Use [Release-Candidate Validation](rc-validation.md) for RC tag and
GoReleaser artifact validation.
Manual QA should cover:

- CLI single download.
- CLI checksum verification.
- TUI single URL checksum verification.
- CLI batch download.
- Dry-run.
- Resume.
- Ctrl+C cancellation.
- Inspect.
- TUI download.
- TUI inspect.
- Config/env.
- Doctor.
- Completion.
- Release-check.

## Release Readiness

No tag should be created automatically from this document. No public release
should be published yet.

GoReleaser local snapshot checks are available. Package manager publishing,
self-update, and other post-1.0 features remain future work and are not
blockers for v1.0.0. v1.0.0 is a stable baseline release of the currently
implemented feature set. See
[Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md) for the full
v1.0.0 criteria and the post-1.0 roadmap.

## Remaining v1.0.0 Validation Steps

These are the remaining steps before v1.0.0 can be tagged. They do not require
new features — only validation and release process work.

Completed:
- Internal validation on `v0.6.0-rc.2` — no blockers found.
- Clean-directory install-and-use validation from source and GoReleaser snapshot
  artifacts on the `v0.6.0-rc.2` tag — PASS WITH NOTES.
- v1.0.0 release notes drafted. See
  [Daryaft v1.0.0 Release Notes (draft)](release-notes-v1.0.0.md).
- Binary asset decision made: binary archives will be attached. See
  [v1.0.0 Release Assets](release-assets.md).

Remaining:
- Confirm GitHub Actions green on the final release commit.
- Finalize release notes for the tagged commit.
- Tag `v1.0.0`, build release artifacts, and upload binary assets.

## Known Limitations at v1.0.0

The following are known limitations that will be clearly documented in the
v1.0.0 release notes. They are not blockers for the stable baseline release.

- Windows is not officially supported or tested.
- Self-update is not implemented.
- Proxy, custom headers, and auth are not implemented.
- Concurrent and segmented downloads are not implemented.
- Batch checksum semantics, checksum file auto-discovery, and signed checksum
  verification are not implemented.
- Queue and history are not implemented.
- Package manager publishing is not enabled at v1.0.0.

## Decision Checklist

- [ ] GitHub Actions green on the final release commit.
- [ ] Local tests pass on the final release commit.
- [ ] Local lint passes on the final release commit.
- [ ] `make rc-check` passes on the final release commit.
- [ ] `make security` passes on the final release commit.
- [x] Real-terminal TUI QA passed (completed at `v0.6.0-rc.2`).
- [x] No blocker findings from RC validation.
- [x] Clean install/use validation from source and snapshot artifacts on
      `v0.6.0-rc.2` — PASS WITH NOTES.
- [x] Release notes drafted with known limitations documented. See
      [Daryaft v1.0.0 Release Notes (draft)](release-notes-v1.0.0.md).
- [x] Binary asset decision made. See [v1.0.0 Release Assets](release-assets.md).
- [ ] Release notes finalized for the tagged commit.
- [ ] `git status` clean on the release commit.
- [ ] No tag/release created accidentally.
