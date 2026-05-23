# Security Policy

## Reporting security issues

Report security issues privately to:

```text
info@xhesam.com
```

Do not open public GitHub issues for vulnerabilities.

## Security model

Daryaft will download files from user-provided URLs. It must never execute
downloaded files automatically. Future features that run hooks, execute commands,
or scan files must be explicit opt-in features.

Read:

- `docs/architecture/security-model.md`
- `docs/features/security-and-validation.md`
- `docs/engineering/error-model.md`

## Update security

Self-update is planned but not implemented. When added, it must verify checksums
before replacing the binary and attempt rollback if replacement fails.
