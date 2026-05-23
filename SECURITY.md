# Security Policy

## Reporting security issues

Report security issues privately to:

```text
info@xhesam.com
```

Do not open public GitHub issues for security vulnerabilities.

## Security model

Daryaft downloads files from user-provided URLs. It must not execute downloaded files automatically.
Features that execute user commands, run hooks, or scan files must be explicit opt-in features.

Read:

- `docs/architecture/security-model.md`
- `docs/features/security-and-validation.md`
- `docs/engineering/error-model.md`

## Update security

Self-update must verify checksums before replacing the binary.
Rollback must be attempted when replacement fails.
