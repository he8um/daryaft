# Agent Navigation

This file tells coding agents which docs to read for each kind of task. It exists because dumping all docs into context is how token budgets go to die.

## Universal first reads

Always read these first:

1. `docs/engineering/project-structure.md`
2. `docs/architecture/module-boundaries.md`
3. `docs/engineering/interfaces-and-contracts.md`
4. `docs/engineering/error-model.md`
5. `docs/roadmap/release-train.md`

## Task routing

| Task | Read these docs |
|---|---|
| CLI command | `command-reference.md`, `engineering/coding-standards.md` |
| Downloader | `features/single-url-download.md`, `architecture/downloader-engine.md` |
| Batch input | `features/batch-downloads.md` |
| Resume/retry | `features/resume-and-retry.md`, `architecture/storage-and-state.md` |
| TUI | `features/terminal-ui.md`, `architecture/tui-architecture.md` |
| Interactive mode | `features/interactive-mode.md` |
| Update | `features/self-update.md`, `architecture/updater-architecture.md`, `roadmap/release-train.md` |
| Packaging | `features/packaging.md`, `engineering/release-process.md`, `installation.md` |
| Proxy | `features/proxy-and-networking.md`, `troubleshooting/proxy.md` |
| JSON automation | `features/automation-and-json.md`, `api/json-events.md` |
| Manifest | `features/manifest-and-lockfile.md`, `api/manifest-schema.md` |
| Security | `architecture/security-model.md`, `features/security-and-validation.md`, `SECURITY.md` |
| Tests | `engineering/testing-strategy.md`, `engineering/definition-of-done.md` |

## Write rules

- Write code only inside `Documents/Daryaft-project/Daryaft`.
- Do not copy private Docs/Caveman content into the repo unless it is meant to become public documentation.
- Keep docs linked from `docs/index.md`.
- When adding a feature, update command help and docs in the same commit.
