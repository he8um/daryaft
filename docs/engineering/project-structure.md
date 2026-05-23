# Project Structure

## Local workspace

```text
Documents/Daryaft-project/
├── Daryaft/      # Git repository root. Commit this folder only.
├── Docs/         # Private agent docs. Never commit.
├── Caveman/      # Token-efficient agent protocol. Never commit.
└── Workspace/    # Local scratch files. Never commit.
```

## Repository structure

```text
Daryaft/
├── cmd/
├── internal/
│   ├── app/
│   ├── downloader/
│   ├── network/
│   ├── tui/
│   ├── updater/
│   ├── config/
│   ├── history/
│   ├── manifest/
│   ├── report/
│   └── utils/
├── pkg/
│   └── version/
├── docs/
├── scripts/
├── testdata/
├── .github/
├── README.md
├── CHANGELOG.md
├── CONTRIBUTING.md
├── SECURITY.md
├── LICENSE
├── Makefile
├── .goreleaser.yml
├── go.mod
└── main.go
```

## Rules

1. `cmd/` wires CLI commands only.
2. `internal/downloader/` owns download logic.
3. `internal/network/` owns HTTP clients, proxy, headers, timeouts.
4. `internal/tui/` owns rendering, screens, components, styles.
5. `internal/updater/` owns GitHub release updates.
6. `internal/config/` owns config loading, defaults, validation.
7. `internal/history/` owns history persistence.
8. `internal/manifest/` owns manifest and lockfile logic.
9. `pkg/version/` exposes build version metadata.
10. Private agent docs must never be copied into this repository.
