# Implementation Plan

This plan is optimized for coding agents.

## Read order for agents

1. [Agent Navigation](agent-navigation.md)
2. [Project Structure](project-structure.md)
3. [Module Boundaries](../architecture/module-boundaries.md)
4. [Interfaces and Contracts](interfaces-and-contracts.md)
5. [Error Model](error-model.md)
6. [Versioning Policy](../roadmap/versioning-policy.md)
7. [Release Train](../roadmap/release-train.md)
8. Feature-specific document for the current task

## Phase 1: Repository skeleton

Implement:

- `go.mod` with module `github.com/he8um/daryaft`
- `main.go`
- `cmd/root.go`
- `pkg/version/version.go`
- Makefile targets
- `.gitignore`
- docs skeleton

Acceptance:

- `go test ./...` passes.
- `go run . -h` works.
- No private Docs/Caveman files are inside repo.

## Phase 2: Downloader core

Read:

- [Single URL Download](../features/single-url-download.md)
- [Downloader Engine](../architecture/downloader-engine.md)
- [Error Model](error-model.md)

Implement:

- Download task model
- HTTP GET download
- Progress events
- Filename detection
- Output directory handling

## Phase 3: TUI foundation

Read:

- [Terminal UI](../features/terminal-ui.md)
- [TUI Architecture](../architecture/tui-architecture.md)
- [Interactive Mode](../features/interactive-mode.md)

Implement:

- Bubble Tea model
- Progress bar
- Spinner
- Footer with hyperlink fallback
- Interactive home screen for bare `daryaft`

## Phase 4: Batch, resume, retry

Read:

- [Batch Downloads](../features/batch-downloads.md)
- [Resume and Retry](../features/resume-and-retry.md)
- [Storage and State](../architecture/storage-and-state.md)

Implement:

- txt URL parser
- `.part` file
- HTTP Range resume
- retry with backoff

## Phase 5: Update system

Read:

- [Self Update](../features/self-update.md)
- [Updater Architecture](../architecture/updater-architecture.md)
- [Release Train](../roadmap/release-train.md)

Implement:

- GitHub release lookup
- semver comparison
- asset selection
- update size and changelog preview
- y/N confirmation
- checksum verification
- binary replacement and rollback

## Phase 6: Packaging and docs gate

Read:

- [Packaging](../features/packaging.md)
- [Release Process](release-process.md)
- [Installation](../installation.md)

Implement:

- GoReleaser config
- package build validation
- pre-1.0 publishing guard
- full `-h` help
- docs audit

## Phase 7: v1.0.0 readiness

Read:

- [Definition of Done](definition-of-done.md)
- [Pre-1 Roadmap](../roadmap/pre-1-roadmap.md)

Ship only when all v1.0.0 acceptance criteria pass.
