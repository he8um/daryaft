# Daryaft v1.7.0

Daryaft v1.7.0 is a stable release focused on checksum verification for single-target CLI downloads.

## Highlights

- Added `--checksum algorithm:hex` for single-target CLI downloads.
- Added SHA-256 checksum verification.
- Added SHA-512 checksum verification.
- Added checksum validation before network requests.
- Added post-download checksum verification after successful downloads.
- Added clear checksum success and mismatch output.
- Added checksum dry-run display.
- Added root URL mode support for `--checksum`.
- Added clear validation for unsupported batch checksum usage.
- Added checksum feature documentation and manual QA.

## Usage

Verify a SHA-256 checksum:

```bash
daryaft download https://example.com/file.zip --checksum sha256:<hex>
```

Verify a SHA-512 checksum:

```bash
daryaft download https://example.com/file.zip --checksum sha512:<hex>
```

Root URL mode:

```bash
daryaft --checksum sha256:<hex> https://example.com/file.zip
```

Dry-run:

```bash
daryaft download https://example.com/file.zip --checksum sha256:<hex> --dry-run
```

## Behavior

On match:

```text
Checksum verified: sha256
```

On mismatch, Daryaft exits non-zero and reports:

* file path
* algorithm
* expected digest
* actual digest

The downloaded file is left in place for inspection/removal.

## Batch limitation

`--checksum` is currently supported for single-target downloads only.

These are intentionally rejected in v1.7.0:

```bash
daryaft download URL1 URL2 --checksum sha256:<hex>
daryaft download --file urls.txt --checksum sha256:<hex>
```

Per-file batch checksums are not supported yet.

## Security note

Checksum verification confirms file integrity against a known digest.

It does not prove who published the file. Publisher authenticity requires
signatures or attestations, which are out of scope for v1.7.0.

## Documentation

Updated or added:

* `docs/features/checksum-verification.md`
* `docs/operations/checksum-verification-qa.md`
* `docs/roadmap/v1.7.0-checksum-verification.md`
* `CHANGELOG.md`

## Known limitations

* Per-file batch checksums are not supported.
* TUI checksum entry is not implemented.
* Signature/PGP/attestation verification is not implemented.
* Checksum verification proves integrity against a known digest, not publisher authenticity.
* Auto-update is still not implemented.
* Config persistence is not implemented.
* Concurrent batch download is not implemented.
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
