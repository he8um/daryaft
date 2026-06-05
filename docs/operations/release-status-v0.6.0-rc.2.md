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
- Clean install validation passed. See
  [Clean Install Validation: v0.6.0-rc.2](clean-install-validation-v0.6.0-rc.2.md).

## Release Artifact Status

The GitHub pre-release at `v0.6.0-rc.2` is **source-code and release-notes
only**. No compiled binary assets are attached.

This is acceptable for internal validation:

- Reviewers can clone the tag and build from source.
- `make release-check` generates local snapshot artifacts for inspection.
- Source-only is a deliberate choice for pre-stable internal RCs.

For the `v0.6.0-rc.2` internal pre-release, source-only is correct. Binary
assets will be attached to the v1.0.0 stable release. See the asset section
below.

## Publish Status

- GitHub pre-release published at the URL above.
- Not marked stable. Not a public stable release.
- No package-manager publishing (Homebrew, deb, rpm, Arch).
- No public install-channel guarantee.
- No self-update support.

## Asset Decision

The decision has been made: binary assets will be attached to the v1.0.0
GitHub release. The `v0.6.0-rc.2` source-only pre-release was the correct
choice for an internal RC.

The full v1.0.0 asset list, build process, and upload/validation commands are
documented in [v1.0.0 Release Assets](release-assets.md).

## Current Recommendation

`v0.6.0-rc.2` RC validation is complete. No blocker bugs have been identified.
The remaining steps before tagging `v1.0.0` are:

- Finalize v1.0.0 release notes (see
  [Daryaft v1.0.0 Release Notes (draft)](release-notes-v1.0.0.md)).
- Re-verify all quality gates on the final release commit.
- Tag `v1.0.0`, build release artifacts with `goreleaser release`, and attach
  binary assets per [v1.0.0 Release Assets](release-assets.md).

v1.0.0 is a **stable baseline release** of the current feature set. No
additional product features are required. Post-1.0 features (Windows,
self-update, proxy/auth, package managers, concurrency, etc.) are deferred and
not blockers.

Public stable release is planned for `v1.0.0`.

See [Release-Candidate Validation](rc-validation.md) for the full validation
workflow. See [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md)
for the exact v1.0.0 criteria and the post-1.0 roadmap.
