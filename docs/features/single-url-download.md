# Feature: Single URL Download

Single URL HTTP/HTTPS download is implemented with simple text progress output.

## Current

```bash
daryaft https://example.com/file.zip
daryaft download https://example.com/file.zip
daryaft https://example.com/file.zip --dry-run
daryaft https://example.com/file.zip --checksum sha256:<hex>
```

Current behavior:

- Accepts `http` and `https` URLs.
- Rejects empty or malformed URLs.
- Supports `--output`, `--name`, `--retries`, `--resume`, and `--no-resume`.
- Performs one HTTP GET for non-dry-run single URL plans.
- Accepts only HTTP 2xx responses.
- Creates the output directory when needed.
- Writes to `<filename>.part`, then renames to the final filename on success.
- Writes `<filename>.part.daryaft.json` sidecar metadata while incomplete.
- Resumes existing `.part` files with HTTP Range when `--resume` is enabled.
- Restarts safely from byte `0` when the server does not support Range.
- Restarts safely from byte `0` when saved `ETag` or `Last-Modified` metadata
  shows the remote file changed.
- Truncates `.part` files and overwrites metadata when `--no-resume` is used.
- Does not overwrite existing final files.
- Verifies manual CLI checksums after successful completed downloads when
  `--checksum sha256:<hex>` or `--checksum sha512:<hex>` is provided.
- Emits structured downloader events for started, progress, resuming,
  restarting, retrying, completed, and failed states.
- Uses simple line-based text progress in the CLI.
- Supports optional custom filename input in the TUI single URL flow.

Example output with a known total:

```text
Downloading: https://example.com/file.zip
Saving to: downloads/file.zip
Progress: 524288 / 1048576 bytes (50.0%) | 1.2 MB/s
Completed: downloads/file.zip
Checksum verified: sha256
```

Checksum verification is CLI-only in this milestone. Daryaft validates the
checksum format before starting the download. Dry-run shows the checksum but
does not compute it. A mismatch returns a non-zero error like `checksum
mismatch: expected <expected>, got <actual>` and leaves the completed final file
in place. Daryaft does not auto-discover checksum files, download checksum
files, verify signed checksums, or expose checksum entry in the TUI yet.

Example resume output:

```text
Downloading: https://example.com/file.zip
Resuming from 524288 bytes
Saving to: downloads/file.zip
Progress: 786432 / 1048576 bytes (75.0%) | 1.2 MB/s
Completed: downloads/file.zip
```

Example output when the server does not provide `Content-Length`:

```text
Progress: 524288 bytes | 1.2 MB/s
```

Filename selection:

1. `Content-Disposition` filename.
2. URL path base name.
3. `download.bin`.

The CLI `--name` option remains the command-line custom filename path for a
single URL. In the TUI, Download from URL follows URL input, output directory
input, custom filename input, then the plan screen. Leaving the TUI filename
field empty means auto-detect. A custom TUI filename is trimmed, lightly
validated, shown on the plan screen, and passed to the same download plan used
by execution.

Filenames are sanitized so path traversal and directory separators cannot escape
the output directory.

## Planned

Rich progress bars, TUI checksum flow, signed checksum handling, and segmented
downloads are planned. The Bubble Tea execution screen consumes the existing
downloader event stream.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Batch Downloads](batch-downloads.md)
- [Downloader Engine](../architecture/downloader-engine.md)
