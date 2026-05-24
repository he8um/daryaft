# Command Reference

## `daryaft`

Implemented. Opens the Bubble Tea interactive home screen when run with no
arguments.

The home screen shows:

- Daryaft
- Modern terminal downloader
- Download from URL
- Download from .txt file
- View help
- Version
- Quit

Use up/down arrows or `k`/`j` to move, enter to select, `esc` or backspace to
return from sub-screens, and `q` or ctrl+c to quit. Download from URL and
Download from .txt file show planned screens and do not start downloads yet.
Existing CLI download commands remain the stable way to download files.

When URL arguments or `--file` are provided, the root command enters the current
download validation mode.

## `daryaft --help`

Implemented. Shows help text with:

- short description
- usage
- available commands
- common flags
- roadmap examples
- footer line

## `daryaft version`

Implemented. Prints:

- Daryaft version
- commit
- build date
- Go version

Default local values are used unless release tooling injects build variables.

## `daryaft [url...] --dry-run`

Implemented. Validates one or more HTTP/HTTPS URLs and prints a dry-run plan.

```bash
daryaft https://example.com/file.zip --dry-run
daryaft https://example.com/a.txt https://example.com/b.txt --dry-run
daryaft -f urls.txt --dry-run
```

## `daryaft <url>`

Implemented for exactly one URL. Performs an HTTP GET, accepts HTTP 2xx
responses, writes to `<filename>.part`, then renames the file after completion.
Incomplete downloads keep sidecar metadata at
`<filename>.part.daryaft.json`.

```bash
daryaft https://example.com/file.zip
daryaft https://example.com/file.zip --output downloads
daryaft https://example.com/file.zip --name file.zip
```

The command does not overwrite existing final files. It uses simple text output:

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

If the server does not provide a known content length, progress uses:

```text
Progress: <downloaded> bytes | <speed>
```

Progress lines are generated from structured downloader events. The current TUI
home screen does not start downloads yet; future TUI download screens will
consume the same event stream.

`--retries` is implemented for transient failures. The value is the number of
retry attempts after the first try, so `--retries 0` means one total attempt and
`--retries 3` means up to four total attempts. Daryaft retries network errors,
timeouts, interrupted response bodies, HTTP `429`, `500`, `502`, `503`, and
`504`.

`--resume` is enabled by default. If a partial file exists, Daryaft sends
`Range: bytes=<partial_size>-` and appends only after `206 Partial Content`.
If the server returns a full response instead, Daryaft truncates the `.part`
file and restarts safely. If saved `ETag` or `Last-Modified` metadata shows the
remote file changed, Daryaft also restarts from byte `0`. `--no-resume` ignores
existing partial data and overwrites the partial file from byte `0`.

## `daryaft download [url...] --dry-run`

Implemented. Explicit form of the same download validation and dry-run planner.

```bash
daryaft download https://example.com/file.zip --dry-run
daryaft download -f urls.txt --dry-run
```

## `daryaft download <url>`

Implemented. Explicit form of single URL real download.

## `daryaft [url...]`

Implemented for one or more URLs. When more than one URL is present, downloads
run sequentially in input order and continue after item failures.

```bash
daryaft https://example.com/a.txt https://example.com/b.txt
daryaft -f urls.txt
daryaft https://example.com/a.txt -f urls.txt
```

Batch output uses item headers, per-item progress, and a final summary:

```text
[1/2] Downloading: <url>
Saving to: <path>
Progress: <downloaded> / <total> bytes (<percent>%) | <speed>
Completed: <path>
Daryaft batch summary
Total: 2
Completed: 2
Failed: 0
```

If any item fails, Daryaft continues with the remaining URLs, lists failed
downloads in the summary, and returns a non-zero exit status at the end. Each
batch item has its own retry cycle.

## `daryaft download [url...]`

Implemented. Explicit form of sequential batch download.

## Download Flags

- `-f`, `--file string`: read URLs from a file.
- `-o`, `--output string`: output directory.
- `--name string`: filename for a single URL.
- `--dry-run`: validate inputs and print the download plan.
- `--retries int`: retry attempts after the initial attempt, default `3`.
- `--resume`: resume interrupted `.part` files with HTTP Range, default `true`.
- `--no-resume`: disable resume and restart partial files from byte `0`.

Validation rules:

- At least one URL arg or `--file` is required for download mode.
- URL args and `--file` are combined.
- Only `http` and `https` URLs are accepted.
- Empty lines and `#` comments are ignored in URL files.
- `--name` is rejected when more than one URL is provided.
- `--retries` must be greater than or equal to `0`.

## Common Flags

Implemented:

- `--no-color`: avoid color styling in the TUI.
- `--no-tui`: skip the no-argument TUI and print the non-interactive placeholder.
- `-v`, `--verbose`: enable verbose output when verbose logging exists.

## Planned Commands And Forms

These are planned and not implemented yet:

```bash
daryaft update
```

Batch concurrency, queue persistence, rich progress bars, TUI download
execution, segmented downloads, and self-update are planned.

Related docs:

- [Usage](usage.md)
- [Architecture Overview](architecture/overview.md)
