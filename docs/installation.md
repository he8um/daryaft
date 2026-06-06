# Installation

Daryaft `v1.0.0` is the first stable release. Install via Homebrew or
download binary archives directly from the
[GitHub releases page](https://github.com/he8um/daryaft/releases/tag/v1.0.0).
Other package manager channels (deb, rpm, Arch) are post-1.0 work.

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

## Homebrew

The `he8um/tap` Homebrew tap is the first live package-manager install channel
for Daryaft.

```bash
brew tap he8um/tap
brew install daryaft
daryaft version
daryaft doctor
```

Alternatively, install with the fully-qualified tap name:

```bash
brew install he8um/tap/daryaft
```

A Homebrew trust warning may appear when tapping because this is a user-owned
custom tap. This is expected for third-party taps.

To upgrade after a new release:

```bash
brew update
brew upgrade daryaft
```

See [Homebrew Tap](operations/homebrew-tap.md) for formula details, SHA-256
checksums, maintenance instructions, and the GoReleaser publishing checklist.

## GitHub Binary Archives

Download binary archives directly for any supported platform:

```bash
# Example: macOS Apple Silicon
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.0.0/daryaft_darwin_arm64.tar.gz
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.0.0/checksums.txt
shasum -a 256 --check checksums.txt
tar -xzf daryaft_darwin_arm64.tar.gz
./daryaft version
```

Available archives: `daryaft_linux_amd64.tar.gz`, `daryaft_linux_arm64.tar.gz`,
`daryaft_darwin_amd64.tar.gz`, `daryaft_darwin_arm64.tar.gz`.

## Other Planned Future Channels (Post-1.0)

These channels are not yet available:

```bash
curl -fsSL https://xhesam.com/daryaft/install.sh | sh
```

Other package manager channels (deb, rpm, Arch, Scoop) are later post-1.0 work.

Related docs:

- [Quick Start](quick-start.md)
- [Versioning Policy](roadmap/versioning-policy.md)
