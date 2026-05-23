# Installation

Daryaft does not have public stable install channels yet.

## Current Local Development

Use the Go toolchain:

```bash
go run . --help
go run . version
go run .
```

Or build the binary:

```bash
make build
./bin/daryaft --help
```

## Public Install Policy

Public install channels are planned for `v1.0.0` and later. Before `v1.0.0`,
Homebrew, Debian, RPM, Arch, GitHub release archives, and install scripts are
configuration stubs only, not official stable install paths.

## Planned Stable Channels

These are planned examples for the `v1.0.0` era:

```bash
brew install he8um/tap/daryaft
```

```bash
curl -fsSL https://xhesam.com/daryaft/install.sh | sh
```

Related docs:

- [Quick Start](quick-start.md)
- [Versioning Policy](roadmap/versioning-policy.md)
