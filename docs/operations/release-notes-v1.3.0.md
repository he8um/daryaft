# Daryaft v1.3.0

> **Status: RELEASED** — Stable release. Tag: `v1.3.0`.

## Summary

Daryaft `v1.3.0` adds HTTP request customization to the CLI download and inspect
flows. It ships `--proxy`, repeatable `--header "Name: Value"`, `--user-agent`,
and HTTP Basic Auth via `--username`/`--password`. All five flags apply to root
URL downloads, the `download` subcommand, `inspect`, and batch downloads. No new
external dependencies are introduced.

The release ships binary archives for linux/amd64, linux/arm64, darwin/amd64,
and darwin/arm64. The Homebrew formula at `he8um/homebrew-tap` has been updated
to `v1.3.0`.

## What's New

### HTTP Request Customization

Five new flags are available on `daryaft download` (and bare URL invocation) and
`daryaft inspect`:

| Flag | Description |
|------|-------------|
| `--proxy <url>` | HTTP or HTTPS proxy URL |
| `--header "Name: Value"` | Custom request header (repeatable) |
| `--user-agent <value>` | Override the default User-Agent |
| `--username <value>` | HTTP Basic Auth username |
| `--password <value>` | HTTP Basic Auth password |

#### Usage examples

```bash
# Proxy
daryaft download https://example.com/file.zip --proxy http://proxy.corp:8080

# Custom headers (repeatable)
daryaft download https://example.com/file.zip \
  --header "X-API-Token: mytoken" \
  --header "Accept: application/zip"

# User-Agent override
daryaft download https://example.com/file.zip --user-agent "MyApp/2.0"

# HTTP Basic Auth
daryaft download https://restricted.example.com/file.zip \
  --username alice --password s3cr3t

# Inspect with custom header
daryaft inspect https://api.example.com/artifact \
  --header "Authorization: Bearer mytoken"

# Dry-run preview (password shown as [REDACTED])
daryaft download https://example.com/file.zip \
  --username alice --password s3cr3t --dry-run
```

#### Batch download support

All flags apply uniformly to every item in a batch download:

```bash
daryaft download -f urls.txt --header "X-Batch: yes" --proxy http://proxy:8080
```

#### Proxy transport

The proxy wraps the existing transport, preserving all existing timeouts. Only
`http://` and `https://` proxy schemes are supported.

### Security and Redaction

- **Passwords are always redacted.** `--password` values never appear in output,
  logs, dry-run summaries, or verbose diagnostics. The raw value only reaches
  `req.SetBasicAuth()` internally.
- **Sensitive headers are redacted.** Header names matching `authorization`,
  `proxy-authorization`, `cookie`, `set-cookie`, `x-api-key`, `x-token`,
  `x-auth-token`, or names containing `secret`, `token`, `password`, `key`, or
  `authorization` are shown as `[REDACTED]` in dry-run and verbose output.

#### Dry-run example

```bash
$ daryaft download https://example.com/file.zip \
    --username alice --password topsecret \
    --header "Authorization: Bearer tok" \
    --header "X-Custom: visible" \
    --dry-run
# Output:
# HTTP
#   Headers:
#     Authorization: [REDACTED]
#     X-Custom: visible
#   Auth: [REDACTED]
```

### Validation Rules

- `--proxy`: `http://` or `https://` only. `socks5://` and other schemes are
  rejected with a clear error.
- `--header "Name: Value"`: colon separator required; empty values are allowed.
  Repeatable.
- `--password` without `--username`: rejected with error `--password requires --username`.
- `--username`/`--password` combined with an explicit `Authorization` header:
  rejected with error `use either --username/--password or an Authorization header, not both`.
- `--user-agent` overrides any `User-Agent` custom header.

### New Internal Package

`internal/httpopts` is a new shared package providing:

- Options types (`Header`, `Options`, `RedactedOptions`)
- Header parsing and validation
- Proxy URL validation
- Auth consistency checks
- Redaction (`Redact()`)
- Request application (`ApplyToRequest()`)
- Proxy client factory (`NewClient()`) — clones base transport, preserves timeouts

### TUI Isolation

