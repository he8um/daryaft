# Feature: Terminal UI

The Bubble Tea terminal UI is implemented for no-argument startup. It can
collect URL or `.txt` file input, show dry-run download plans, and start real
downloads from the plan screen. It also has a read-only Inspect URL flow for
remote metadata checks.

## Navigation

- [Documentation Index](../index.md)
- [Implementation Plan](../engineering/implementation-plan.md)
- [Agent Navigation](../engineering/agent-navigation.md)
- [Versioning Policy](../roadmap/versioning-policy.md)
- [Release Train](../roadmap/release-train.md)

## Project constants

```text
Project: Daryaft
Binary: daryaft
Module: github.com/he8um/daryaft
Repository: git@github.com:he8um/daryaft.git
Author: AmirHesam Piri <info@xhesam.com>
Website: https://xhesam.com
Project page: https://xhesam.com/daryaft
License: MIT
Footer: Developed with <3 by AmirHesam Piri
```

## Current Behavior

Running `daryaft` with no URL arguments opens an interactive home screen. It is
implemented with Bubble Tea and styled with Lip Gloss.

The home screen shows:

- app title: `Daryaft`
- subtitle: `Modern terminal downloader`
- menu items:
  1. Download from URL
  2. Download from .txt file
  3. Inspect URL
  4. View help
  5. Version
  6. Quit
- footer: `Developed with <3 by AmirHesam Piri`

Navigation:

- up/down arrows or `k`/`j` move through menu items
- enter selects an item
- `esc` returns from sub-screens
- backspace edits non-empty text inputs and returns only when the current input
  is empty
- `q` quits unless a download is running
- ctrl+c quits from anywhere

The View help and Version menu items render simple in-TUI screens. Download
from URL and Download from .txt file render text input forms. Pressing enter
validates the URL or URL file through the existing download planner, then opens
an output directory input screen. The default/current value is config
`download_dir` when set, otherwise `.`. The single URL flow then opens an
`Enter custom filename` input screen with `Leave empty to auto-detect` help
text. Empty filename input means auto-detect; a custom filename is shown on the
plan and passed to the existing download plan. The `.txt` batch flow skips this
filename screen and keeps `Filename: auto-detect`. Pressing enter from the last
input shows a dry-run plan with URL count, the first URLs, selected output
directory, filename, retries, and resume settings. TUI retries and resume use
config defaults unless the code path receives explicit values from the CLI.

Inspect URL prompts `Enter URL to inspect`, validates HTTP/HTTPS input, then
runs the shared `internal/inspect` probe through an injectable runner. It shows
the final URL, status, inferred filename, content length, content type,
`Accept-Ranges`, resume support, `ETag`, and `Last-Modified`. This flow is
read-only: it does not start downloader execution, create `.part` files, or
write metadata sidecars. JSON inspect output remains available through the CLI
`daryaft inspect <url> --json` command.

The TUI handles terminal resize messages. Its bordered panel uses available
terminal width with sensible lower and upper bounds, and text inputs resize
with the panel so narrow terminals remain usable while wide terminals stay
readable.

On the plan screen, enter starts the real download. The execution screen shows
the current item, URL, target path when known, status, downloaded bytes, total
bytes when known, percent, speed, recent messages, and a final batch summary.
It consumes the same downloader event stream as the CLI and supports sequential
batch execution for `.txt` input. TUI downloads honor the selected output
directory for both single URL and `.txt` batch flows. CLI `-o`/`--output` and
`--name` behavior is unchanged.

While a download is running, `q` cancels it and shows `Cancelling...`. After the
downloader stops, the status becomes `Cancelled` with the message
`Download cancelled. Partial file kept for resume.` Cancelled downloads keep
their `.part` file and metadata sidecar and are not retried.

CLI commands remain fully supported:

```bash
daryaft https://example.com/file.zip
daryaft -f urls.txt
daryaft download https://example.com/file.zip
```

`--no-color` keeps the panel layout but avoids color styling. Config `no_color`
can provide the same default for plain `daryaft`. Config `theme: mono` also
uses monochrome styling. `theme: default` is the normal styled mode.
`animations` and `hyperlinks` are reserved config fields and do not currently
change TUI behavior. `--no-tui` skips the no-argument TUI and prints the
non-interactive placeholder; config `no_tui` can also make that the default.

Existing CLI download commands are unchanged and remain the stable way to
download files.

## Planned

Queue persistence, concurrent downloads, history, and rich progress bars are
not implemented yet.

## Examples

```bash
daryaft -h
daryaft
daryaft https://example.com/file.zip
daryaft -f urls.txt
```
