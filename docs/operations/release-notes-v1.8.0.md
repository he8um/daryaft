# Daryaft v1.8.0

**Release date:** 2026-06-10
**Status:** Stable

Daryaft v1.8.0 adds batch checksum verification via a manifest file and TUI
checksum status display. It builds on v1.7.0 which added single-target
`--checksum` verification.

## Highlights

**`--checksum-file <path>`** — verify each download target against a manifest
of `<algorithm>:<hex> <url>` entries. Every target must have exactly one entry;
the command validates the manifest before any network request starts.

**TUI checksum status** — after a checksum-backed batch run, the queue shows
`Checksum OK` or `Checksum Failed` per item and the summary appends
`Checksum verified: N`.

## What's New

### `--checksum-file` flag

```bash
daryaft download URL1 URL2 --checksum-file checksums.txt
daryaft download --file urls.txt --checksum-file checksums.txt
```

Manifest format — one entry per line, blank lines and `#` comments ignored:

```
sha256:<hex> <url>
sha512:<hex> <url>
```

Validation rules (all checked before any network request):

- Every download target must have exactly one manifest entry.
- Every manifest entry must match a download target (no orphan entries).
- Duplicate URLs in the manifest are rejected.
- Malformed lines produce a line-numbered error.
- `--checksum` and `--checksum-file` cannot be used together.

### Checksum verification behavior

- On match: item marked `Checksum OK`; batch summary shows `Checksum verified: N`.
- On mismatch: item marked `Checksum Failed`, file kept in place, command exits non-zero.
- Verification runs after each successful download; failed downloads are not verified.

### TUI checksum status

After a checksum-backed batch completes:

- Queue items show `Checksum OK` (✓ in color mode, `[ok]` in no-color mode) or
  `Checksum Failed` (✗ / `[!]`).
- The execution summary appends `Checksum verified: N`.

## Install

**Homebrew:**

```bash
brew update && brew upgrade daryaft
daryaft version
```

**GitHub binary archives:**

```bash
# Example: macOS Apple Silicon
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.8.0/daryaft_darwin_arm64.tar.gz
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.8.0/checksums.txt
shasum -a 256 --check checksums.txt
tar -xzf daryaft_darwin_arm64.tar.gz
./daryaft version
```

Available archives: `daryaft_linux_amd64.tar.gz`, `daryaft_linux_arm64.tar.gz`,
`daryaft_darwin_amd64.tar.gz`, `daryaft_darwin_arm64.tar.gz`.

**Source:**

```bash
git pull
go build .
```

## Known Limitations

- Checksum file URL matching is exact; URLs are not normalized. The URL in the
  manifest must match the URL passed to the command character-for-character.
- GNU `sha256sum`-format files (digest-only lines without `<algorithm>:` prefix)
  are not supported.
- Checksum auto-discovery (fetching `.sha256` sidecar files) is not implemented.
- TUI batch checksum input forms are not implemented; the TUI displays results
  only (results come from CLI `--checksum-file`).
- Signature, PGP, and attestation verification are out of scope.
- A checksum mismatch does not delete the downloaded file.

## Upgrade Notes

- No breaking changes. `--checksum` single-target behavior is unchanged.
- `--checksum` and `--checksum-file` are mutually exclusive — passing both exits
  non-zero before any network request.
