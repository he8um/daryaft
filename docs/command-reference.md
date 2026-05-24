# Command Reference

## `daryaft`

Implemented. Prints the current placeholder message and exits successfully.
Interactive TUI mode is planned.

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
daryaft -f urls.txt --dry-run
```

## `daryaft <url>`

Implemented for exactly one URL. Performs an HTTP GET, accepts HTTP 2xx
responses, writes to `<filename>.part`, then renames the file after completion.

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
Completed: <path>
```

If the server does not provide a known content length, progress uses:

```text
Progress: <downloaded> bytes | <speed>
```

Progress lines are generated from structured downloader events. The future TUI
will consume the same event stream, but it is not implemented yet.

## `daryaft download [url...] --dry-run`

Implemented. Explicit form of the same download validation and dry-run planner.

```bash
daryaft download https://example.com/file.zip --dry-run
daryaft download -f urls.txt --dry-run
```

## `daryaft download <url>`

Implemented. Explicit form of single URL real download.

Batch real downloads are not implemented yet. A non-dry-run plan with more than
one URL returns:

```text
batch downloading is not implemented yet; use --dry-run to inspect the plan
```

## Download Flags

- `-f`, `--file string`: read URLs from a file.
- `-o`, `--output string`: output directory.
- `--name string`: filename for a single URL.
- `--dry-run`: validate inputs and print the download plan.
- `--retries int`: included in the plan, default `3`; retry execution is planned.
- `--resume`: included in the plan, default `true`; resume execution is planned.
- `--no-resume`: disable planned resume support in the plan.

Validation rules:

- At least one URL arg or `--file` is required for download mode.
- URL args and `--file` are combined.
- Only `http` and `https` URLs are accepted.
- Empty lines and `#` comments are ignored in URL files.
- `--name` is rejected when more than one URL is provided.
- `--retries` must be greater than or equal to `0`.

## Common Flags

Implemented harmless placeholders:

- `--no-color`: disable colored output when colorized output exists.
- `--no-tui`: disable terminal UI when the TUI exists.
- `-v`, `--verbose`: enable verbose output when verbose logging exists.

## Planned Commands And Forms

These are planned and not implemented yet:

```bash
daryaft update
```

Rich progress bars, TUI rendering, resume execution, retry execution, segmented
downloads, and self-update are planned.

Related docs:

- [Usage](usage.md)
- [Architecture Overview](architecture/overview.md)
