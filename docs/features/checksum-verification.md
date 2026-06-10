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

## Batch limitation

`--checksum` is currently supported for single-target downloads only. Providing
`--checksum` with multiple URLs or a `--file` batch list is rejected:

```text
--checksum is currently supported only for single URL downloads
```

Per-file batch checksums are not supported in v1.7.0.

## TUI

TUI checksum entry is not implemented in v1.7.0. Using `--checksum` from the
root command routes to CLI download mode rather than the TUI.

## Security note

Checksum verification confirms the downloaded file matches a known digest. It
does not prove who published the file. For publisher authenticity, digital
signatures or attestations would be needed, which are out of scope for v1.7.0.

Related docs:

- [Single URL Download](single-url-download.md)
- [Inspect and Dry-run](inspect-and-dry-run.md)
- [Checksum Verification QA](../operations/checksum-verification-qa.md)
