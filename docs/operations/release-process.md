# Release Process

**`v1.0.0` is the first public stable release.** `main` is now on
`1.1.0-dev` post-release development.

See [Daryaft v1.0.0 Release Notes](release-notes-v1.0.0.md) for the stable
release notes and known limitations. The binary asset strategy for v1.0.0 is
in [v1.0.0 Release Assets](release-assets.md). Historical RC and pre-release
readiness docs remain available for reference:
- [Release-Candidate Validation](rc-validation.md)
- [Pre-Release Readiness](pre-release-readiness.md)
- [v0.6.0-rc.2 Release Status](release-status-v0.6.0-rc.2.md)
- [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md)

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
1.1.0-dev-SNAPSHOT-<short-commit>
```

GoReleaser normally derives versions from Git tags. The snapshot template keeps
local dry-run metadata aligned with the current `1.1.0-dev` development version
without creating or deleting tags.

## Release Metadata

GoReleaser injects build metadata into `github.com/he8um/daryaft/pkg/version`
with linker flags:

- `Version`
- `Commit`
- `Date`
- `BuiltBy`

Source builds keep the default development metadata:

```text
version: 1.1.0-dev
commit: local
date: unknown
built_by: source
```

## Publishing Policy

`v1.0.0` is the stable baseline release and has been published. It includes
binary assets for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, and
`checksums.txt`. `he8um/homebrew-tap` is live and is the first package-manager
install channel — the formula is manually maintained and GoReleaser `brews:`
publishing remains disabled. Other package-manager channels (deb, rpm, Arch)
are later post-1.0 work. See [v1.0.0 Release Assets](release-assets.md) and
[Homebrew Tap](homebrew-tap.md) for details.

### Future Release Process (Post-1.0)

When publishing a new Daryaft release, the Homebrew formula must be updated
manually until GoReleaser tap publishing is enabled:

1. Build and publish the new GitHub release assets.
2. Download the new `checksums.txt` from the GitHub release.
3. Update `Formula/daryaft.rb` in `he8um/homebrew-tap` with the new `version`,
   `url`, and `sha256` values.
4. Push the updated formula to `he8um/homebrew-tap`.
5. Verify with `brew update && brew upgrade daryaft && daryaft version`.

## Post-1.0 Release Work

After a stable v1.0.0 baseline is shipped, the following release process work
can be addressed:

- Enable GoReleaser `brews:` publishing to automate Homebrew formula updates.
- Package manager publishing for apt/deb, rpm, Arch, Scoop.
- Automated binary asset upload in the release pipeline.
- Windows CI and verified binary builds for Windows.
- Self-update mechanism.
- Long-term compatibility and backport policy.

See [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md) for the
full post-1.0 roadmap.
