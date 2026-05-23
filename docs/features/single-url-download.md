# Feature: Single URL Download

Single URL download behavior is partly implemented as validation and dry-run
planning. Real HTTP downloading is planned for the next downloader engine
milestone.

## Current

```bash
daryaft https://example.com/file.zip --dry-run
daryaft download https://example.com/file.zip --dry-run
```

Current behavior:

- Accepts `http` and `https` URLs.
- Rejects empty or malformed URLs.
- Supports `--output`, `--name`, `--retries`, `--resume`, and `--no-resume`.
- Prints a dry-run plan.
- Does not make network calls.

Without `--dry-run`, Daryaft validates input and returns:

```text
download engine is not implemented yet; use --dry-run to inspect the download plan
```

## Planned

The downloader engine milestone will add real HTTP downloads, output filename
detection, filesystem writes, progress events, and error handling.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Downloader Engine](../architecture/downloader-engine.md)
