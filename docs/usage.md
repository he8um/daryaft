# Usage

The current implementation is intentionally small.

## Implemented

```bash
daryaft --help
```

Shows the CLI description, usage, available commands, common flags, examples,
and footer line.

```bash
daryaft version
```

Prints the Daryaft version, commit, build date, and Go version.

```bash
daryaft
```

Opens the first Bubble Tea interactive home screen. The home screen includes
Download from URL, Download from .txt file, View help, Version, and Quit.
Download actions inside the TUI now open input forms. Entering a URL or a path
to a `.txt` URL file validates the input with the existing download planning
logic and shows a dry-run plan. Press enter on the plan screen to start a real
download. The TUI supports one URL and sequential `.txt` batch execution using
the same downloader event stream as the CLI. Cancellation is planned; while a
download is running, `q` shows that cancellation is not implemented and ctrl+c
can still terminate the program.

```bash
daryaft https://example.com/file.zip --dry-run
daryaft https://example.com/a.txt https://example.com/b.txt --dry-run
daryaft -f urls.txt --dry-run
daryaft download https://example.com/file.zip --dry-run
daryaft download -f urls.txt --dry-run
```

Validates input and prints a download plan. URL arguments and `--file` can be
combined. URL files are read line by line; empty lines and lines starting with
`#` are ignored.

Current flags:

- `-f`, `--file`: read URLs from a file.
- `-o`, `--output`: output directory.
- `--name`: filename for a single URL.
- `--dry-run`: print the plan without attempting a download.
- `--retries`: retry attempts after the initial attempt, default `3`.
- `--resume`: resume interrupted `.part` files, default `true`.
- `--no-resume`: ignore existing partial state and restart from byte `0`.

Common root flags:

- `--no-color`: avoid color styling in the TUI.
- `--no-tui`: skip the TUI and print the non-interactive placeholder.
- `-v`, `--verbose`: reserved for future verbose output.

```bash
daryaft https://example.com/file.zip
daryaft download https://example.com/file.zip
```

Downloads one HTTP/HTTPS URL with simple text output. Daryaft creates the output
directory when needed, writes to `<filename>.part`, then renames it to the final
filename when complete. While the partial file exists, Daryaft also writes
`<filename>.part.daryaft.json` metadata with the URL, target path, partial path,
byte counts, `ETag`, `Last-Modified`, `Accept-Ranges`, and timestamps. Existing
final files are not overwritten.

Resume is enabled by default. If `<filename>.part` exists, Daryaft checks the
local file size and sends `Range: bytes=<partial_size>-`. It appends only when
the server returns `206 Partial Content`. Progress starts at the existing byte
count:

```text
Resuming from 524288 bytes
Progress: 786432 / 1048576 bytes (75.0%) | 1.2 MB/s
```

If the server ignores Range and returns a full response, Daryaft safely
truncates the partial file and restarts:

```text
Resume not supported by server; restarting download
```

If saved `ETag` or `Last-Modified` metadata no longer matches the server
response, Daryaft does not append stale bytes:

```text
Remote file changed; restarting download
```

`--no-resume` ignores existing `.part` data for resume, truncates the partial
file, overwrites the sidecar metadata, and downloads from byte `0`. Existing
final target files are still rejected before Daryaft writes to the partial file.

During real single URL downloads, the CLI prints line-based progress from
structured downloader events:

```text
Downloading: https://example.com/file.zip
Saving to: downloads/file.zip
Progress: 524288 / 1048576 bytes (50.0%) | 1.2 MB/s
Completed: downloads/file.zip
```

Retry execution is implemented for transient network failures and temporary
server responses. `--retries 0` means one attempt total. `--retries 3` means the
initial attempt plus up to three retries, for four total attempts.

```text
Downloading: https://example.com/file.zip
Retrying 2/4 in 1s: temporary server error: 503 Service Unavailable
Saving to: downloads/file.zip
Progress: 524288 / 1048576 bytes (50.0%) | 1.2 MB/s
Completed: downloads/file.zip
```

Daryaft retries network errors, timeouts, HTTP `429`, `500`, `502`, `503`, and
`504`, plus interrupted response bodies such as unexpected EOF. When resume is
enabled, a retry after a partial body failure can continue from the current
`.part` size. With `--no-resume`, each retry restarts and truncates the `.part`
file. It does not retry client errors such as `404`, existing final files,
invalid output paths, filename safety failures, or local filesystem permission
errors.

When the server does not provide a known `Content-Length`, progress omits the
total and percent:

```text
Progress: 524288 bytes | 1.2 MB/s
```

Sequential batch downloads are also implemented:

```bash
daryaft https://example.com/a.txt https://example.com/b.txt
daryaft -f urls.txt
daryaft https://example.com/a.txt -f urls.txt
```

Batch downloads run one at a time in input order. URL arguments are processed
first, then URLs from `--file`. Each item uses the same filename selection,
partial file, existing final file rejection, and progress events as a single URL
download. Each item also gets its own retry cycle. `--name` remains invalid when
more than one URL is present.

Batch output includes an item header and final summary:

```text
[1/2] Downloading: https://example.com/a.txt
Saving to: downloads/a.txt
Progress: 1024 / 1024 bytes (100.0%) | 120.0 KB/s
Completed: downloads/a.txt
[2/2] Downloading: https://example.com/b.txt
Saving to: downloads/b.txt
Progress: 2048 / 2048 bytes (100.0%) | 200.0 KB/s
Completed: downloads/b.txt
Daryaft batch summary
Total: 2
Completed: 2
Failed: 0
```

By default, batch downloads continue after an item fails. If any item fails,
Daryaft prints failed URLs and returns a non-zero exit status after the summary.

Filenames are selected in this order:

1. `Content-Disposition` filename.
2. URL path base name.
3. `download.bin`.

For failed batches, the summary includes:

```text
Failed downloads:
- https://example.com/missing.txt: download "https://example.com/missing.txt" failed: server returned 404 Not Found
```

Sequential batch downloads inherit the same resume behavior per item. A failed
item can resume its own `.part` file during retries without affecting the next
item. If the final target file appears between attempts, Daryaft fails that item
instead of overwriting it.

## Planned Examples

These examples are roadmap examples and are not implemented yet:

```bash
daryaft update
```

Concurrency, queue persistence, TUI cancellation, rich progress bars, segmented
downloads, and self-update are planned. CLI download commands remain fully
supported alongside the TUI.

Related docs:

- [Command Reference](command-reference.md)
- [Configuration](configuration.md)
- [Roadmap](roadmap/index.md)
