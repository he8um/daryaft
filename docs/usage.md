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

Prints a placeholder explaining that interactive mode is planned for the TUI
milestone.

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
- `--retries`: planned retry count, default `3`.
- `--resume`: planned resume support, default `true`.
- `--no-resume`: disable planned resume support.

```bash
daryaft https://example.com/file.zip
daryaft download https://example.com/file.zip
```

Downloads one HTTP/HTTPS URL with simple text output. Daryaft creates the output
directory when needed, writes to `<filename>.part`, then renames it to the final
filename when complete. Existing final files are not overwritten. Because resume
is not implemented yet, an existing `.part` file is restarted/truncated.

During real single URL downloads, the CLI prints line-based progress from
structured downloader events:

```text
Downloading: https://example.com/file.zip
Saving to: downloads/file.zip
Progress: 524288 / 1048576 bytes (50.0%) | 1.2 MB/s
Completed: downloads/file.zip
```

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
download. `--name` remains invalid when more than one URL is present.

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

## Planned Examples

These examples are roadmap examples and are not implemented yet:

```bash
daryaft update
```

Concurrency, queue persistence, TUI, rich progress bars, resume, retry
execution, segmented downloads, and self-update are planned. The current
downloader event stream is the foundation for the future TUI, but no Bubble Tea
interface is implemented yet.

Related docs:

- [Command Reference](command-reference.md)
- [Configuration](configuration.md)
- [Roadmap](roadmap/index.md)
