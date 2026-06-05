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
daryaft version --json
```

Prints the Daryaft version, commit, build date, build source, and Go version.
Use `daryaft version --json` for stable machine-readable metadata:

```json
{
  "version": "0.6.0-dev",
  "commit": "local",
  "date": "unknown",
  "built_by": "source",
  "go_version": "go1.xx.x"
}
```

Source builds default to `0.6.0-dev`, commit `local`, date `unknown`, and
built by `source`. Release builds inject metadata with linker flags.

```bash
daryaft inspect https://example.com/file.zip
daryaft inspect https://example.com/file.zip --json
```

Inspects one HTTP/HTTPS URL and prints metadata without saving a file. Inspect
follows redirects and reports the final URL, HTTP status, inferred filename,
content length when known, content type, `Accept-Ranges`, resume support,
`ETag`, and `Last-Modified`. Daryaft tries `HEAD` first and may fall back to a
small `Range: bytes=0-0` probe when the server does not support `HEAD` or omits
useful metadata. Metadata may be unknown if the server does not send the
corresponding headers. Use `--json` for stable machine-readable output.

```bash
daryaft doctor
daryaft doctor --json
daryaft doctor --strict
daryaft doctor --json --strict
```

Prints a local diagnostic report. Use `daryaft doctor` for human-readable text
and `daryaft doctor --json` for automation or CI. The report includes Go
runtime details, Daryaft version metadata, config path and config loading,
default download directory status, terminal environment values, optional
`clamscan` detection, and a skipped GitHub release check. Invalid config and
unwritable active output directories are critical failures. Missing configured
output directories are warnings and are not created by `doctor`. Both output
modes return non-zero only when critical failures exist, unless `--strict` is
set. Strict mode is useful in CI because warnings also cause a non-zero exit
status. JSON strict mode keeps warning checks as `warning` but sets `ok` to
`false`.

```bash
daryaft completion bash
daryaft completion zsh
daryaft completion fish
daryaft completion powershell
```

Generates shell completion scripts using Cobra's standard generators. These
commands do not require shell binaries to be present; they print the script to
stdout.

Example setup commands:

```bash
daryaft completion zsh > "${fpath[1]}/_daryaft"
daryaft completion bash > /etc/bash_completion.d/daryaft
daryaft completion fish > ~/.config/fish/completions/daryaft.fish
```

Installation paths vary by OS, shell, and user permissions.

```bash
daryaft config path
daryaft config show
daryaft config init
daryaft config init --force
daryaft config get retries
daryaft config set retries 5
daryaft config reset
daryaft config keys
```

Manages the YAML config file at `<UserConfigDir>/daryaft/config.yaml`.
`config path` prints the exact path, `config show` prints the effective config,
and `config init` creates the default file. `config init --force` overwrites an
existing config. `config get` prints one effective value, `config set` writes
one file value, `config reset` overwrites the file with defaults, and
`config keys` lists supported keys and types. CLI flags have priority over
environment variables, environment variables have priority over config file
values, and config file values have priority over built-in defaults.

Default config:

```yaml
download_dir: ""
retries: 3
resume: true
no_color: false
no_tui: false
theme: default
animations: true
hyperlinks: true
```

Supported set/get keys are `download_dir`, `retries`, `resume`, `no_color`,
`no_tui`, `theme`, `animations`, and `hyperlinks`. Shell completion suggests
these keys for `config get` and `config set`; boolean `config set` values also
suggest `true` and `false`.
`theme` supports `default` and `mono`; `mono` uses monochrome TUI styling.
`animations` and `hyperlinks` are reserved fields stored for future behavior
and do not currently change runtime output.

```bash
daryaft
```

Opens the first Bubble Tea interactive home screen. The home screen includes
Download from URL, Download from .txt file, Inspect URL, View help, Version,
and Quit.
Download actions inside the TUI now open input forms. Entering a URL or a path
to a `.txt` URL file validates the input with the existing download planning
logic, then opens an output directory input before showing the dry-run plan.
For single URL downloads, the TUI then prompts `Enter custom filename` and
`Enter checksum`; leaving those fields empty means auto-detect filename and
skip checksum verification. The `.txt` batch flow skips custom filename and
checksum input because one value cannot safely apply to multiple downloads.
Leaving the output directory empty means `.`, the current directory. Press
enter on the plan screen to start a real download. The TUI supports one URL and
sequential `.txt` batch execution using the same downloader event stream as the
CLI, and both flows honor the selected output directory. The TUI resizes its
panel and text inputs to the terminal window. Escape navigates back; Backspace
edits a non-empty text input and navigates back only when the current input is
empty.
Press `q` while a TUI download is running to cancel it. Cancelled downloads
keep the `.part` file and sidecar metadata for resume and are not retried. CLI
`-o`/`--output` and `--name` behavior is unchanged. CLI Ctrl+C/SIGTERM also
cancels through the downloader context, keeps the `.part` file and sidecar
metadata, stops remaining batch items, and exits non-zero. If `download_dir` is
set in config, the TUI output directory input starts with that value; otherwise
it starts with `.`.

The Inspect URL menu item prompts for one HTTP/HTTPS URL, runs the same
read-only metadata probe as `daryaft inspect <url>`, and shows the final URL,
status, inferred filename, content length, content type, resume support,
`ETag`, and `Last-Modified`. It does not download files, create `.part` files,
or write metadata sidecars. JSON inspect output remains a CLI-only
`daryaft inspect <url> --json` mode.

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
- `--checksum`: verify a completed single URL download with
  `sha256:<hex>` or `sha512:<hex>`.
- `--retries`: retry attempts after the initial attempt, default `3`, valid
  range `0` through `20`.
- `--resume`: resume interrupted `.part` files, default `true`.
- `--no-resume`: ignore existing partial state and restart from byte `0`.

If `-o`/`--output` is not provided and `download_dir` is set in config, Daryaft
uses that directory. `DARYAFT_DOWNLOAD_DIR` can override the config file value.
If `--retries` is not provided, Daryaft uses `DARYAFT_RETRIES` or config
`retries`. If neither `--resume` nor `--no-resume` is provided, Daryaft uses
`DARYAFT_RESUME` or config `resume`.

Common root flags:

- `--no-color`: avoid color styling in the TUI.
- `--no-tui`: skip the TUI and print the non-interactive placeholder.
- `-v`, `--verbose`: print extra CLI download diagnostics.

If `no_color` is true in config, the TUI uses no-color styling by default. If
`no_tui` is true in config, plain `daryaft` prints the non-interactive fallback
instead of launching the TUI. `DARYAFT_NO_COLOR` and `DARYAFT_NO_TUI` can
override those config file values. `DARYAFT_THEME` accepts `default` or `mono`.
Boolean environment values accept `true`, `1`, `yes`, `y`, `on`, `false`, `0`,
`no`, `n`, and `off`.

Environment examples:

```bash
DARYAFT_DOWNLOAD_DIR=~/Downloads daryaft https://example.com/file.zip
DARYAFT_RETRIES=5 daryaft https://example.com/file.zip
DARYAFT_NO_TUI=true daryaft
```

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

For manual integrity checks, pass `--checksum sha256:<hex>` or
`--checksum sha512:<hex>` with a single URL. Daryaft validates the checksum
format before downloading, computes the digest only after the final file has
completed and been renamed into place, and prints:

```text
Checksum verified: sha256
```

If the digest does not match, the command returns non-zero with an error like:

```text
checksum mismatch: expected <expected>, got <actual>
```

The completed final file is not deleted on mismatch in this milestone. Dry-run
validates and shows the checksum but does not compute it. Checksum verification
is currently single URL only. The TUI supports optional checksum input for
single URL downloads, while batch input and `--file` remain unsupported for
checksums.

With `--verbose` or `-v`, CLI downloads also print lines prefixed with
`Verbose:` for the effective URL with user info, query, and fragment redacted,
output directory, selected filename, checksum details when provided, HTTP
status when known, target path, resume/retry decisions, and completion
duration. Normal output is unchanged when verbose mode is not enabled.

Retry execution is implemented for transient network failures and temporary
server responses. `--retries 0` means one attempt total. `--retries 3` means the
initial attempt plus up to three retries, for four total attempts. Valid retry
values are `0` through `20`.

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

Concurrency, queue persistence, rich progress bars, segmented downloads, and
self-update are planned. CLI download commands remain fully supported alongside
the TUI.

Related docs:

- [Command Reference](command-reference.md)
- [Configuration](configuration.md)
- [Roadmap](roadmap/index.md)
