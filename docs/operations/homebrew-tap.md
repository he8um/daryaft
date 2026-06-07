# Homebrew Tap

Daryaft is available via a Homebrew tap. The `he8um/tap` tap is the first live
package-manager install channel for Daryaft.

## Status

- `v1.1.0` is the current stable release. GitHub binary assets are available.
- `he8um/homebrew-tap` exists at https://github.com/he8um/homebrew-tap.
- `Formula/daryaft.rb` is live in the tap repository, updated to `v1.1.0`.
- Homebrew tap installation has been validated:
  - `brew tap he8um/tap` succeeded.
  - `brew install daryaft` (or `brew install he8um/tap/daryaft`) succeeded.
  - `daryaft version` reports `1.1.0`.
  - `daryaft doctor` runs successfully.
  - `daryaft update --check` works correctly.
- GoReleaser Homebrew publishing is **not yet enabled** — the `brews:` block in
  `.goreleaser.yml` remains commented out.
- The formula is manually maintained. Future releases require updating
  `Formula/daryaft.rb` in `he8um/homebrew-tap` with the new version and
  checksums unless GoReleaser tap publishing is enabled later.

## Install

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

A Homebrew trust warning may appear when tapping because this is a
user-owned custom tap, not a Homebrew core formula. This is expected behavior
for third-party taps.

## Upgrade

After a new release is published and the formula is updated in `he8um/homebrew-tap`:

```bash
brew update
brew upgrade daryaft
```

## Checking for Updates

Use `daryaft update --check` to verify whether a newer release is available
before upgrading:

```bash
daryaft update --check
```

When installed via Homebrew, the output shows:

```text
Install channel:  homebrew
Update command:   brew update && brew upgrade daryaft
```

`daryaft update --check` is read-only and does not install anything.

## Verify

After install or upgrade:

```bash
daryaft version
daryaft doctor
```

To run the Homebrew formula test:

```bash
HOMEBREW_NO_AUTO_UPDATE=1 brew test --verbose he8um/tap/daryaft
```

## Tap Repository

| Field | Value |
|-------|-------|
| Tap owner | `he8um` |
| Tap repository | `he8um/homebrew-tap` |
| Repository URL | https://github.com/he8um/homebrew-tap |
| Formula path | `Formula/daryaft.rb` |
| `brew tap` command | `brew tap he8um/tap` |
| `brew install` command | `brew install daryaft` |
| Explicit install | `brew install he8um/tap/daryaft` |

## Formula Details

| Field | Value |
|-------|-------|
| Current formula version | `1.1.0` |
| Formula source | GitHub v1.1.0 release assets |
| macOS Apple Silicon | `daryaft_darwin_arm64.tar.gz` |
| macOS Intel | `daryaft_darwin_amd64.tar.gz` |
| Install method | Pre-built binary archive; does not build from source |

### v1.1.0 Asset SHA-256 Checksums

