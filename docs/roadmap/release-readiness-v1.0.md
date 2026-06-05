# Release Readiness: v1.0.0

This document defines what is required before a public stable `v1.0.0` release
and what is explicitly deferred to post-1.0. It is a planning reference, not an
implementation spec.

## Release Strategy

v1.0.0 is a **stable baseline release**, not a feature-complete release.

The goal is to ship the currently implemented and validated feature set under a
stable version, with known limitations clearly documented. Features outside that
baseline are post-1.0 work.

Current state: `v0.6.0-rc.2` is the current internal release candidate.
GitHub pre-release published. See
[v0.6.0-rc.2 Release Status](../operations/release-status-v0.6.0-rc.2.md).

## v1.0.0 Release Criteria

All of the following must be satisfied before tagging `v1.0.0`.

### Quality Gates
- [x] GitHub Actions green on the release commit (Go test/build matrix on Linux
      and macOS, goreleaser-check, lint, security) — green at `v0.6.0-rc.2`.
- [x] `govulncheck` blocking in CI — no vulnerabilities found.
- [x] `gosec` blocking in CI — Issues: 0.
- [x] `make rc-check` passes with blocking security checks.
- [x] `make release-check` passes (local GoReleaser snapshot, no publishing).
- [x] `go test ./...` passes.
- [x] `go build ./...` passes.
- [x] Race tests pass: `go test -race ./internal/downloader` and
      `go test -race ./internal/tui`.

Note: these gates must be re-verified green on the final release commit before
tagging `v1.0.0`.

### Validation
- [x] Real-terminal interactive TUI QA passed (completed at `v0.6.0-rc.2`).
- [ ] No blocker bugs identified from internal RC validation.
- [x] Clean-directory install-and-use validation from source and GoReleaser
      snapshot artifacts on the `v0.6.0-rc.2` tag — PASS WITH NOTES. See
      [Clean Install Validation: v0.6.0-rc.2](../operations/clean-install-validation-v0.6.0-rc.2.md).

### Release Artifacts and Notes
- [x] Binary asset strategy decided — binary archives will be attached to the
      v1.0.0 GitHub release. See
      [v1.0.0 Release Assets](../operations/release-assets.md).
- [x] v1.0.0 release notes drafted. See
      [Daryaft v1.0.0 Release Notes (draft)](../operations/release-notes-v1.0.0.md).
- [ ] Release notes finalized and accurate for the tagged commit.
- [x] Known limitations clearly documented in release notes.

### Documentation
- [ ] README reflects the stable release.
- [ ] Installation docs reflect the actual install method for v1.0.0.
- [ ] Known limitations listed (see below).

## Known Limitations to Document at v1.0.0

The following are known limitations that must be clearly documented in the
v1.0.0 release notes and README, but do not block the release:

- Windows is not officially supported or tested.
- Self-update is not implemented.
- Proxy, custom headers, and auth are not implemented.
- Concurrent and segmented downloads are not implemented.
- Checksum verification is supported for single URL downloads only; batch
  checksum semantics, checksum file auto-discovery, and signed checksum
  verification are not implemented.
- Queue and history are not implemented.
- Package manager publishing is not enabled at v1.0.0.

## Explicitly NOT Required Before v1.0.0

The following items are deferred to post-1.0. They are not blockers.

- Windows official support and CI.
- Package manager publishing (Homebrew, deb, rpm, Arch).
- Self-update.
- Proxy, custom headers, and auth.
- Concurrent downloads.
- Segmented downloads.
- Queue and history.
- Checksum file auto-discovery.
- Signed checksum verification.
- Batch checksum semantics.
- Decisions on all of the above — they can remain "deferred" in v1.0.0.

## Post-1.0 Roadmap

Features and decisions deferred until after a stable v1.0.0 baseline:

- **Install channel expansion**: Homebrew tap, apt/deb repo, rpm, Arch,
  Scoop, or other package manager publishing.
- **Windows support**: Add Windows CI and verify binary builds, or explicitly
  document the support tier.
- **Self-update**: Decide on a self-update mechanism and implement it.
- **Proxy, custom headers, and auth**: HTTP proxy support, per-request custom
  headers, and authentication flows.
- **Concurrency and segmentation**: Parallel multi-segment downloads for faster
  throughput on large files.
- **Queue and history**: Download queue management and persistent history.
- **Advanced checksum support**: Checksum file auto-discovery, signed checksum
  verification, and batch checksum semantics.
- **Broader platform QA**: Testing on additional Linux distributions, macOS
  versions, and (if supported) Windows.
- **Long-term compatibility policy**: Minimum Go version policy, OS support
  matrix, backport policy.

## Reference

- [v1.0.0 Go/No-Go Checklist](../operations/v1.0.0-go-no-go.md)
- [v1.0.0 Release Plan](../operations/v1.0.0-release-plan.md)
- [v1.0.0 Release Notes (draft)](../operations/release-notes-v1.0.0.md)
- [v1.0.0 Release Assets](../operations/release-assets.md)
- [v0.6.0-rc.2 Release Status](../operations/release-status-v0.6.0-rc.2.md)
- [Release-Candidate Validation](../operations/rc-validation.md)
- [Pre-Release Readiness](../operations/pre-release-readiness.md)
- [Pre-1.0 Roadmap](pre-1-roadmap.md)
- [Post-1.0 Feature Packs](post-1-feature-packs.md)
