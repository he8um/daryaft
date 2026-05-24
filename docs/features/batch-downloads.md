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

Each item uses the normal single URL behavior:

- accepts only HTTP 2xx responses
- chooses `Content-Disposition`, URL basename, or `download.bin`
- writes to `<filename>.part`
- renames to the final path on success
- rejects existing final files
- emits downloader events that the CLI renders as text progress

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

`--name` remains rejected when multiple URLs are present because one filename
cannot safely apply to multiple downloads.

## Planned

Concurrency, persistent queue state, history, retry execution, TUI rendering,
and richer batch formats are planned. The current implementation deliberately
does not do concurrent downloads or queue persistence.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Queue and History](queue-and-history.md)
