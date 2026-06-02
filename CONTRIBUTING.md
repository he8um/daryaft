# Contributing to Daryaft

Daryaft is in pre-1.0 development. Contributions should keep the codebase small,
testable, and honest about which features are implemented.

## Branches

Use focused branches named for the change:

```text
feature/downloader-core
feature/tui-progress
fix/resume-range-validation
chore/goreleaser-config
```

## Commit Messages

Prefer conventional commits:

```text
feat: add basic single URL download
fix: handle missing content length
chore: add release workflow
```

## Tests

Run these before opening a pull request:

```bash
go test ./...
go build ./...
make lint
```

Run `make security` before security-sensitive or downloader changes when
`govulncheck` and `gosec` are available locally. GitHub Actions runs separate
`lint` and `security` jobs on push and pull request, so local runs help catch
the same issues before CI.

Add tests for downloader, updater, parser, config, storage, and event behavior.
CLI-only text changes may use focused command tests once command behavior grows.

## Documentation

Update docs when changing commands, flags, config, storage, release behavior, or
the public roadmap. Do not document planned commands as active commands unless
they are implemented and tested.

## Project Rules

- Keep business logic out of `cmd/`.
- Keep TUI rendering independent from downloader internals.
- Use events to communicate downloader state when the engine is introduced.
- Keep public install instructions disabled until `v1.0.0`.
