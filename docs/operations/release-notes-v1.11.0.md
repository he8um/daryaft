# Daryaft v1.11.0

Daryaft v1.11.0 is a stable release focused on read-only TUI awareness of safe configuration settings.

## Highlights

- Added a read-only TUI Settings screen.
- Added a `Settings` item to the TUI home menu before `Quit`.
- Added `c` shortcut from the TUI home screen to open Settings.
- Added TUI display of active config path.
- Added TUI display of whether config was loaded from file or defaults are being used.
- Added TUI display of safe effective config values:
  - `download_dir`
  - `retries`
  - `resume`
  - `no_color`
  - `no_tui`
  - `theme`
  - `animations`
  - `hyperlinks`
  - `user_agent`
  - `timeout`
- Added display markers for default, unset, and reserved values.
- Added tests for Settings rendering, Settings navigation, `c` shortcut, no-color rendering, and secret exclusion.
- Updated TUI/config documentation and manual QA.

## TUI Settings screen

The TUI now includes a read-only Settings screen.

Open it from the TUI home screen by either:

```text
selecting Settings
pressing c
```

Return home with:

```text
esc
backspace
```

Quit with:

```text
q
```

## What it shows

The Settings screen shows the active configuration path and whether a config file was loaded.

It also shows safe effective settings such as:

```text
download_dir
retries
resume
no_color
no_tui
theme
animations
hyperlinks
user_agent
timeout
```

Empty/default values are presented clearly:

```text
(default)
(none)
(reserved)
```

## Security model

The Settings screen is read-only and intentionally excludes secrets.

It does not display:

* `DARYAFT_USERNAME`
* `DARYAFT_PASSWORD`
* usernames
* passwords
* tokens
* API keys
* Authorization headers
* Proxy-Authorization headers
* cookies
* proxy values
* arbitrary header values

It does not edit config files.

It does not provide config set/reset actions.

It does not introduce any TUI credential or auth flow.

## Documentation

Updated or added:

* `docs/configuration.md`
* `docs/command-reference.md`
* `docs/usage.md`
* `docs/operations/manual-qa.md`
* `docs/roadmap/v1.11.0-tui-config-awareness.md`
* `docs/index.md`
* `CHANGELOG.md`

## Known limitations

* Settings screen is read-only.
* TUI config editor is not implemented.
* Config set/reset from TUI is not implemented.
* Credential display is intentionally excluded.
* Auth flow in TUI is not implemented.
* Proxy/header persistence is not implemented.
* Config profiles are not implemented.
* Live config reload is not implemented.
* Auto-update is still not implemented.
* GoReleaser Homebrew publishing remains disabled.
* Windows is not officially supported.

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