HTTP customization flags are CLI-only in `v1.3.0`. All five flags are registered
in `hasDownloadFlagChanges()`, so passing any of them on the command line routes
to CLI download mode, not the TUI. The TUI is unaffected.

## Upgrade Instructions

### Homebrew

```bash
brew update
brew upgrade daryaft
daryaft version
daryaft download --help
```

### GitHub Binary Archives

Download `v1.3.0` assets from the
[v1.3.0 GitHub release](https://github.com/he8um/daryaft/releases/tag/v1.3.0):

```bash
# Example: macOS Apple Silicon
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.3.0/daryaft_darwin_arm64.tar.gz
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.3.0/checksums.txt
shasum -a 256 --check checksums.txt
tar -xzf daryaft_darwin_arm64.tar.gz
./daryaft version
```

Available archives:

| File | Platform |
|------|----------|
| `daryaft_linux_amd64.tar.gz` | Linux x86-64 |
| `daryaft_linux_arm64.tar.gz` | Linux ARM64 |
| `daryaft_darwin_amd64.tar.gz` | macOS Intel |
| `daryaft_darwin_arm64.tar.gz` | macOS Apple Silicon |
| `checksums.txt` | SHA-256 checksums for all archives |

### Source Build

```bash
git pull
go build .
./daryaft version
```

## Known Limitations

- **TUI HTTP customization not implemented.** `--proxy`, `--header`,
  `--user-agent`, `--username`, and `--password` are CLI-only. The interactive
  TUI does not support HTTP options in this release.
- **SOCKS5 proxy not supported.** `socks5://` scheme is rejected. SOCKS5
  requires `golang.org/x/net/proxy`; adding it would introduce an external
  dependency which is out of scope.
- **OAuth and token manager not implemented.** Bearer token support requires an
  explicit `--header "Authorization: Bearer <token>"`.
- **Cookie jar and session persistence not implemented.**
- **Credential storage not implemented.** Passwords are not persisted to keychain,
  config file, or environment. Use shell environment variables or OS-level
  credential managers externally if needed.
- **`.netrc` not implemented.**
- **Auto-update not implemented.** `daryaft update` without `--check` exits
  non-zero. A future release will implement in-place binary upgrade.
- **Windows not officially supported.** No Windows binaries are published.
- **GoReleaser Homebrew publishing remains disabled.** The `brews:` block in
  `.goreleaser.yml` is still commented out. The Homebrew formula is manually
  maintained, assisted by `scripts/update-homebrew-formula.sh`.
- **No signed checksums.** `checksums.txt` contains plain SHA-256 hashes.

## Validation Summary

All required checks passed before release:

| Check | Result |
|-------|--------|
| `go test ./...` | PASS |
| `go build ./...` | PASS |
| `go test -race ./internal/downloader` | PASS |
| `go test -race ./internal/tui` | PASS |
| `make lint` (`golangci-lint`) | 0 issues |
| `make security` (`govulncheck`) | no vulnerabilities |
| `make security` (`gosec`) | Issues: 0 |
| `goreleaser check` | 1 config validated |
| `make rc-check` | PASS |
| `git diff --check` | PASS |
| Invalid header rejection smoke | PASS |
| SOCKS5 rejection smoke | PASS |
| Password redaction smoke | PASS |
| GoReleaser artifact build | PASS |
| Extracted binary version | `1.3.0` |
| `built_by` | `goreleaser` |
| Checksum verification | PASS |
| `daryaft doctor` smoke | PASS |
| Homebrew formula update | PASS |
| `brew upgrade daryaft` | PASS |

## References

- [Changelog](../../CHANGELOG.md)
- [v1.3.0 Scope](../roadmap/v1.3.0-http-customization.md)
- [HTTP Request Customization](../features/http-request-customization.md)
- [HTTP Customization QA](http-customization-qa.md)
- [Manual QA Checklist](manual-qa.md)
- [Homebrew Tap](homebrew-tap.md)
- [v1.2.0 Release Notes](release-notes-v1.2.0.md)
- [Release Assets Strategy](release-assets.md)
- [Post-1.0 Feature Packs](../roadmap/post-1-feature-packs.md)
