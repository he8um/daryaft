# Daryaft Documentation

Daryaft is a modern terminal downloader written in Go. **`v1.11.0` is the
current stable release.**

Install via Homebrew (`brew tap he8um/tap && brew install daryaft`) or download
binary archives from the
[GitHub releases page](https://github.com/he8um/daryaft/releases/tag/v1.11.0).
Other package manager channels are post-1.0 work.

## Core Docs

- [Quick Start](quick-start.md): run the local development CLI.
- [Installation](installation.md): Homebrew tap install, GitHub binary archives, and source build instructions.
- [Usage](usage.md): implemented commands and planned examples.
- [Command Reference](command-reference.md): current command behavior.
- [Configuration](configuration.md): default metadata and planned config locations.
- [Inspect and Dry Run](features/inspect-and-dry-run.md): URL metadata inspection and dry-run preflight behavior.
- [HTTP Request Customization](features/http-request-customization.md): proxy, custom headers, user-agent, and Basic Auth for download and inspect.
- [Checksum Verification](features/checksum-verification.md): single-target `--checksum` and batch `--checksum-file` manifest verification, plus TUI checksum status.
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
- [v1.11.0 Release Notes](operations/release-notes-v1.11.0.md): stable release notes for `v1.11.0` — read-only TUI Settings screen.
- [v1.10.0 Release Notes](operations/release-notes-v1.10.0.md): stable release notes for `v1.10.0` — user_agent/timeout config keys and --config/--timeout flags.
- [v1.9.0 Release Notes](operations/release-notes-v1.9.0.md): stable release notes for `v1.9.0` — download retry/resume reliability hardening.
- [v1.8.0 Release Notes](operations/release-notes-v1.8.0.md): stable release notes for `v1.8.0` — batch checksum verification and TUI checksum status.
- [v1.7.0 Release Notes](operations/release-notes-v1.7.0.md): stable release notes for `v1.7.0` — single-target checksum verification.
- [v1.1.0 Release Notes](operations/release-notes-v1.1.0.md): stable release notes for `v1.1.0` — read-only update check feature.
- [v1.4.0 Release Notes](operations/release-notes-v1.4.0.md): stable release notes for `v1.4.0` — reliability and test determinism.
- [v1.3.0 Release Notes](operations/release-notes-v1.3.0.md): stable release notes for `v1.3.0` — HTTP request customization.
- [v1.2.0 Release Notes](operations/release-notes-v1.2.0.md): stable release notes for `v1.2.0` — update check UX polish and install-channel hardening.
- [Update Check QA](operations/update-check-qa.md): manual QA checklist for `daryaft update --check`.
- [HTTP Customization QA](operations/http-customization-qa.md): manual QA checklist for HTTP request customization flags.
- [Checksum Verification QA](operations/checksum-verification-qa.md): manual QA checklist for `--checksum` and `--checksum-file` verification.
- [v1.2.0 Scope](roadmap/v1.2.0-update-ux.md): v1.2.0 update UX polish scope and quality gates.
- [Homebrew Tap](operations/homebrew-tap.md): live tap at `he8um/homebrew-tap`, install instructions, formula details, maintenance guide, and GoReleaser publishing checklist.
- [Homebrew Release Automation](operations/homebrew-release-automation.md): helper script for updating the tap formula after a release; dry-run support, safety rules, and future GoReleaser automation path.
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
- [v1.3.0 HTTP Customization Scope](roadmap/v1.3.0-http-customization.md): scope, quality gates, and future track for HTTP request customization.
- [v1.7.0 Checksum Verification Scope](roadmap/v1.7.0-checksum-verification.md): scope, quality gates, and limitations for single-target checksum verification.
- [v1.8.0 Batch Checksum + TUI Checksum Status Scope](roadmap/v1.8.0-batch-checksum-tui-checksum-ux.md): scope, data model, quality gates, and limitations for batch checksum support.
- [v1.9.0 Download Reliability Hardening Scope](roadmap/v1.9.0-download-reliability-hardening.md): scope, quality gates, and limitations for retry/resume reliability hardening.
- [v1.10.0 Config Persistence Safe Core Scope](roadmap/v1.10.0-config-persistence-safe-core.md): scope, quality gates, and limitations for user_agent/timeout config keys and --config/--timeout flags.
- [v1.11.0 TUI Config Awareness Scope](roadmap/v1.11.0-tui-config-awareness.md): scope, security policy, and quality gates for the read-only TUI Settings screen.

## Related Existing Docs

The repository also contains deeper planning docs under `docs/features/`,
`docs/operations/`, `docs/troubleshooting/`, and `docs/engineering/`. Those
documents describe future work and should be kept aligned as features become
real.
