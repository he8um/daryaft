# Architecture Overview

Daryaft should stay small and testable. The current implementation only has the
CLI foundation, version package, and default metadata.

## CLI Layer

`cmd/` owns Cobra commands, flags, help text, and command dispatch. It should not
contain downloader business logic.

## Downloader Engine

Planned. The engine will own HTTP requests, resume support, retry behavior,
checksum validation, and file writes.

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
