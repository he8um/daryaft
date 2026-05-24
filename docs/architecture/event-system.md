# Event System

Downloader-to-interface event contracts and lifecycle.

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

## Current implementation

The downloader emits structured events from `internal/downloader`. The
implemented event types are:

- `started`: emitted after the target path and partial file are ready, before body copy starts.
- `progress`: emitted during the controlled copy loop, throttled to avoid excessive output.
- `completed`: emitted after the `.part` file is renamed to the final target path.
- `failed`: emitted when a download attempt fails after entering the downloader path.

Events carry simple fields that can be tested and reused by later interfaces:

- URL
- target path
- downloaded bytes
- total bytes, using `0` when the server does not provide a usable content length
- percent, only populated when total bytes are known and greater than zero
- approximate bytes per second
- error, for failed events
- timestamp

The CLI currently consumes these events for line-based progress output. Single
URL output uses:

```text
Downloading: <url>
Saving to: <path>
Progress: <downloaded> / <total> bytes (<percent>%) | <speed>
Completed: <path>
```

When the total size is unknown:

```text
Progress: <downloaded> bytes | <speed>
```

Sequential batch downloads reuse the same per-item downloader events. The batch
runner adds item metadata so the CLI can print headers and collect summary
counts:

```text
[1/3] Downloading: <url>
Saving to: <path>
Progress: <downloaded> / <total> bytes (<percent>%) | <speed>
Completed: <path>
```

The batch runner continues after item failures and returns a final result that
contains completed and failed item counts plus failure reasons.

## Boundaries

The event system is deliberately small. It does not implement Bubble Tea, a full
terminal UI, JSON output, concurrent batch download orchestration, queue
persistence, resume, retry, or segmented downloads yet.

Future TUI work should subscribe to the downloader event stream instead of
reaching into downloader internals. Future automation output can use the same
event data with a separate renderer.

## Examples

```bash
daryaft -h
daryaft https://example.com/file.zip
daryaft -f urls.txt
daryaft update --check
```
