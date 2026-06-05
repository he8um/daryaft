# Daryaft Documentation

Daryaft is a modern terminal downloader written in Go. The project is currently
pre-1.0 and only the CLI foundation is implemented.

Public install channels begin at `v1.0.0`. Until then, use local development
commands only.

## Core Docs

- [Quick Start](quick-start.md): run the local development CLI.
- [Installation](installation.md): current local-only policy and future stable install plan.
- [Usage](usage.md): implemented commands and planned examples.
- [Command Reference](command-reference.md): current command behavior.
- [Configuration](configuration.md): default metadata and planned config locations.
- [Inspect and Dry Run](features/inspect-and-dry-run.md): URL metadata inspection and dry-run preflight behavior.
- [Architecture Overview](architecture/overview.md): planned high-level components.
- [Testing Strategy](engineering/testing-strategy.md): local and CI checks.
- [Manual QA Checklist](operations/manual-qa.md): local pre-release validation checklist.
- [Pre-Release Readiness](operations/pre-release-readiness.md): `0.6.0-dev` internal validation status and blockers.
- [QA Results: 0.6.0-dev](operations/qa-results-0.6.0-dev.md): completed internal validation readiness QA pass.
- [Release-Candidate Validation](operations/rc-validation.md): internal RC tag validation workflow.
- [RC Release Notes: v0.6.0-rc.2](operations/release-notes-v0.6.0-rc.2.md): current internal RC notes (unpublished).
- [RC Release Notes: v0.6.0-rc.1](operations/release-notes-v0.6.0-rc.1.md): superseded RC notes (historical reference).
- [Release Process](operations/release-process.md): local release readiness checks.
- [Branch Protection](operations/branch-protection.md): recommended `main` protection.

## Roadmap

- [Roadmap](roadmap/index.md): pre-1.0 milestones and post-1 feature packs.
- [Versioning Policy](roadmap/versioning-policy.md): release version rules.

## Related Existing Docs

The repository also contains deeper planning docs under `docs/features/`,
`docs/operations/`, `docs/troubleshooting/`, and `docs/engineering/`. Those
documents describe future work and should be kept aligned as features become
real.
