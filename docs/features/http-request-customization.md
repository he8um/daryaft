# HTTP Request Customization

Daryaft v1.3.0 adds HTTP request customization to the CLI download and inspect flows.

## Flags

| Flag | Description |
|------|-------------|
| `--proxy <url>` | HTTP or HTTPS proxy URL |
| `--header "Name: Value"` | Custom request header (repeatable) |
| `--user-agent <value>` | Override the default `User-Agent` header |
| `--username <value>` | HTTP Basic Auth username |
| `--password <value>` | HTTP Basic Auth password |

These flags apply to:
- `daryaft <url> ...` (root download mode)
- `daryaft download ...`
- `daryaft inspect ...`
- Single URL and batch downloads
- Dry-run output (with safe redaction)

These flags do **not** apply to:
- TUI download flows
- `daryaft config`, `daryaft update`, `daryaft doctor`, `daryaft version`

## Examples

### Custom header

```bash
daryaft download https://example.com/file.zip --header "X-Request-ID: demo"
```

### Multiple headers

```bash
daryaft download https://example.com/file.zip \
  --header "X-Request-ID: demo" \
  --header "Accept: application/octet-stream"
```

### User-Agent override

```bash
daryaft download https://example.com/file.zip --user-agent "Daryaft/1.3"
```

### Proxy

```bash
daryaft download https://example.com/file.zip --proxy http://127.0.0.1:8080
```

### Basic Auth

```bash
daryaft download https://example.com/private.zip \
  --username "$DARYAFT_USER" \
  --password "$DARYAFT_PASSWORD"
```

### Inspect with headers

```bash
daryaft inspect https://example.com/private.zip \
  --header "Authorization: Bearer $TOKEN"
```

### Dry-run with redaction

```bash
daryaft download https://example.com/private.zip \
  --username alice \
  --password secret \
  --dry-run
```

The dry-run output shows `[REDACTED]` for the password and sensitive header values:

```
Daryaft download plan
URLs: 1
1. https://example.com/private.zip
Output: current directory
Filename: auto-detect
Retries: 3
Resume: true
HTTP
  Auth: [REDACTED]
Mode: dry-run only, no network request performed
```

## Validation Rules

- **Proxy**: scheme must be `http` or `https`; SOCKS5 is not supported in v1.3.0.
- **Header**: must use `Name: Value` format; name must be a valid HTTP token; control characters rejected.
- **User-Agent**: `--user-agent` overrides any `User-Agent` custom header.
- **Basic Auth**: `--password` requires `--username`. `--username` with empty password is valid. Combining `--username/--password` with a custom `Authorization` header is rejected.

## Security Warnings

- Shell history may store `--password` values if passed directly. Prefer environment variables:
  ```bash
  --password "$DARYAFT_PASSWORD"
  ```
- No credential storage, OAuth, cookie jar, `.netrc`, or TUI auth is implemented in v1.3.0.
- Proxy URL credentials (if embedded in the URL) are not redacted from shell history. Avoid embedding credentials in proxy URLs.
- SOCKS5 proxy is not supported in v1.3.0.

## Limitations

- TUI download and inspect flows do not support these flags in v1.3.0.
- No per-URL header overrides in batch downloads; all URLs receive the same headers.
- No SOCKS5 proxy support.
- No credential storage or `.netrc` support.
