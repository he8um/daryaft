# Daryaft

Daryaft is a beautiful, fast, terminal-based download manager for macOS and Linux.
It is inspired by tools like `wget`, but adds a modern TUI, queue management, self-update, package releases, and automation-friendly output.

> Developed with <3 by AmirHesam Piri

## Status

Daryaft is planned as a public open-source project at:

- Repository: https://github.com/he8um/daryaft
- Project page: https://xhesam.com/daryaft
- Author: AmirHesam Piri <info@xhesam.com>
- License: MIT

## Pre-1.0 policy

The repository may be public before `v1.0.0`, but pre-1.0 versions are not public installation releases.
Do not publish Homebrew, Debian, RPM, Arch, or one-line install channels before `v1.0.0`.

Before `v1.0.0`, releases may exist as internal/local tags or GitHub pre-releases for development validation only.
Users must not be directed to install Daryaft before the first stable release.

## Stable installation policy

From `v1.0.0` onward, users should be able to install the latest stable version with one command and optionally install a specific version.

Examples planned for stable releases:

```bash
brew install he8um/tap/daryaft
```

```bash
curl -fsSL https://xhesam.com/daryaft/install.sh | sh
```

Specific version installation is planned through installer flags and package manager versions:

```bash
curl -fsSL https://xhesam.com/daryaft/install.sh | sh -s -- --version v1.0.0
```

## Documentation

Repository documentation lives in `docs/`.
Private implementation docs for agents live outside this repository in `Documents/Daryaft-project/Docs` and must not be committed.

Start here:

- [Quick Start](docs/quick-start.md)
- [Command Reference](docs/command-reference.md)
- [Architecture Overview](docs/architecture/overview.md)
- [Roadmap](docs/roadmap/index.md)
- [Release Policy](docs/roadmap/versioning-policy.md)

## Core commands

```bash
daryaft
```

```bash
daryaft https://example.com/file.zip
```

```bash
daryaft -f urls.txt
```

```bash
daryaft update
```

## Footer

The TUI footer must show:

```text
Developed with <3 by AmirHesam Piri
```

If terminal hyperlink support is available, `AmirHesam Piri` should link to:

```text
https://xhesam.com
```
