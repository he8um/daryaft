# Feature: Single URL Download

Single URL HTTP/HTTPS download is implemented with simple text progress output.

## Current

```bash
daryaft https://example.com/file.zip
daryaft download https://example.com/file.zip
daryaft https://example.com/file.zip --dry-run
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
- Emits structured downloader events for started, progress, resuming,
  restarting, retrying, completed, and failed states.
- Uses simple line-based text progress in the CLI.

Example output with a known total:

```text
Downloading: https://example.com/file.zip
Saving to: downloads/file.zip
Progress: 524288 / 1048576 bytes (50.0%) | 1.2 MB/s
Completed: downloads/file.zip
```

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

Filenames are sanitized so path traversal and directory separators cannot escape
the output directory.

## Planned

TUI rendering, rich progress bars, checksum validation, and segmented downloads
are planned. The event stream is in place as the foundation for the future TUI;
Bubble Tea is not integrated yet.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Batch Downloads](batch-downloads.md)
- [Downloader Engine](../architecture/downloader-engine.md)
