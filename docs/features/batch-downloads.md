# Feature: Batch Downloads

Batch input is implemented for URL file parsing, validation, dry-run planning,
and simple sequential real downloads.

## Current

```bash
daryaft -f urls.txt --dry-run
daryaft download -f urls.txt --dry-run
daryaft https://example.com/a.txt https://example.com/b.txt
daryaft -f urls.txt
daryaft https://example.com/a.txt -f urls.txt
```

URL files are read line by line:

- Empty lines are ignored.
- Lines starting with `#` are ignored.
- Whitespace is trimmed.
- URL arguments and `--file` input are combined.
- URL arguments are downloaded first, followed by URLs from `--file`.
- Only `http` and `https` URLs are accepted.

Real batch execution is sequential in this milestone. Daryaft downloads one URL
at a time, waits for it to complete or fail, then moves to the next URL.
The TUI uses the same sequential batch runner after a `.txt` file is validated
and confirmed from the plan screen. The TUI `.txt` batch flow is file input,
output directory input, then plan screen. It honors the selected output
directory, skips custom filename input, and keeps `Filename: auto-detect`
because one custom filename cannot safely apply to multiple downloads.

If a TUI batch is cancelled, Daryaft cancels the current item, keeps its partial
state for resume, and does not start remaining URLs. The final TUI summary
shows total, completed, failed, cancelled, and skipped counts when relevant.

Each item uses the normal single URL behavior:

- accepts only HTTP 2xx responses
- chooses `Content-Disposition`, URL basename, or `download.bin`
- writes to `<filename>.part`
- writes `<filename>.part.daryaft.json` metadata while incomplete
- resumes existing `.part` files with HTTP Range when `--resume` is enabled
- restarts safely from byte `0` when Range is unsupported or the remote changed
- restarts from byte `0` when `--no-resume` is used
- renames to the final path on success
- rejects existing final files
- retries transient network and server failures according to `--retries`
- emits downloader events that the CLI renders as text progress and the TUI
  renders as status/progress fields

Batch output starts each item with a clear header:

```text
[1/3] Downloading: https://example.com/a.txt
Saving to: downloads/a.txt
```

By default, batch downloads continue after failures. A failed item does not stop
the next URL from downloading. At the end, Daryaft prints:

```text
Daryaft batch summary
Total: 3
Completed: 2
Failed: 1

Failed downloads:
- https://example.com/missing.txt: download "https://example.com/missing.txt" failed: server returned 404 Not Found
```

If one or more items fail, Daryaft returns a non-zero exit status after printing
the summary.

Each batch item has its own retry cycle. `--retries 0` means one attempt for
each item. `--retries 3` means each item can make up to four attempts total.
Retryable failures are network errors, timeouts, interrupted response bodies,
HTTP `429`, `500`, `502`, `503`, and `504`.

Each batch item has independent resume state. If a retry follows a partial body
failure and `--resume` is enabled, that item can continue from its current
`.part` size. If a final target file exists or appears between attempts, that
item fails without retrying and the batch continues.

`--name` remains rejected when multiple URLs are present because one filename
cannot safely apply to multiple downloads. CLI `--name` behavior is unchanged
for single URL commands.

`--checksum` is also rejected for batch input and `--file` input in this
milestone because one checksum cannot safely apply to multiple files. Manual
checksum verification is currently single URL only in both CLI and TUI flows.

## Planned

Concurrency, persistent queue state, history, and richer batch formats are
planned. The current implementation deliberately does not do concurrent
downloads, queue persistence, or one custom filename for batch downloads.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Queue and History](queue-and-history.md)
