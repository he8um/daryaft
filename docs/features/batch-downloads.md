# Feature: Batch Downloads

Batch input is partly implemented as URL file parsing and dry-run planning. Real
batch downloading is planned for the downloader engine and queue milestones.

## Current

```bash
daryaft -f urls.txt --dry-run
daryaft download -f urls.txt --dry-run
```

URL files are read line by line:

- Empty lines are ignored.
- Lines starting with `#` are ignored.
- Whitespace is trimmed.
- URL arguments and `--file` input are combined.
- Only `http` and `https` URLs are accepted.

`--name` is rejected when multiple URLs are present because one filename cannot
safely apply to multiple downloads.

## Planned

The downloader engine milestone will add actual downloads. Later queue work will
add concurrency, persistent queue state, history, and richer batch formats.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Queue and History](queue-and-history.md)
