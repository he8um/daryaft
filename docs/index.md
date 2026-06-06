# Daryaft Documentation

Daryaft is a modern terminal downloader written in Go. **`v1.0.0` is the first
stable release.** The source development version is now `1.1.0-dev`.

Install via Homebrew (`brew tap he8um/tap && brew install daryaft`) or download
binary archives from the
[GitHub releases page](https://github.com/he8um/daryaft/releases/tag/v1.0.0).
Other package manager channels are post-1.0 work.

## Core Docs

- [Quick Start](quick-start.md): run the local development CLI.
- [Installation](installation.md): Homebrew tap install, GitHub binary archives, and source build instructions.
- [Usage](usage.md): implemented commands and planned examples.
- [Command Reference](command-reference.md): current command behavior.
- [Configuration](configuration.md): default metadata and planned config locations.
- [Inspect and Dry Run](features/inspect-and-dry-run.md): URL metadata inspection and dry-run preflight behavior.
- [Architecture Overview](architecture/overview.md): planned high-level components.
- [Testing Strategy](engineering/testing-strategy.md): local and CI checks.
- [Manual QA Checklist](operations/manual-qa.md): local pre-release validation checklist.
- [Pre-Release Readiness](operations/pre-release-readiness.md): historical pre-1.0 readiness review (v1.0.0 has shipped).
- [QA Results: 0.6.0-dev](operations/qa-results-0.6.0-dev.md): completed internal validation readiness QA pass.
- [Release-Candidate Validation](operations/rc-validation.md): internal RC tag validation workflow.
- [RC Release Notes: v0.6.0-rc.2](operations/release-notes-v0.6.0-rc.2.md): current internal RC notes (GitHub pre-release published).
- [v0.6.0-rc.2 Release Status](operations/release-status-v0.6.0-rc.2.md): CI status, QA status, asset decision, and recommendation.
- [Clean Install Validation: v0.6.0-rc.2](operations/clean-install-validation-v0.6.0-rc.2.md): clean-clone/build/artifact validation pass — PASS WITH NOTES.
- [v1.0.0 Release Notes](operations/release-notes-v1.0.0.md): stable release notes with highlights, known limitations, install, and upgrade notes.
- [v1.0.0 Release Assets](operations/release-assets.md): binary asset strategy, GoReleaser build process, and upload/validation commands.
- [Homebrew Tap](operations/homebrew-tap.md): live tap at `he8um/homebrew-tap`, install instructions, formula details, maintenance guide, and GoReleaser publishing checklist.
- [v1.0.0 Go/No-Go Checklist](operations/v1.0.0-go-no-go.md): final pre-tag checklist — validated baseline, required checks, asset decision, go/no-go criteria.
- [v1.0.0 Release Plan](operations/v1.0.0-release-plan.md): step-by-step release execution plan including version policy, tagging, artifact build, publish, and post-release verification.
- [RC Release Notes: v0.6.0-rc.1](operations/release-notes-v0.6.0-rc.1.md): superseded RC notes (historical reference).
- [Release Process](operations/release-process.md): local release readiness checks.
- [Branch Protection](operations/branch-protection.md): recommended `main` protection.

## Roadmap

- [Roadmap](roadmap/index.md): pre-1.0 milestones and post-1 feature packs.
- [Release Readiness: v1.0](roadmap/release-readiness-v1.0.md): v1.0.0 criteria (stable baseline), required steps, and post-1.0 roadmap.
- [Post-1.0 Feature Packs](roadmap/post-1-feature-packs.md): features deferred until after stable baseline.
- [Versioning Policy](roadmap/versioning-policy.md): release version rules.
- [Self-Update Roadmap](roadmap/self-update.md): `update --check` current state and auto-update plan.

## Related Existing Docs

The repository also contains deeper planning docs under `docs/features/`,
`docs/operations/`, `docs/troubleshooting/`, and `docs/engineering/`. Those
documents describe future work and should be kept aligned as features become
real.
