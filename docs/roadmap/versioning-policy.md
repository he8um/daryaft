# Versioning Policy

Daryaft uses milestone-oriented semantic versioning.

The latest stable release is `v1.11.0`. The current source development metadata
version is `1.12.0-dev`, representing active post-1.11.0 development on `main`.

## Version History Note

`v1.3.0` was tagged in source history but was intentionally not published as a
GitHub Release. The first GitHub Release after `v1.2.0` is `v1.4.0`. Do not
backfill or recreate `v1.3.0`. The `release-preflight` script enforces this
by detecting unexpected version gaps and requiring `--allow-skip` when a skip
is intentional.

## Rules

- `v0.x.0`: pre-1.0 local development / preview versions.
- `v1.0.0`: first official stable public installable release.
- `v1.0.x`: hotfixes.
- `v1.x.0`: bugfixes, polish, and small improvements.
- `v2.0.0` and later: feature-pack releases with 3 to 4 major new capabilities.

## Public Install Rule

`v1.0.0` is the first public stable release. GitHub release archives are the
supported install method. Package manager channels (Homebrew, deb, rpm, Arch)
are post-1.0 work.

## Build Metadata

Source builds default to:

```text
version: 1.12.0-dev
commit: local
date: unknown
built_by: source
```

Release builds inject version metadata through linker flags into
`github.com/he8um/daryaft/pkg/version`. GoReleaser sets version, commit, date,
and `built_by`.

## Release Preflight

Before tagging a release, always run:

```bash
make release-preflight VERSION=X.Y.Z
```

The preflight guardrail validates: clean tree, correct branch, source dev
version alignment, release notes existence, CHANGELOG entry, and absence of
the local/remote tag and GitHub release. It never creates tags or pushes.

If intentionally skipping a version, use:

```bash
make release-preflight-allow-skip VERSION=X.Y.Z
```

## Local Release Checks

`make release-check` runs GoReleaser v2 in snapshot mode with publishing
skipped:

```bash
goreleaser release --snapshot --clean --skip=publish
```

This check is local validation only. It does not publish a release, create tags,
or enable package-manager publishing. Snapshot versions are named
`1.12.0-dev-SNAPSHOT-<short-commit>` so they align with the current development
metadata instead of deriving from release tags. Snapshot output is local and
ignored by Git.

Related docs:

- [Roadmap](index.md)
- [Installation](../installation.md)
