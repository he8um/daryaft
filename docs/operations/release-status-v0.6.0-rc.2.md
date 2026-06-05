# Daryaft v0.6.0-rc.2 Release Status

## Release Summary

| Field            | Value                                                         |
|------------------|---------------------------------------------------------------|
| Tag              | `v0.6.0-rc.2`                                                 |
| Release name     | Daryaft v0.6.0-rc.2                                           |
| Release type     | GitHub pre-release (internal/validation; not stable)          |
| Draft            | No                                                            |
| Pre-release flag | Yes                                                           |
| Created          | 2026-06-05T21:11:34Z                                          |
| Published        | 2026-06-05T22:15:30Z                                          |
| URL              | https://github.com/he8um/daryaft/releases/tag/v0.6.0-rc.2    |
| Public stable    | Not yet — planned for `v1.0.0`                                |

## CI Status

- GitHub Actions workflow ran on the `v0.6.0-rc.2` tag push.
- All jobs passed:
  - `Go test/build (ubuntu-latest)` — PASS
  - `Go test/build (macos-latest)` — PASS
  - `goreleaser-check` — PASS
  - `lint` — PASS
  - `security` — PASS

## Security Status

- `govulncheck ./...` is blocking in CI. No vulnerabilities found.
- `gosec ./...` is blocking in CI. Issues: 0.
- `make security` passes locally with Go `1.26.4` or newer.
- The previous Go `1.26.3` standard-library advisory gap (GO-2026-5039 and
  GO-2026-5037) is resolved.

## QA Status

- Full automated QA pass: tests, builds, race tests, lint, govulncheck, gosec,
  goreleaser check — all passed.
- Real-terminal interactive TUI QA passed: download, inspect, checksum, batch,
  cancel, and config flows verified in a live terminal.
- `make rc-check` passed with blocking security checks.
- `make release-check` passed; local GoReleaser snapshot artifacts generated
  and verified without publishing.

## Release Artifact Status

The GitHub pre-release at `v0.6.0-rc.2` is **source-code and release-notes
only**. No compiled binary assets are attached.

This is acceptable for internal validation:

- Reviewers can clone the tag and build from source.
- `make release-check` generates local snapshot artifacts for inspection.
- Source-only is a deliberate choice for pre-stable internal RCs.

For a wider pre-release intended for external users, consider attaching
GoReleaser-built binary archives and a `checksums.txt` file. This should be an
explicit decision, not automatic. See the asset recommendation section below.

## Publish Status

- GitHub pre-release published at the URL above.
- Not marked stable. Not a public stable release.
- No package-manager publishing (Homebrew, deb, rpm, Arch).
- No public install-channel guarantee.
- No self-update support.

## Asset Recommendation

For the current internal validation pre-release, source-tag-only is sufficient.

When the project is ready for a wider pre-release or public release, attach
compiled binary artifacts. The recommended process:

1. Build with `goreleaser release` (not `--snapshot`) on a clean tag.
2. Upload `dist/*.tar.gz`, `dist/*.zip`, and `dist/checksums.txt` to the
   GitHub release.
3. Verify checksums against the uploaded files before announcing.

Do not implement binary asset upload from this RC. Make it an explicit release
decision when the project is ready for wider distribution.

## Current Recommendation

- Continue internal validation on `v0.6.0-rc.2`.
- If blockers are found, fix on `main`, confirm quality gates and CI pass,
  and create `v0.6.0-rc.3`.
- If no blockers are found, proceed toward `v1.0.0` stable release.

v1.0.0 is a **stable baseline release** of the current feature set. No
additional product features are required. The remaining steps are: completing
RC validation, clean install/use testing, release notes, and deciding on binary
assets. Post-1.0 features (Windows, self-update, proxy/auth, package managers,
concurrency, etc.) are deferred and not blockers.

Public stable release is planned for `v1.0.0`.

See [Release-Candidate Validation](rc-validation.md) for the full validation
workflow. See [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md)
for the exact v1.0.0 criteria and the post-1.0 roadmap.
