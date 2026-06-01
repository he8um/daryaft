# Versioning Policy

Daryaft uses milestone-oriented semantic versioning.

The current development metadata version is `0.5.0-dev`. It reflects the
project's current pre-1.0 capability level and is not a public stable release.

## Rules

- `v0.x.0`: pre-1.0 local development / preview versions.
- `v1.0.0`: first official stable public installable release.
- `v1.0.x`: hotfixes.
- `v1.x.0`: bugfixes, polish, and small improvements.
- `v2.0.0` and later: feature-pack releases with 3 to 4 major new capabilities.

## Public Install Rule

Before `v1.0.0`, do not implement public install channels as active stable
channels. The repository can contain future packaging configuration, but
user-facing install instructions must clearly state that public install is
available from `v1.0.0` onward.

## Build Metadata

Source builds default to:

```text
version: 0.5.0-dev
commit: local
date: unknown
built_by: source
```

Release builds inject version metadata through linker flags into
`github.com/he8um/daryaft/pkg/version`. GoReleaser sets version, commit, date,
and `built_by`.

## Local Release Checks

`make release-check` runs GoReleaser in snapshot mode with publishing skipped:

```bash
goreleaser release --snapshot --clean --skip=publish
```

This check is local release-readiness validation only. It must not be treated as
a public release, must not create tags, and must not publish package-manager
artifacts before `v1.0.0`. Snapshot output is local and ignored by Git.

Related docs:

- [Roadmap](index.md)
- [Installation](../installation.md)
