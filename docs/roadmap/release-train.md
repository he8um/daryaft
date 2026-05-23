# Release Train

This document defines how Daryaft releases move from local development to public installation.

## Repository

```text
git@github.com:he8um/daryaft.git
```

## Public repository policy

The repository is public, but public installation channels must stay disabled until `v1.0.0`.

This means:

- Code can be pushed to GitHub.
- Documentation can exist publicly.
- Issues and project boards can exist.
- Pre-1.0 tags can exist only as development/pre-release markers.
- README must clearly say installation starts from `v1.0.0`.

## Installation channel gate

Before `v1.0.0`:

```text
build: yes
local test: yes
GitHub pre-release: optional
Homebrew publish: no
deb/rpm/arch public promotion: no
install.sh public use: no
```

From `v1.0.0`:

```text
GitHub stable release: yes
Homebrew tap: yes
install.sh latest: yes
install.sh --version: yes
package assets: yes
checksums: yes
```

## Latest install command

From `v1.0.0` onward:

```bash
brew install he8um/tap/daryaft
```

Alternative:

```bash
curl -fsSL https://xhesam.com/daryaft/install.sh | sh
```

## Specific version install

```bash
curl -fsSL https://xhesam.com/daryaft/install.sh | sh -s -- --version v1.0.0
```

## GoReleaser guard

Release automation must check the tag. If tag starts with `v0.`, do not publish public package channels.

Pseudo-code:

```text
version = git tag
if version starts with v0.:
    build archives only
    mark as prerelease if published
    skip Homebrew and package publishing
else:
    publish full stable release
```
