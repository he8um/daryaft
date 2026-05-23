# Contributing to Daryaft

Thank you for considering contributing to Daryaft.

## Development rules

1. Keep business logic out of `cmd/`.
2. Keep TUI rendering independent from downloader internals.
3. Use events to communicate downloader state.
4. Add tests for any downloader, updater, or parser behavior.
5. Update documentation when changing commands, flags, config, storage, or release behavior.

## Branching

Use focused branches:

```text
feature/downloader-core
feature/tui-progress
fix/resume-range-validation
chore/goreleaser-config
```

## Commit style

Prefer conventional commits:

```text
feat: add basic single URL download
fix: handle missing content length
chore: add release workflow
```

## Pre-1.0 note

Before v1.0.0, the project may be public but should not publish official user installation channels.
