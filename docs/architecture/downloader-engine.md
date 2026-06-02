# Downloader Engine

HTTP download engine, resume, retry, filename, and conflict behavior.

## Navigation

- [Documentation Index](../index.md)
- [Implementation Plan](../engineering/implementation-plan.md)
- [Agent Navigation](../engineering/agent-navigation.md)
- [Versioning Policy](../roadmap/versioning-policy.md)
- [Release Train](../roadmap/release-train.md)

## Project constants

```text
Project: Daryaft
Binary: daryaft
Module: github.com/he8um/daryaft
Repository: git@github.com:he8um/daryaft.git
Author: AmirHesam Piri <info@xhesam.com>
Website: https://xhesam.com
Project page: https://xhesam.com/daryaft
License: MIT
Footer: Developed with <3 by AmirHesam Piri
```

## Requirements

1. Implement this area using clean, isolated packages.
2. Keep command wiring in `cmd/`; do not put business logic there.
3. Use typed errors and user-safe messages.
4. Update `daryaft -h` help text when user-facing commands or flags change.
5. Update tests and documentation in the same change.
6. Do not commit private agent docs from `Documents/Daryaft-project/Docs` or `Documents/Daryaft-project/Caveman`.

## Implementation notes

The agent must treat this file as a contract. If behavior is ambiguous, prefer the behavior documented in:

- `../engineering/interfaces-and-contracts.md`
- `../engineering/error-model.md`
- `../architecture/module-boundaries.md`

The default HTTP client uses transport phase timeouts rather than a global
`http.Client.Timeout`. Dial, TLS handshake, response-header, and idle
connection phases are bounded, but Daryaft does not impose a fixed total
download-duration timeout by default. Long body streams should continue as long
as the request context is active and the server keeps sending data.

Retry backoff is context-aware. Cancellation during backoff returns promptly
and cancelled downloads are not classified as retryable. Retry counts are
bounded to `0` through `20`; `0` means the first attempt only.

Resume/restart helpers own response bodies explicitly. If a Range response must
be replaced by a new full request, the old response body is closed once and the
new active response is returned to the caller. The caller closes only the
active response body after a successful preparation step.

## Acceptance criteria

- The feature is implemented in the correct module.
- The feature is covered by tests where practical.
- Errors are clear and actionable.
- The command help reflects the implemented behavior.
- The documentation includes examples and known limitations.

## Examples

```bash
daryaft -h
daryaft https://example.com/file.zip
daryaft -f urls.txt
# planned: daryaft update --check
```
