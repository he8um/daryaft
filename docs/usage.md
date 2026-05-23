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
daryaft -f urls.txt --dry-run
daryaft download https://example.com/file.zip --dry-run
daryaft download -f urls.txt --dry-run
```

Validates input and prints a download plan. URL arguments and `--file` can be
combined. URL files are read line by line; empty lines and lines starting with
`#` are ignored.

Current flags:

- `-f`, `--file`: read URLs from a file.
- `-o`, `--output`: planned output directory.
- `--name`: planned filename for a single URL.
- `--dry-run`: print the plan without attempting a download.
- `--retries`: planned retry count, default `3`.
- `--resume`: planned resume support, default `true`.
- `--no-resume`: disable planned resume support.

Non-dry-run download attempts validate input, then return:

```text
download engine is not implemented yet; use --dry-run to inspect the download plan
```

## Planned Examples

These examples are roadmap examples and are not implemented yet:

```bash
daryaft https://example.com/file.zip
daryaft -f urls.txt
daryaft update
```

Real downloading is planned for the next downloader engine milestone.

Related docs:

- [Command Reference](command-reference.md)
- [Configuration](configuration.md)
- [Roadmap](roadmap/index.md)
