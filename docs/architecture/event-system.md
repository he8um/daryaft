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
- `resuming`: emitted when Daryaft will append to an existing `.part` file after a valid `206 Partial Content` response.
- `restarting`: emitted when Daryaft discards partial state and restarts from byte `0`, such as unsupported Range or changed remote validators.
- `warning`: reserved for non-fatal downloader warnings that interfaces should show plainly.
- `retrying`: emitted before a retry delay when a failure is retryable and attempts remain.
- `completed`: emitted after the `.part` file is renamed to the final target path.
- `failed`: emitted when a download fails after retries are exhausted or the error is not retryable.
- `cancelled`: emitted when a context cancellation stops the download and partial state is preserved.

Events carry simple fields that can be tested and reused by later interfaces:

- URL
- target path
- partial path
- downloaded bytes
- total bytes, using `0` when the server does not provide a usable content length
- percent, only populated when total bytes are known and greater than zero
- approximate bytes per second
- message, for resuming, restarting, warning, and cancelled events
- error, for failed, retrying, and cancelled events
- attempt and max attempts, for retrying and attempt-scoped events
- next delay, for retrying events
- timestamp

The CLI consumes these events for line-based progress output. The TUI execution
screen consumes the same events through an injectable execution runner and a
goroutine/channel bridge, then renders status, target path, byte progress,
percent, speed, recent messages, completion, failure, and batch summaries. The
production TUI runner calls the existing batch downloader path; tests can
inject a runner to assert plans or cancellation without network downloads.
User-facing event behavior is unchanged. Single URL CLI output uses:

```text
Downloading: <url>
Saving to: <path>
Progress: <downloaded> / <total> bytes (<percent>%) | <speed>
Resuming from <bytes> bytes
Resume not supported by server; restarting download
Remote file changed; restarting download
Retrying <attempt>/<max> in <delay>: <reason>
Completed: <path>
```

When the total size is unknown:

```text
Progress: <downloaded> bytes | <speed>
```

Sequential batch downloads reuse the same per-item downloader events. The batch
runner adds item metadata so the CLI can print headers and the TUI can show
`Item N of M` while both interfaces collect summary counts:

```text
[1/3] Downloading: <url>
Saving to: <path>
Progress: <downloaded> / <total> bytes (<percent>%) | <speed>
Completed: <path>
```

The batch runner continues after item failures and returns a final result that
contains completed and failed item counts plus failure reasons.

## Boundaries

Cancellation emits `cancelled` with URL, target path, partial path, downloaded
bytes, total bytes, and the message `Download cancelled. Partial file kept for
resume.` Cancellation is not retryable and does not emit a failed event.

The event system is deliberately small. It does not implement JSON output,
concurrent batch download orchestration, queue persistence, or segmented
downloads yet.

The TUI subscribes to the downloader event stream instead of reaching into
downloader internals. Future automation output can use the same event data with
a separate renderer.

## Examples

```bash
daryaft -h
daryaft https://example.com/file.zip
daryaft -f urls.txt
# planned: daryaft update --check
```
