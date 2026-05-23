# Feature: Batch Downloads

Batch input is implemented for URL file parsing, validation, and dry-run
planning. Real batch downloading is not implemented yet.

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

If a URL file contains exactly one effective URL, non-dry-run mode can download
that single URL:

```bash
daryaft -f one-url.txt
```

If a validated plan contains more than one URL and `--dry-run` is false, Daryaft
returns:

```text
batch downloading is not implemented yet; use --dry-run to inspect the plan
```

`--name` is rejected when multiple URLs are present because one filename cannot
safely apply to multiple downloads.

## Planned

Real batch downloads, concurrency, persistent queue state, history, retry
execution, and richer batch formats are planned.

Related docs:

- [Usage](../usage.md)
- [Command Reference](../command-reference.md)
- [Queue and History](queue-and-history.md)
