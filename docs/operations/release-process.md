# Release Process

Daryaft is still pre-1.0. Public stable release is planned for `v1.0.0`.
v1.0.0 is a stable baseline release of the current implemented feature set;
it does not require additional product features before shipping. See
[Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md) for exact
v1.0.0 criteria and the post-1.0 roadmap.

See [Pre-Release Readiness](pre-release-readiness.md) for the current
`0.6.0-dev` internal validation verdict and remaining validation steps.
For internal RC tags, use
[Release-Candidate Validation](rc-validation.md).
The current RC is `v0.6.0-rc.2`; its release notes are in
[Daryaft v0.6.0-rc.2 Internal Release Candidate](release-notes-v0.6.0-rc.2.md)
and its GitHub pre-release status is in
[v0.6.0-rc.2 Release Status](release-status-v0.6.0-rc.2.md).

The v1.0.0 stable release notes draft is at
[Daryaft v1.0.0 Release Notes (draft)](release-notes-v1.0.0.md). The binary
asset strategy is documented in [v1.0.0 Release Assets](release-assets.md).

The final pre-tag go/no-go checklist is at
[v1.0.0 Go/No-Go Checklist](v1.0.0-go-no-go.md). The step-by-step release
execution plan is at [v1.0.0 Release Plan](v1.0.0-release-plan.md).

## CI Validation

The test workflow includes a `goreleaser-check` job on push and pull request.
The workflow's `push` trigger is not limited to branches, so tag pushes such as
`v0.6.0-rc.2` also run the workflow. This validates the RC tag without adding a
tag-triggered release or publishing job.
It installs GoReleaser v2 and runs:

```bash
goreleaser check
```

This validates `.goreleaser.yml` only. It does not run `goreleaser release`,
does not publish, does not create tags, and does not use publishing secrets.

The Go test/build workflow runs on both Linux and macOS. It checks module
tidiness, runs `go test ./...`, runs `go build ./...`, and runs
`go test -race ./internal/tui`.

The CI workflow also runs separate `lint` and `security` jobs. The `lint` job
runs blocking `golangci-lint run` with `.golangci.yml`. The `security` job
installs and runs `govulncheck ./...` and `gosec ./...`; both are blocking. The
previous Go 1.26.3 standard-library advisory gap is resolved by using Go 1.26.4
or newer, which the CI `security` job now does. These jobs do not publish
releases, create tags, or run snapshot builds.

Recommended branch protection checks:

- `Go test/build (ubuntu-latest)`
- `Go test/build (macos-latest)`
- `goreleaser-check`
- `lint`
- `security`

See [Branch Protection](branch-protection.md) for recommended `main` branch
settings.

## Local Snapshot Check

Before release-configuration work, run the normal local checks plus the local
quality gates when the tools are installed:

```bash
make ci
make rc-check
make rc-info
make lint
make security
```

`make lint` requires `golangci-lint`. `make rc-check` runs the release-candidate
readiness checks including `govulncheck` and `gosec`. `make security` requires
`govulncheck` and `gosec`, and runs both as blocking checks. `make lint` and
`make security` are the local equivalents of the CI `lint` and `security` jobs.
`make rc-info` prints the current Git describe value, local RC tags, source
version metadata, and reminders for `make rc-check` and `make release-check`.

Use the local release check before changing release configuration:

```bash
make release-check
```

The target verifies that GoReleaser v2 is installed. If it is missing, it
prints:

```text
GoReleaser is required. Install it with: brew install goreleaser
```

When GoReleaser is available, the target runs:

```bash
goreleaser release --snapshot --clean --skip=publish
```

This is local-only validation. It does not publish a GitHub release, does not
create tags, and does not enable package-manager publishing. Snapshot artifacts
are written under ignored local build directories such as `dist/`.

Snapshot versions are intentionally named:

```text
0.6.0-dev-SNAPSHOT-<short-commit>
```

GoReleaser normally derives versions from Git tags. The snapshot template keeps
local dry-run metadata aligned with Daryaft's current `0.6.0-dev` development
version without creating or deleting tags.

## Release Metadata

GoReleaser injects build metadata into `github.com/he8um/daryaft/pkg/version`
with linker flags:

- `Version`
- `Commit`
- `Date`
- `BuiltBy`

Source builds keep the default development metadata:

```text
version: 0.6.0-dev
commit: local
date: unknown
built_by: source
```

## Publishing Policy

Do not publish a stable release or enable package-manager publishing before
`v1.0.0`. The GitHub pre-release for `v0.6.0-rc.2` is published and marked
pre-release (not stable). Package-manager publishing (Homebrew, deb, rpm,
Arch) is post-1.0 work and not required for the v1.0.0 baseline release.

The v1.0.0 stable release will attach compiled binary assets (linux/amd64,
linux/arm64, darwin/amd64, darwin/arm64) and `checksums.txt`. Build with
`goreleaser release` (not `--snapshot`) on the `v1.0.0` tag. See
[v1.0.0 Release Assets](release-assets.md) for the full process.

## Post-1.0 Release Work

After a stable v1.0.0 baseline is shipped, the following release process work
can be addressed:

- Package manager publishing (Homebrew tap, apt/deb, rpm, Arch, Scoop).
- Automated binary asset upload in the release pipeline.
- Windows CI and verified binary builds for Windows.
- Self-update mechanism.
- Long-term compatibility and backport policy.

See [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md) for the
full post-1.0 roadmap.
