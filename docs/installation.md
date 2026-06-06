# Installation

Daryaft `v1.0.0` is the first stable release. Download binary archives from
the [GitHub releases page](https://github.com/he8um/daryaft/releases/tag/v1.0.0).
Package manager channels (Homebrew, deb, rpm, Arch) are post-1.0 work.

## Current Local Development

Use the Go toolchain:

```bash
go run . --help
go run . version
go run . version --json
go run .
```

Or build the binary:

```bash
make build
./bin/daryaft --help
```

For a local binary with injected build metadata:

```bash
make build-local
./bin/daryaft version
```

Source builds report version `1.1.0-dev`, commit `local`, build date
`unknown`, and built by `source`. Local ldflags builds can inject the current
git commit, UTC build time, and `built by` value. Release builds use GoReleaser
ldflags for the same metadata fields.

## Local Release Check

Use `make release-check` to run a local GoReleaser snapshot check:

```bash
make release-check
```

This requires GoReleaser v2. If it is missing, the target prints:

```text
GoReleaser is required. Install it with: brew install goreleaser
```

The target runs `goreleaser release --snapshot --clean --skip=publish`. It is
local only: it does not publish releases, create tags, or enable package-manager
publishing. Snapshot versions are named like
`1.1.0-dev-SNAPSHOT-<short-commit>`, and snapshot artifacts are written under
ignored local build directories such as `dist/`.

## Shell Completion

Daryaft can generate completion scripts for bash, zsh, fish, and PowerShell.
The commands print scripts to stdout; install locations vary by OS and shell
setup.

Examples:

```bash
daryaft completion zsh > "${fpath[1]}/_daryaft"
daryaft completion bash > /etc/bash_completion.d/daryaft
daryaft completion fish > ~/.config/fish/completions/daryaft.fish
```

For PowerShell:

```powershell
daryaft completion powershell
```

## Public Install Policy

`v1.0.0` is the first stable public release. GitHub release archives are the
supported install method at v1.0.0. Package manager channels (Homebrew, deb,
rpm, Arch, Scoop) are post-1.0 work and are not yet available.

## Planned Future Channels (Post-1.0)

These are post-1.0 planned channels, not yet available:

```bash
brew install he8um/tap/daryaft
```

```bash
curl -fsSL https://xhesam.com/daryaft/install.sh | sh
```

Related docs:

- [Quick Start](quick-start.md)
- [Versioning Policy](roadmap/versioning-policy.md)
