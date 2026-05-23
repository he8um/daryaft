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

## `daryaft download [url...] --dry-run`

Implemented. Explicit form of the same download validation and dry-run planner.

```bash
daryaft download https://example.com/file.zip --dry-run
daryaft download -f urls.txt --dry-run
```

The current milestone does not perform network calls. Without `--dry-run`, the
command validates input and then reports that the downloader engine is not
implemented yet.

## Download Flags

- `-f`, `--file string`: read URLs from a file.
- `-o`, `--output string`: planned output directory.
- `--name string`: planned filename for a single URL.
- `--dry-run`: validate inputs and print the download plan.
- `--retries int`: planned retry count, default `3`.
- `--resume`: enable planned resume support, default `true`.
- `--no-resume`: disable planned resume support.

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

Real downloading for `daryaft <url>` and `daryaft download <url>` is planned for
the next downloader engine milestone.

Related docs:

- [Usage](usage.md)
- [Architecture Overview](architecture/overview.md)
