# Daryaft v1.9.0

Daryaft v1.9.0 is a stable release focused on download retry and resume reliability hardening.

## Highlights

- Added HTTP `408 Request Timeout` to retryable download responses.
- Locked retry classification with deterministic tests.
- Added deterministic integration coverage for `408` followed by successful retry.
- Hardened resume behavior for oversized local partial files.
- Locked missing/corrupt resume sidecar behavior with tests.
- Documented safe restart behavior for `416`, partial-overflow, and sidecar-loss cases.
- Documented cancellation exit behavior.
- Documented that `429 Retry-After` is not honored yet.
- Added v1.9.0 roadmap/scope documentation.

## Retry behavior

Daryaft now retries these HTTP statuses:

```text
408
429
500
502
503
504
```

These common client-side statuses remain non-retryable:

```text
400
401
403
404
410
```

`429 Retry-After` headers are not honored yet. Daryaft continues to use its standard exponential backoff.

## Resume behavior

Daryaft already supported safe resume behavior for `Range` downloads.

This release additionally hardens the case where a local partial file is larger than the known remote file size. In that case, Daryaft restarts safely from byte `0` instead of appending to an unsafe partial file.

Resume behavior remains:

* `206 Partial Content`: append when validators and content range are safe.
* `200 OK` after a `Range` request: restart from byte `0`.
* `416 Requested Range Not Satisfiable`: restart from byte `0` as a safe default.
* missing/corrupt sidecar metadata: restart safely or treat as non-resumable; no unsafe append.

## Cancellation behavior

Cancellation exits non-zero and preserves partial download state where applicable.

Cancelled batch items remain cancelled, and remaining items are reported as not started.

## Documentation

Updated or added:

* `docs/features/resume-and-retry.md`
* `docs/command-reference.md`
* `docs/usage.md`
* `docs/operations/manual-qa.md`
* `docs/roadmap/v1.9.0-download-reliability-hardening.md`
* `docs/index.md`
* `CHANGELOG.md`

## Known limitations

* `Retry-After` headers are not honored yet.
* `416` always restarts from byte `0` instead of treating the partial file as complete.
* Concurrent batch downloads are not implemented.
* Config persistence is not implemented.
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
