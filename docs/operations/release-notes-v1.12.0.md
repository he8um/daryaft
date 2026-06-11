# Daryaft v1.12.0

Daryaft v1.12.0 is a stable release focused on polishing the TUI download input experience.

## Highlights

- Added cleaner URL validation before TUI download planning.
- Added cleaner `.txt` file path validation before TUI download planning.
- Added actionable validation messages for empty, malformed, unsupported, missing, directory, and unreadable inputs.
- Added clearer `.txt` file path guidance with full-path examples.
- Added one-URL-per-line guidance for `.txt` download files.
- Added compact safe defaults preview on URL and file input screens.
- Improved visibility of configured output root / `download_dir`.
- Updated TUI help text to mention Settings and the `c` shortcut.
- Exported `download.ValidateURL` for shared URL validation.
- Added tests for TUI validation, defaults preview, no-color rendering, prompt/help text, and exported URL validation.
- Updated usage, command reference, manual QA, roadmap, index, and changelog docs.

## TUI input validation

The TUI now catches common input mistakes earlier and shows shorter, more actionable messages.

Examples include:

```text
URL cannot be empty.
URL must start with http:// or https://.
Invalid URL: host is required.
File path cannot be empty.
File not found: /path/to/urls.txt
Path is a directory: /path/to/dir
Cannot read file: /path/to/urls.txt
```

The core download planner remains authoritative for final validation.

## .txt file guidance

The TUI now clarifies that file-based downloads expect a full path to a `.txt` file containing URLs.

Example:

```text
/Users/amirhesampiri/Downloads/urls.txt
```

The `.txt` file should contain one URL per line.

## Defaults preview

The URL and file input screens now show a compact safe defaults preview before the download setup continues.

The preview includes:

```text
Output dir
Retries
Resume
User-Agent
Timeout
```

This helps users confirm the configured output root and common defaults before proceeding.

## Security model

The defaults preview intentionally shows only safe non-secret values.

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

This release does not add credential input or credential persistence.

## Documentation

Updated or added:

* `docs/usage.md`
* `docs/command-reference.md`
* `docs/operations/manual-qa.md`
* `docs/roadmap/v1.12.0-tui-download-input-ux.md`
* `docs/index.md`
* `CHANGELOG.md`

## Known limitations

* Queue/session persistence is not implemented.
* Config editing from TUI is not implemented.
* Credential/auth/proxy/header input from TUI is not implemented.
* File picker/autocomplete is not implemented.
* Mouse support is not implemented.
* `user_agent` and `timeout` are not newly wired into actual TUI download HTTP options unless already supported.
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