Extracted from `checksums.txt` attached to the
[v1.1.0 GitHub release](https://github.com/he8um/daryaft/releases/tag/v1.1.0):

See [v1.1.0 Release Notes](release-notes-v1.1.0.md) for the full asset list.
SHA-256 values are available in `checksums.txt` attached to the release.

### v1.0.0 Asset SHA-256 Checksums (historical)

Extracted from `checksums.txt` attached to the
[v1.0.0 GitHub release](https://github.com/he8um/daryaft/releases/tag/v1.0.0):

| Archive | SHA-256 |
|---------|---------|
| `daryaft_darwin_arm64.tar.gz` | `5874045f452016dd2ccb61e347ab584deb50678b6d64f60d430bdf10fdcb1be3` |
| `daryaft_darwin_amd64.tar.gz` | `711f6e5cffe77c2d87534119443c1837946e59d02d318b0fc4d9f8c52faa3eca` |

## Formula Maintenance

The formula is manually maintained for now. When a new Daryaft release is
published:

1. Build and publish the new GitHub release assets.
2. Clone the tap: `git clone https://github.com/he8um/homebrew-tap.git /tmp/homebrew-tap`
3. Run the helper script:
   ```bash
   scripts/update-homebrew-formula.sh --version X.Y.Z --tap-dir /tmp/homebrew-tap
   ```
   The script fetches `checksums.txt` from the GitHub release, updates
   `Formula/daryaft.rb` in the local clone, and prints next steps.
4. Review the diff, validate with `ruby -c` and `brew reinstall`.
5. In the tap clone: `git add Formula/daryaft.rb && git commit -m "Update daryaft to vX.Y.Z" && git push`
6. Verify with `brew update && brew upgrade daryaft && daryaft version`.

For a dry run preview without modifying files:

```bash
scripts/update-homebrew-formula.sh --version X.Y.Z --tap-dir /tmp/homebrew-tap --dry-run
```

Or use the Makefile helper:

```bash
make homebrew-formula-update VERSION=X.Y.Z TAP_DIR=/tmp/homebrew-tap
make homebrew-formula-update-dry-run VERSION=X.Y.Z TAP_DIR=/tmp/homebrew-tap
```

The script never pushes, commits, or creates releases. See
[Homebrew Release Automation](homebrew-release-automation.md) for full details.

A formula draft reference is kept in this repository at
[docs/operations/homebrew-formula-draft/daryaft.rb](homebrew-formula-draft/daryaft.rb)
for review, but the canonical formula lives in `he8um/homebrew-tap`.

## Linuxbrew

Linux archive assets (`daryaft_linux_amd64.tar.gz`,
`daryaft_linux_arm64.tar.gz`) are available in the v1.0.0 release. Adding
Linuxbrew support requires an `on_linux` block in the formula with the
corresponding URL and SHA-256 values. Linuxbrew support is future work.

## Security Requirements

- Formula pins exact version and SHA-256 checksums. It does not use `latest`.
- SHA-256 values are taken from the `checksums.txt` attached to the official
  GitHub release.
- Formula references official GitHub release archives from
  `https://github.com/he8um/daryaft/releases/download/`.
- On each new release, update `version`, both `url` entries, and both `sha256`
  entries in the formula.

## GoReleaser Homebrew Publishing

GoReleaser can automate formula updates via a `brews:` block in
`.goreleaser.yml`. This is not yet enabled.

### Prerequisite Checklist Before Enabling

- [x] `he8um/homebrew-tap` repository exists and is public.
- [x] `Formula/daryaft.rb` is in the tap repository and manually validated.
- [x] `brew install`, `daryaft version`, and `daryaft doctor` all pass.
- [ ] A `HOMEBREW_TAP_GITHUB_TOKEN` (or equivalent) secret is configured in
      the `he8um/daryaft` repository CI secrets — with write access to
      `he8um/homebrew-tap` only.
- [ ] Token scope and rotation policy are documented.
- [ ] The GoReleaser `brews:` block has been reviewed for correctness:
      `skip_upload`, `commit_author`, `folder`, `homepage`, `description`.
- [ ] A dry-run of the GoReleaser release pipeline (with `--skip=publish`)
      confirms the formula would be generated correctly before a live run.
- [ ] The release workflow is trusted and the publish step is reviewed.

The commented-out `brews:` block in `.goreleaser.yml` is the placeholder:

```yaml
# brews:
#   - name: daryaft
#     repository:
#       owner: he8um
#       name: homebrew-tap
#     homepage: https://xhesam.com/daryaft
#     description: Modern terminal downloader written in Go.
#     license: MIT
```

Do not uncomment this block until all prerequisite checks above are satisfied.

## References

- [Homebrew Release Automation](homebrew-release-automation.md)
- [v1.1.0 Release Notes](release-notes-v1.1.0.md)
- [v1.0.0 Release Assets](release-assets.md)
- [v1.0.0 Release Notes](release-notes-v1.0.0.md)
- [Release Process](release-process.md)
- [Installation](../installation.md)
- [Post-1.0 Feature Packs](../roadmap/post-1-feature-packs.md)
- [Versioning Policy](../roadmap/versioning-policy.md)
- [Formula Draft](homebrew-formula-draft/daryaft.rb)
