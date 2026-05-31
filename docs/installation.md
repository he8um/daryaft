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
