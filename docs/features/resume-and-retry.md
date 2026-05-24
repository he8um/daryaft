# Feature: Resume and Retry

Retry execution is implemented. Resume is still planned.

## Current Retry Behavior

`--retries` means retry attempts after the first try:

- `--retries 0`: one total attempt.
- `--retries 3`: one initial attempt plus up to three retries, for four total attempts.

The default is `3`.

Daryaft retries:

- network errors from the HTTP request
- timeouts
- HTTP `429`
- HTTP `500`
- HTTP `502`
- HTTP `503`
- HTTP `504`

Daryaft does not retry:

- HTTP `400`
- HTTP `401`
- HTTP `403`
- HTTP `404`
- existing final target files
- invalid output paths
- invalid URLs
- filename or path safety failures
- local filesystem permission errors

Retry backoff is exponential and capped:

- first failed attempt: wait `1s`
- second failed attempt: wait `2s`
- third failed attempt: wait `4s`
- later failures: wait up to `8s`

The CLI prints retry events:

```text
Retrying 2/4 in 1s: temporary server error: 503 Service Unavailable
```

Attempt numbering includes the initial attempt. `Retrying 2/4` means the next
attempt is attempt 2 of 4 total possible attempts.

## Batch Downloads

Sequential batch downloads use the same retry behavior for each item. One item
can retry and succeed or fail without changing the retry cycle for the next
item. If an item still fails after all retries, the batch continues and the
final summary reports that item failure.

## Resume Status

Resume is not implemented yet. Daryaft still writes to `<filename>.part`, but it
does not issue Range requests or continue from a previous partial file.

Until resume exists, every retry restarts the download and truncates the `.part`
file. Final files are never overwritten. If the final target file appears
between attempts, Daryaft fails that item without retrying again.

## Planned

Planned resume work includes Range requests, partial metadata validation,
checksum-aware completion, and safer recovery after interrupted downloads.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Batch Downloads](batch-downloads.md)
