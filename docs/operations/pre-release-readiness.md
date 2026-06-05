# Pre-Release Readiness Review

## Status

- Daryaft current development version: `0.6.0-dev`.
- This is not a public stable release.
- Public stable release remains planned for `v1.0.0`.
- The project is suitable for continued local/internal pre-release validation.

## Readiness Verdict

Verdict: Ready for internal/manual `0.6.0-dev` validation, not ready for public
stable `v1.0.0`.

Core CLI and TUI functionality exists, stabilization work is complete, quality
gates exist, and release tooling is validated locally. Public install channels,
Windows support, and final `v1.0.0` guarantees are not ready.

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

`make release-check` remains local/manual and does not publish. `gosec` is
blocking in CI. `govulncheck` is temporarily advisory in CI only until Go
`1.26.4` or newer is available, while local `make security` remains strict.

## Known Toolchain/Security Note

On Go `1.26.3`, `govulncheck` may report vulnerabilities in Go standard
library packages. The reported fixed version is Go `1.26.4`.

This should be resolved by using a patched Go toolchain when available. This is
not currently identified as a Daryaft source-code vulnerability. Do not suppress
it permanently.

Re-run `make security` after the local and CI Go toolchains receive the patch.
Restore CI `govulncheck` to blocking after patched Go is available.

## Manual QA

Use the [Manual QA Checklist](manual-qa.md) for local pre-release validation.
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

GoReleaser local snapshot checks are available. Package manager publishing
remains disabled/future work. Self-update remains future work.

## Not Ready For v1.0.0 Because

- Public install channels are not enabled.
- Windows is not officially tested/supported in CI.
- Self-update is not implemented.
- Proxy/custom headers/auth are not implemented.
- Batch checksum semantics, signed checksums, and checksum file auto-discovery
  are not implemented.
- Concurrent/segmented downloads are not implemented.
- Queue/history is not implemented.
- More manual QA is needed before public release.
- Security gate should be re-run with patched Go toolchain when available.

## Recommended Next Work

- Run full manual QA.
- Fix any QA findings.
- Add inspect polish if needed.
- Later: proxy/auth, batch checksum semantics, signed checksums, concurrent
  batch, segmented downloads, queue/history.
- Consider Windows CI only when ready to support Windows officially.
- Restore `govulncheck` blocking in CI after Go `1.26.4` or newer is
  available.

## Decision Checklist

- [ ] GitHub Actions green.
- [ ] Local tests pass.
- [ ] Local lint passes.
- [ ] Local security passes or known Go toolchain issue documented.
- [ ] Manual QA completed.
- [ ] No unexpected generated artifacts.
- [ ] `git status` clean.
- [ ] No tag/release created accidentally.
- [ ] Known limitations documented.
