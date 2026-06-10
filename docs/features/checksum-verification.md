# Checksum Verification

Daryaft can verify the integrity of a completed download against a known digest.
Pass `--checksum algorithm:hex` to enable verification after the file is written.

## Usage

```bash
daryaft download https://example.com/file.zip --checksum sha256:<hex>
daryaft https://example.com/file.zip --checksum sha256:<hex>
```

## Supported Algorithms

| Algorithm | Hex length | Notes |
|-----------|-----------|-------|
| `sha256`  | 64        | Recommended |
| `sha512`  | 128       | Supported |

## Dry-run

Dry-run validates the checksum format but does not download or compute the digest:

```bash
daryaft download https://example.com/file.zip --checksum sha256:<hex> --dry-run
```

Output includes the checksum spec in the plan:

```text
Checksum: sha256:<hex>
```

## Success output

On a successful match:

```text
Checksum verified: sha256
```

## Mismatch behavior

If verification fails, Daryaft exits non-zero and reports:

```text
checksum mismatch: expected <expected>, got <actual>
```

The downloaded file is left in place for inspection or removal. Daryaft does not
delete it automatically.

## Single-target limitation of `--checksum`

`--checksum` is for single-target downloads only. Providing `--checksum` with
multiple URLs or a `--file` batch list is rejected:

```text
--checksum is currently supported only for single URL downloads
```

For per-target checksums across multiple downloads, use `--checksum-file`
(below).

## Batch checksums with `--checksum-file`

`--checksum-file <path>` verifies every download in a batch against a manifest
file that maps one checksum to each target URL. It works with multiple URL
arguments, with `--file` URL lists, and with a single URL.

```bash
daryaft download URL1 URL2 --checksum-file checksums.txt
daryaft download --file urls.txt --checksum-file checksums.txt
daryaft --checksum-file checksums.txt URL
```

`--checksum` and `--checksum-file` cannot be used together:

```text
--checksum and --checksum-file cannot be used together
```

### Manifest format

One entry per line, in the form `<algorithm>:<hex> <url>`:

```text
# Comments and blank lines are ignored.
sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa https://example.com/file1.zip
sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb https://example.com/file2.zip
sha512:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc https://example.com/file3.tar.gz
```

Rules:

- Blank lines and lines beginning with `#` are ignored.
- Each remaining line must have exactly two whitespace-separated fields: the
  checksum spec and the URL.
- The same `sha256`/`sha512` algorithms and validation as `--checksum` apply.
- Every target URL must have exactly one checksum entry.
- Every manifest URL must match a planned target URL.
- Duplicate URLs in the manifest are rejected.
- Malformed lines are rejected and the error includes the line number.

### Exact URL matching

The URL in the checksum file must match the download target URL **exactly**.
Daryaft does not normalize, canonicalize, or re-encode URLs. A trailing slash,
different escaping, or different casing will not match.

Validation errors include:

```text
checksum file: no checksum entries found
checksum file: no checksum provided for URL: <url>
checksum file: manifest URL not in download targets: <url>
checksum file: manifest line 3: expected "<algorithm>:<hex> <url>" format
checksum file: manifest line 4: duplicate URL <url>
```

All `--checksum-file` validation happens before any network request.

### Batch verification behavior

- Each file is verified after it downloads successfully.
- A failed download is not verified.
- A checksum mismatch fails that item, counts toward the failed total, and the
  command exits non-zero. The downloaded file is left in place.
- The batch summary reports a `Checksum verified: N` count when at least one
  file passed verification.

## TUI checksum status

The TUI does not collect checksum input for batch downloads, and it does not
perform checksum verification itself. When a checksum-backed download runs, the
TUI displays the final checksum result from the execution model in the queue:

```text
✓ file.zip — Checksum OK
✗ file.zip — Checksum Failed
```

The final summary shows the `Checksum verified: N` count. Live "verifying"
progress is not shown in this release.

## Security note

Checksum verification confirms the downloaded file matches a known digest. It
does not prove who published the file. For publisher authenticity, digital
signatures or attestations would be needed, which are out of scope. Signature,
PGP, and attestation verification are not implemented.

Related docs:

- [Single URL Download](single-url-download.md)
- [Inspect and Dry-run](inspect-and-dry-run.md)
- [Checksum Verification QA](../operations/checksum-verification-qa.md)
