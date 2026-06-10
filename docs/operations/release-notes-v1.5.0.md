# Daryaft v1.5.0

Daryaft v1.5.0 is a stable release focused on HTTP customization QA, security, and UX polish.

## Highlights

- Improved dry-run Basic Auth display to show the username while redacting the password.
- Added optional `DARYAFT_USERNAME` and `DARYAFT_PASSWORD` fallback for CLI Basic Auth.
- Improved invalid header and proxy error messages.
- Expanded sensitive-header redaction test coverage.
- Expanded HTTP customization QA documentation.

## HTTP credential fallback

`DARYAFT_USERNAME` and `DARYAFT_PASSWORD` can now be used as fallback Basic Auth credentials for CLI download and inspect flows.

CLI flags still take precedence:

- `--username`
- `--password`

The env fallback applies to:

- root URL download
- `daryaft download`
- `daryaft inspect`

It does not apply to TUI.

## Security

- Raw CLI passwords are not printed.
- Raw env passwords are not printed.
- Sensitive headers remain redacted.
- `DARYAFT_PASSWORD` without a username is rejected.
- CLI flags override env credentials.
- Basic Auth plus explicit `Authorization` header remains rejected.
- SOCKS5 proxy remains unsupported.

## Documentation

Updated:

- `docs/features/http-request-customization.md`
- `docs/operations/http-customization-qa.md`
- `CHANGELOG.md`

## Known limitations

- No OAuth support.
- No token manager support.
- No cookie jar/session persistence.
- No credential storage.
- No `.netrc` support.
- No TUI auth/config flow.
- No SOCKS5 proxy support.
- Auto-update is still not implemented.
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
