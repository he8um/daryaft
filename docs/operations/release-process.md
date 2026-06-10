# Release Process

**`v1.5.0` is the current stable release.** `main` is now on `1.6.0-dev`
post-v1.5.0 development.

> **Version skip note:** `v1.3.0` was tagged in source history but no GitHub
> Release was published for it. The tag exists locally and on the remote as
> part of normal commit history, but it was intentionally not promoted to a
> stable release. `v1.4.0` is the first release after `v1.2.0` that was
> published as a stable GitHub release with binary assets. Do **not** backfill
> or recreate a `v1.3.0` release.

See [Daryaft v1.5.0 Release Notes](release-notes-v1.5.0.md) for the latest
stable release notes and known limitations. Earlier releases:
- [v1.4.0 Release Notes](release-notes-v1.4.0.md): HTTP customization polish.
- [v1.3.0 Release Notes](release-notes-v1.3.0.md): HTTP request customization
  (tag exists; no GitHub release published — intentionally skipped).
- [v1.2.0 Release Notes](release-notes-v1.2.0.md): update check UX polish.
- [v1.1.0 Release Notes](release-notes-v1.1.0.md): read-only update check.
- [v1.0.0 Release Notes](release-notes-v1.0.0.md): initial stable baseline.

The binary asset strategy for v1.0.0 is in [v1.0.0 Release Assets](release-assets.md).
Historical RC and pre-release readiness docs remain available for reference:
- [Release-Candidate Validation](rc-validation.md)
- [Pre-Release Readiness](pre-release-readiness.md)
- [v0.6.0-rc.2 Release Status](release-status-v0.6.0-rc.2.md)
- [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md)

## Release Preflight Guardrail

Before creating a release tag, always run the release preflight script to
validate the target version against the current repository state:

```bash
make release-preflight VERSION=X.Y.Z
```

For example, to validate that `v1.5.0` is ready to tag:

```bash
make release-preflight VERSION=1.5.0
```

The preflight script checks:

1. Working tree is clean and on `main`, up to date with `origin/main`.
2. Source dev version in `pkg/version/version.go` matches `X.Y.Z-dev`.
3. Latest stable tag is the expected predecessor (no accidental version skip).
4. `docs/operations/release-notes-vX.Y.Z.md` exists and mentions the version.
5. `CHANGELOG.md` has a `[X.Y.Z]` section.
6. Local tag `vX.Y.Z` does not exist.
7. Remote tag `refs/tags/vX.Y.Z` does not exist.
8. GitHub release `vX.Y.Z` does not exist.

The script is read-only. It never creates tags, pushes, or modifies files.

### Intentional Version Skip

If a version is intentionally skipped (as happened with `v1.3.0`), use:

```bash
make release-preflight-allow-skip VERSION=X.Y.Z
```

This allows the skip but prints an explicit warning. Document the reason in
the release notes and CHANGELOG before tagging.

### Release Execution Order

```text
release notes finalized (docs/operations/release-notes-vX.Y.Z.md)
→ CHANGELOG.md [X.Y.Z] section ready
→ make release-preflight VERSION=X.Y.Z    ← must pass
→ quality gates (make rc-check)
→ git tag -a vX.Y.Z -m "Daryaft vX.Y.Z"
→ git push origin vX.Y.Z
→ wait for tag CI (green)
→ goreleaser release --clean --skip=publish
→ gh release create vX.Y.Z --notes-file ...
→ gh release upload vX.Y.Z dist/*.tar.gz dist/checksums.txt
→ post-release verification
→ Homebrew formula update
→ advance source version to next -dev
→ commit and push
```

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
1.5.0-dev-SNAPSHOT-<short-commit>
```

GoReleaser normally derives versions from Git tags. The snapshot template keeps
local dry-run metadata aligned with the current `1.5.0-dev` development version
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
version: 1.5.0-dev
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

When publishing a new Daryaft release, use the helper script to update the
Homebrew formula after GitHub release assets are published:

1. Build and publish the new GitHub release assets.
2. Clone the tap: `git clone https://github.com/he8um/homebrew-tap.git /tmp/homebrew-tap`
3. Run the helper: `scripts/update-homebrew-formula.sh --version X.Y.Z --tap-dir /tmp/homebrew-tap`
4. Review the diff and validate with `ruby -c` and `brew reinstall`.
5. In the tap clone: commit and push the formula update.
6. Verify with `brew update && brew upgrade daryaft && daryaft version`.

The script fetches checksums from the GitHub release, updates the formula
locally, runs a Ruby syntax check, and prints next steps. It never pushes or
commits automatically. See [Homebrew Release Automation](homebrew-release-automation.md)
for full details, dry-run usage, and safety rules.

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
