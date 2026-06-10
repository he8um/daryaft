# Feature: Resume and Retry

Retry and single URL resume execution are implemented. Sequential batch
downloads inherit the same behavior for each item.

## Current Retry Behavior

`--retries` means retry attempts after the first try:

- `--retries 0`: one total attempt.
- `--retries 3`: one initial attempt plus up to three retries, for four total attempts.

The default is `3`.
Valid retry values are `0` through `20`.

Daryaft retries:

- network errors from the HTTP request
- timeouts
- interrupted response bodies, including unexpected EOF
- HTTP `408`
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
- HTTP `410`
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

HTTP `429 Too Many Requests` uses the same exponential backoff as the retryable
5xx responses. Daryaft does not yet honor a `Retry-After` response header; the
backoff schedule above always applies.

The CLI prints retry events:

```text
Retrying 2/4 in 1s: temporary server error: 503 Service Unavailable
```

Attempt numbering includes the initial attempt. `Retrying 2/4` means the next
attempt is attempt 2 of 4 total possible attempts.

The TUI execution screen consumes the same retry events and renders `Retrying`
status with the retry attempt, delay, and error message.

Cancellation is not retryable. When a download context is cancelled, Daryaft
emits a cancelled event and returns immediately without waiting for backoff or
starting another attempt.

In CLI mode, Ctrl+C or SIGTERM cancels through that context path. Daryaft keeps
the `.part` file and `.part.daryaft.json` metadata sidecar, does not rename the
partial file to the final target, prints a cancellation message, and exits with
code `1`. For batch downloads, the current item is cancelled and remaining URLs
are not started. TUI cancellation remains `q` from the progress screen.

Checksum verification (`--checksum` and `--checksum-file`) runs only after a
download completes successfully and reads the completed local file
synchronously. A cancelled download never reaches checksum verification, so a
cancelled item is always reported as cancelled rather than as a checksum or
generic failure.

## Batch Downloads

Sequential batch downloads use the same retry behavior for each item. One item
can retry and succeed or fail without changing the retry cycle for the next
item. If an item still fails after all retries, the batch continues and the
final summary reports that item failure.

## Resume Behavior

Daryaft writes incomplete downloads to:

```text
<filename>.part
<filename>.part.daryaft.json
```

The sidecar metadata stores the URL, target path, partial path, total bytes,
downloaded bytes, `ETag`, `Last-Modified`, `Accept-Ranges`, `created_at`, and
`updated_at`. Metadata is written through a temporary file and renamed into
place.

`--resume` defaults to `true`. When a `.part` file exists, Daryaft uses the
local partial size and sends:

```text
Range: bytes=<partial_size>-
```

If the server responds with `206 Partial Content`, Daryaft appends to the
partial file, progress starts at the existing byte count, and completion renames
the `.part` file to the final target. The sidecar metadata is removed after a
successful rename.

If the server does not support resume and returns a full response, Daryaft does
not append. It truncates the `.part` file, overwrites metadata, emits:

```text
Resume not supported by server; restarting download
```

and downloads from byte `0`.

If the server rejects the range request with `416 Requested Range Not
Satisfiable`, Daryaft always restarts the download from byte `0` as a safe
default. It does not treat `416` as proof that the partial file is already
complete.

If the local partial file is larger than the known remote total size, Daryaft
does not append past the end of the remote file. It restarts from byte `0` and
emits:

```text
Partial file is larger than remote file; restarting download
```

If the `.part` file exists but its `.part.daryaft.json` sidecar is missing or
corrupt, Daryaft does not fail. It treats the download as a fresh resume
candidate based on the partial file size when the server supports range
requests, and otherwise restarts safely. The final downloaded file is correct
in all of these cases.

The TUI execution screen renders resume and restart events as `Resuming` and
`Restarting` statuses with the same messages.

Cancelled downloads keep the `.part` file and `.part.daryaft.json` sidecar.
The next run with resume enabled can continue from the partial size when the
server supports HTTP Range requests.

If metadata contains an `ETag` or `Last-Modified` value and the resume response
returns a different value, Daryaft treats the remote file as changed, emits:

```text
Remote file changed; restarting download
```

and restarts from byte `0`.

`--no-resume` disables append behavior. Existing `.part` data is ignored for
resume, the partial file is truncated, metadata is overwritten, and the final
target overwrite rules still apply.

Retries work with resume. If an attempt writes some bytes and then fails, the
next retry can send a Range request for the current `.part` size. With
`--no-resume`, retries restart from byte `0`.

Final files are never overwritten. If the final target file exists before a
partial write, Daryaft fails that item without retrying again.

## Planned

Planned follow-up work includes checksum-aware completion and richer progress
rendering.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Batch Downloads](batch-downloads.md)
