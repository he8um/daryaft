# Doctor Diagnostics

`daryaft doctor` prints a local environment report in plain text.
`daryaft doctor --json` prints the same checks as machine-readable JSON for
automation and CI. It does not start the TUI and does not run network-heavy
checks.

## Checks

The report includes:

- Go runtime: OS, architecture, and Go version.
- Daryaft version metadata: version, commit, and build date.
- Config path and whether the config directory is writable or appears
  creatable.
- Effective config loading, including environment overrides.
- Default download directory. Empty `download_dir` means the current directory.
- Terminal hints: `TERM`, `NO_COLOR`, and stdout terminal status when available.
- Optional `clamscan` detection.
- GitHub release check status.

## Status Markers

```text
✓ ok
✗ critical failure
! warning
- informational
```

Critical failures return a non-zero exit code. These include a config path that
cannot be determined, invalid config YAML, an unwritable current directory when
`download_dir` is empty, and an existing configured `download_dir` that is not
writable.

Warnings do not fail the command. A configured download directory that does not
exist is a warning; `doctor` reports it but does not create it.

## JSON Output

Use JSON mode for scripts and CI:

```bash
daryaft doctor --json
```

The JSON structure is stable:

```json
{
  "ok": true,
  "summary": {
    "failures": 0,
    "warnings": 0,
    "checks": 12
  },
  "sections": [
    {
      "name": "System",
      "checks": [
        {
          "status": "ok",
          "label": "OS",
          "message": "darwin"
        }
      ]
    }
  ]
}
```

Status values are `ok`, `warning`, `failure`, `info`, and `skipped`. JSON mode
uses the same exit code behavior as the text report: non-zero for critical
failures, zero for warnings and informational checks.

## Optional Tools

`clamscan` is optional. If it is not found in `PATH`, `doctor` reports that as
informational only. The tool is reserved for future scan features and is not
required for current downloads.

## Skipped Checks

The GitHub release check is currently skipped. The doctor foundation does not
make a network request.
