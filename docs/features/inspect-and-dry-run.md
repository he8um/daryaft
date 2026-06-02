# Feature: Inspect and Dry Run

`daryaft inspect <url>` and download dry-runs are preflight tools. They help
check inputs and remote metadata before starting a real download.

## Inspect

```bash
daryaft inspect https://example.com/file.zip
daryaft inspect https://example.com/file.zip --json
```

Inspect accepts exactly one HTTP or HTTPS URL. It follows redirects and prints
metadata without saving a file. It does not create final files, `.part` files,
or `.part.daryaft.json` metadata sidecars.

Human output includes:

- original URL
- final URL after redirects
- HTTP status
- inferred filename
- content length when known
- content type when known
- `Accept-Ranges`
- resume support as `yes`, `no`, or `unknown`
- `ETag`
- `Last-Modified`

JSON output is stable for automation:

```json
{
  "url": "https://example.com/file.zip",
  "final_url": "https://cdn.example.com/file.zip",
  "status": "200 OK",
  "status_code": 200,
  "filename": "file.zip",
  "content_length": 1048576,
  "content_length_known": true,
  "content_type": "application/zip",
  "accept_ranges": "bytes",
  "resume_supported": true,
  "resume_support_known": true,
  "etag": "\"abc123\"",
  "last_modified": "Tue, 01 Jun 2026 12:00:00 GMT"
}
```

Daryaft tries `HEAD` first. If `HEAD` returns `405 Method Not Allowed` or omits
useful metadata, Daryaft may use a small `GET` request with
`Range: bytes=0-0`. Servers can still omit metadata, so unknown fields are
reported as `unknown` in human output and as empty values with `*_known: false`
where applicable in JSON.

Inspect currently does not implement checksum verification, proxy settings,
custom headers, auth, or a TUI flow.

## Dry Run

```bash
daryaft https://example.com/file.zip --dry-run
daryaft download https://example.com/file.zip --dry-run
daryaft -f urls.txt --dry-run
```

Dry-run validates URL and option parsing and prints the planned download
configuration. It does not make a network request and does not write files.

## Release Status

Daryaft is still pre-1.0. Public stable install channels remain planned for
`v1.0.0`; package-manager install channels are not stable before that release.
