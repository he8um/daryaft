# Checksum Verification QA

Manual QA checklist for `--checksum` verification behavior.

## Setup

Start a local server with a known file:

```bash
tmpdir="$(mktemp -d)"
printf 'hello daryaft checksum\n' > "$tmpdir/file.txt"
expected_sha256="$(shasum -a 256 "$tmpdir/file.txt" | awk '{print $1}')"
expected_sha512="$(shasum -a 512 "$tmpdir/file.txt" | awk '{print $1}')"
echo "SHA-256: $expected_sha256"
echo "SHA-512: $expected_sha512"
cd "$tmpdir"
python3 -m http.server 18171 --bind 127.0.0.1 &
SERVER_PID=$!
```

## Matching SHA-256

```bash
go run . download http://127.0.0.1:18171/file.txt \
  --output "$tmpdir/out-sha256.txt" \
  --checksum "sha256:$expected_sha256"
```

Expected:

- Exit 0.
- Output contains `Checksum verified: sha256`.
- File exists at `$tmpdir/out-sha256.txt`.

## Matching SHA-512

```bash
go run . download http://127.0.0.1:18171/file.txt \
  --output "$tmpdir/out-sha512.txt" \
  --checksum "sha512:$expected_sha512"
```

Expected:

- Exit 0.
- Output contains `Checksum verified: sha512`.

## Mismatched SHA-256

```bash
go run . download http://127.0.0.1:18171/file.txt \
  --output "$tmpdir/out-mismatch.txt" \
  --checksum "sha256:0000000000000000000000000000000000000000000000000000000000000000"
echo "exit code: $?"
```

Expected:

- Exit non-zero.
- Error includes `checksum mismatch: expected`.
- Error includes the expected (`0000...`) and actual hex values.
- `$tmpdir/out-mismatch.txt` exists (file left in place).

## Dry-run with checksum

```bash
go run . download http://127.0.0.1:18171/file.txt \
  --output "$tmpdir/out-dry.txt" \
  --checksum "sha256:$expected_sha256" \
  --dry-run
```

Expected:

- Exit 0.
- Output contains `Checksum: sha256:<hex>`.
- No file created at `$tmpdir/out-dry.txt`.
- Output contains `Mode: dry-run only, no network request performed`.

## Root URL mode with --checksum

```bash
go run . http://127.0.0.1:18171/file.txt \
  --output "$tmpdir/out-root.txt" \
  --checksum "sha256:$expected_sha256"
```

Expected:

- Exit 0 (routes to CLI download mode, not TUI).
- Output contains `Checksum verified: sha256`.

## Invalid checksum format

```bash
go run . download http://127.0.0.1:18171/file.txt \
  --output "$tmpdir/out-invalid.txt" \
  --checksum "sha256:not-hex"
echo "exit code: $?"
```

Expected:

- Exit non-zero before any network request.
- Error includes `valid hexadecimal`.

## Batch rejection

```bash
sum="$(printf '%0.s0' {1..64})"
go run . download \
  http://127.0.0.1:18171/file.txt \
  http://127.0.0.1:18171/file.txt \
  --output "$tmpdir/out-batch.txt" \
  --checksum "sha256:$sum"
echo "exit code: $?"
```

Expected:

- Exit non-zero.
- Error: `--checksum is currently supported only for single URL downloads`.

## Teardown

```bash
kill "$SERVER_PID" 2>/dev/null
rm -rf "$tmpdir"
```

## Checklist

- [ ] SHA-256 match succeeds
- [ ] SHA-512 match succeeds
- [ ] SHA-256 mismatch exits non-zero with expected/actual
- [ ] Mismatch leaves file in place
- [ ] Dry-run shows checksum, no download performed
- [ ] Root URL mode routes to CLI with `--checksum`
- [ ] Invalid format fails before network
- [ ] Batch rejection with clear error
