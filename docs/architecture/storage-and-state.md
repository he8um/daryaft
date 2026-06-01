# Storage and State

Config, history, queue, partial metadata, and OS-specific paths.

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

## Partial Downloads

Single URL downloads store incomplete bytes next to the target file:

```text
file.zip.part
file.zip.part.daryaft.json
```

The sidecar metadata contains:

- URL
- final target path
- partial path
- total bytes
- downloaded bytes
- `ETag`
- `Last-Modified`
- `Accept-Ranges`
- created and updated timestamps

Metadata writes use a temporary JSON file followed by a rename to keep updates
simple and recoverable. On successful completion, Daryaft renames the `.part`
file to the final target and removes the sidecar metadata.

On cancellation, Daryaft does not rename the `.part` file. It leaves both the
partial file and sidecar metadata in place so a later run can resume from the
saved byte count when the server supports HTTP Range requests.

Resume uses filesystem size as the source of truth for the local byte offset.
When `--resume` is enabled, Daryaft sends `Range: bytes=<partial_size>-` and
appends only after a `206 Partial Content` response. If the server returns a
full response, Daryaft truncates the partial file and restarts from byte `0`.
If saved `ETag` or `Last-Modified` metadata differs from the resume response,
Daryaft also restarts from byte `0`.

`--no-resume` ignores existing partial state for append decisions, truncates the
partial file, overwrites the sidecar metadata, and preserves the existing final
file protection.

## User Config

The YAML user config lives at `<UserConfigDir>/daryaft/config.yaml`. Saves
create the parent directory as needed, write a temporary file in the same
directory, set private `0600` permissions, and rename it into place. This avoids
direct partial writes to the final config path during normal updates.

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
daryaft update --check
```
