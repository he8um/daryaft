# Daryaft v1.6.0

Daryaft v1.6.0 is a stable release focused on CLI download UX polish and TUI queue UX polish.

## Highlights

- Improved CLI download completion messages with final size and elapsed time.
- Added clearer in-band failure messages for single-download failures.
- Added clearer `Resuming:` and `Restarting:` CLI messages.
- Renamed unfinished batch-summary items from `Skipped` to `Not started`.
- Added TUI queue history during execution.
- Added TUI queue status markers for completed, failed/cancelled, in-progress, and other states.
- Improved TUI post-run wording and summary labels.

## CLI download UX

Single-download completion output now includes both final size and elapsed time, for example:

```text
Completed: file.zip (512 B in 1.2s)
```

Failure events now show an in-band message:

```text
Failed: <error>
```

Resume/restart events are easier to scan:

```text
Resuming: <message>
Restarting: <message>
```

Batch summaries now use:

```text
Not started
```

for items that never ran.

## TUI queue UX

During batch execution, the TUI now shows a queue history above the current item detail block.

Status markers:

| State              | Color mode | No-color mode |
| ------------------ | ---------: | ------------: |
| Completed          |        `✓` |        `[ok]` |
| Failed / Cancelled |        `✗` |         `[!]` |
| In progress        |        `→` |         `[>]` |
| Other              |        `·` |         `[-]` |

The post-run hint now says:

```text
enter/h new download • q quit
```

## Documentation

Updated:

- `CHANGELOG.md`
- `docs/usage.md`
- `docs/command-reference.md`
- `docs/operations/manual-qa.md`
- `docs/roadmap/v1.6.0-download-tui-ux-polish.md`

## Known limitations

- Auto-update is still not implemented.
- No OAuth support.
- No token manager support.
- No cookie jar/session persistence.
- No credential storage.
- No `.netrc` support.
- No TUI auth/config flow.
- No SOCKS5 proxy support.
- No config persistence.
- No checksum verification.
- No concurrent batch download.
- GoReleaser Homebrew publishing remains disabled.
- Windows is not officially supported.

## Upgrade

Homebrew:

```bash
brew update && brew upgrade daryaft
```

GitHub binary users:

Download the matching archive from GitHub Releases.

Source users:

```bash
git pull
go build .
```
