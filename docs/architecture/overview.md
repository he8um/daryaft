# Architecture Overview

Daryaft should stay small and testable. The current implementation has the CLI
foundation, version package, default metadata, validation and planning packages,
and a first single URL downloader engine.

## CLI Layer

`cmd/` owns Cobra commands, flags, help text, and command dispatch. It should not
contain downloader business logic.

## Downloader Engine

Partly implemented. The current engine performs one HTTP/HTTPS GET, accepts
HTTP 2xx responses, chooses a safe filename, writes to a `.part` file, and
renames it on success. Batch downloads, progress events, resume execution, retry
execution, checksum validation, and richer file conflict behavior are planned.

## Event System

Planned. Downloader state should be emitted as events so CLI output, TUI
rendering, JSON output, and tests can observe behavior without coupling to
engine internals.

## TUI Renderer

Planned. The renderer will subscribe to events and display progress, queue
state, errors, and history.

## Updater

Planned. The updater will check release metadata, download selected binaries,
verify checksums, and replace the executable safely.

## Config

Planned. Configuration should merge defaults, config files, environment
variables, and flags in a predictable order.

## Storage

Planned. Storage may persist history, queue state, partial download metadata,
and lockfiles.

Related docs:

- [Command Reference](../command-reference.md)
- [Configuration](../configuration.md)
- [Roadmap](../roadmap/index.md)
