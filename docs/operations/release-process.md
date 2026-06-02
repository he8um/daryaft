# Release Process

Daryaft is still pre-1.0. Public stable install channels and public release
publishing remain planned for `v1.0.0` and later.

## CI Validation

The test workflow includes a `goreleaser-check` job on push and pull request.
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
runs `golangci-lint run` with `.golangci.yml`. The `security` job installs and
runs `govulncheck ./...` and `gosec ./...`. These jobs do not publish releases,
create tags, or run snapshot builds.

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
make lint
make security
```

`make lint` requires `golangci-lint`. `make security` requires `govulncheck`
and `gosec`. These are the local equivalents of the CI `lint` and `security`
jobs.

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

Do not publish releases, create release tags, or enable Homebrew, deb, rpm, or
Arch package publishing before `v1.0.0`.
