# Testing Strategy

Unit, integration, TUI, release, update, and package tests.

## Navigation

- [Documentation Index](../index.md)
- [Implementation Plan](../engineering/implementation-plan.md)
- [Agent Navigation](../engineering/agent-navigation.md)
- [Versioning Policy](../roadmap/versioning-policy.md)
- [Release Train](../roadmap/release-train.md)

## Project constants

```text
Project: Daryaft
Binary: daryaft
Module: github.com/he8um/daryaft
Repository: git@github.com:he8um/daryaft.git
Author: AmirHesam Piri <info@xhesam.com>
Website: https://xhesam.com
Project page: https://xhesam.com/daryaft
License: MIT
Footer: Developed with <3 by AmirHesam Piri
```

## Requirements

1. Implement this area using clean, isolated packages.
2. Keep command wiring in `cmd/`; do not put business logic there.
3. Use typed errors and user-safe messages.
4. Update `daryaft -h` help text when user-facing commands or flags change.
5. Update tests and documentation in the same change.
6. Do not commit private agent docs from `Documents/Daryaft-project/Docs` or `Documents/Daryaft-project/Caveman`.

## Implementation notes

The agent must treat this file as a contract. If behavior is ambiguous, prefer the behavior documented in:

- `../engineering/interfaces-and-contracts.md`
- `../engineering/error-model.md`
- `../architecture/module-boundaries.md`

## Acceptance criteria

- The feature is implemented in the correct module.
- The feature is covered by tests where practical.
- Errors are clear and actionable.
- The command help reflects the implemented behavior.
- The documentation includes examples and known limitations.

## Local CI

Run the local CI target before opening a PR:

```bash
make ci
```

This runs:

- `go mod tidy`
- `git diff --exit-code go.mod go.sum`
- `go test ./...`
- `go build ./...`
- `go test -race ./internal/tui`
- `git diff --check`
- `goreleaser check` when GoReleaser is installed

Downloader stabilization tests cover resume restart cases for `416`, full
`200` responses after Range requests, changed `ETag` and `Last-Modified`
metadata, context-aware retry backoff cancellation, HTTP response-header
timeouts, long body streaming without a total timeout, retry bounds, and output
directory traversal guards. Config tests cover atomic-style temp-file saves,
private config permissions, retry bounds, environment override validation, and
theme validation. TUI behavior tests cover responsive window sizing,
input-width adaptation, and Backspace editing versus empty-input navigation.
CLI tests cover verbose download diagnostics and unchanged non-verbose output.

If GoReleaser is not installed, `make ci` prints a warning and continues. The
strict release snapshot target remains `make release-check`, which requires
GoReleaser and does not publish.

## GitHub Actions

The test workflow runs on push and pull request. The Go test/build job runs on
Linux and macOS using the Go version declared in `go.mod`. It verifies module
tidiness, runs all Go tests, builds all packages, and runs the TUI race test.

The separate `goreleaser-check` job validates `.goreleaser.yml` with
`goreleaser check` only. It does not run snapshot builds, publish releases,
create tags, or use publishing secrets.

Workflow actions should stay on stable majors supported by GitHub-hosted
runners. The current workflow uses Node 24-compatible action majors for
checkout, Go setup, and GoReleaser validation.

Recommended branch protection checks:

- `Go test/build (ubuntu-latest)`
- `Go test/build (macos-latest)`
- `goreleaser-check`

See [Branch Protection](../operations/branch-protection.md) for the full
recommended `main` settings.

## Examples

```bash
daryaft -h
daryaft https://example.com/file.zip
daryaft -f urls.txt
daryaft update --check
```
