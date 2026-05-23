# Homebrew Tap

Tap repo rules and post-v1 publishing.

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
