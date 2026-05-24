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
- Leaves `.part` files in place on failure for future resume work.
- Restarts/truncates existing `.part` files because resume is not implemented yet.
- Does not overwrite existing final files.
- Emits structured downloader events for started, progress, completed, and failed states.
- Uses simple line-based text progress in the CLI.

Example output with a known total:

```text
Downloading: https://example.com/file.zip
Saving to: downloads/file.zip
Progress: 524288 / 1048576 bytes (50.0%) | 1.2 MB/s
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

TUI rendering, rich progress bars, resume execution, retry execution, checksum
validation, and segmented downloads are planned. The event stream is in place as
the foundation for the future TUI; Bubble Tea is not integrated yet.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Downloader Engine](../architecture/downloader-engine.md)
