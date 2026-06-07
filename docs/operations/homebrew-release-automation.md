# Homebrew Release Automation

This document describes the Homebrew formula update workflow for Daryaft
releases: the current manual process, the intermediate helper script, and the
future path toward GoReleaser `brews:` publishing.

## Current Manual Process

Each release currently requires editing `Formula/daryaft.rb` in
`he8um/homebrew-tap` by hand:

1. Publish the GitHub release with binary assets and `checksums.txt`.
2. Open `he8um/homebrew-tap/Formula/daryaft.rb`.
3. Update `version "X.Y.Z"`.
4. Update both `url` entries to the new tag path.
5. Update both `sha256` values from the release `checksums.txt`.
6. Run `ruby -c Formula/daryaft.rb` to check Ruby syntax.
7. Commit and push the formula update.
8. Verify with `brew update && brew upgrade daryaft && daryaft version`.

This works but requires looking up checksums and carefully editing four
string values per release. The `scripts/update-homebrew-formula.sh` helper
automates steps 2–6.

## Why GoReleaser `brews:` Publishing Remains Disabled

GoReleaser can push formula updates to the tap automatically after a release.
This is not yet enabled because:

- It requires a `HOMEBREW_TAP_GITHUB_TOKEN` secret with write access to
  `he8um/homebrew-tap`.
- Token scope, rotation policy, and audit logging are not yet documented.
- The release pipeline has not been reviewed for automated publishing safety.
- The manual helper workflow needs to be trusted through at least one release
  cycle before automation is added.

The GoReleaser `brews:` block placeholder is already in `.goreleaser.yml`
(commented out). See [Homebrew Tap](homebrew-tap.md) for the enabling
prerequisite checklist.

## Recommended Intermediate Workflow

```text
GitHub release created and assets uploaded
  ↓
Run: scripts/update-homebrew-formula.sh --version X.Y.Z --tap-dir /tmp/homebrew-tap
  ↓
Script fetches checksums.txt from GitHub release
  ↓
Script updates Formula/daryaft.rb in local tap checkout
  ↓
Maintainer reviews diff
  ↓
Maintainer runs: ruby -c Formula/daryaft.rb
Maintainer runs: brew reinstall he8um/tap/daryaft
Maintainer runs: daryaft version && daryaft update --check && daryaft doctor
  ↓
Maintainer commits and pushes tap update manually
```

The script never pushes, commits, or creates releases. Every step that changes
the live tap requires a deliberate manual action.

## Using the Helper Script

### Prerequisites

- A local clone of `he8um/homebrew-tap` (working tree must be clean).
- `curl` available.
- The GitHub release for the target version must already exist and have assets.

### Clone the tap

```bash
git clone https://github.com/he8um/homebrew-tap.git /tmp/homebrew-tap
```

### Dry run (preview changes without editing)

```bash
scripts/update-homebrew-formula.sh --version 1.2.0 --tap-dir /tmp/homebrew-tap --dry-run
```

### Apply changes

```bash
scripts/update-homebrew-formula.sh --version 1.2.0 --tap-dir /tmp/homebrew-tap
```

### Via make

```bash
make homebrew-formula-update VERSION=1.2.0 TAP_DIR=/tmp/homebrew-tap
make homebrew-formula-update-dry-run VERSION=1.2.0 TAP_DIR=/tmp/homebrew-tap
```

### Script output

On success the script prints:

```text
=== Formula update complete ===
  Version : 1.2.0
  arm64   : <sha256>
  amd64   : <sha256>

=== Diff ===
<git diff of Formula/daryaft.rb>

=== Next steps (manual) ===
  1. Review diff above.
  2. Run: ruby -c Formula/daryaft.rb
  3. Run: brew reinstall he8um/tap/daryaft
  4. Run: daryaft version && daryaft update --check && daryaft doctor
  5. Run in tap dir:
       git add Formula/daryaft.rb
       git commit -m "Update daryaft to v1.2.0"
       git push
```

### Idempotence

If the formula already matches the target version and checksums, the script
exits with a message and makes no changes. Running it twice for the same
version is safe.

## Script Safety Rules

The script enforces these rules before modifying anything:

- Version must be in `X.Y.Z` form. Dev suffixes (`-dev`, `-rc.1`) are rejected.
- Tap directory must exist and be a git repository.
- `Formula/daryaft.rb` must be present.
- Tap working tree must be clean before editing.
- Both `daryaft_darwin_arm64.tar.gz` and `daryaft_darwin_amd64.tar.gz`
  checksums must be present in the downloaded `checksums.txt`.
- Each checksum must match the 64-character hex SHA-256 pattern.
- Formula must contain recognizable version, CPU guard, and archive references.
- Script exits non-zero on any failure.

The script never:
- Pushes to the tap.
- Commits to the tap.
- Creates GitHub releases.
- Creates tags.
- Uploads release assets.
- Modifies the Daryaft source repository.

## Required Inputs

| Input | Example | Notes |
|-------|---------|-------|
| `--version` | `1.2.0` or `v1.2.0` | Leading `v` is stripped automatically |
| `--tap-dir` | `/tmp/homebrew-tap` | Must be a clean git checkout of `he8um/homebrew-tap` |
| `--dry-run` | (flag) | Optional; previews diff without writing |
| `--repo` | `he8um/daryaft` | Optional; defaults to `he8um/daryaft` |

## Fields Updated in the Formula

The script updates exactly these fields:

```ruby
version "X.Y.Z"

if Hardware::CPU.arm?
  url "https://github.com/he8um/daryaft/releases/download/vX.Y.Z/daryaft_darwin_arm64.tar.gz"
  sha256 "<arm64-sha256>"
else
  url "https://github.com/he8um/daryaft/releases/download/vX.Y.Z/daryaft_darwin_amd64.tar.gz"
  sha256 "<amd64-sha256>"
end
```

All other formula content (`desc`, `homepage`, `license`, `install`, `test`) is
preserved unchanged.

## Post-Update Validation Commands

```bash
# In the tap checkout
ruby -c Formula/daryaft.rb

# After brew upgrade or reinstall
brew reinstall he8um/tap/daryaft
daryaft version
daryaft update --check
daryaft doctor

# Formula test
HOMEBREW_NO_AUTO_UPDATE=1 brew test --verbose he8um/tap/daryaft
```

## Future Path: GoReleaser `brews:` Publishing

Once the manual helper workflow has been used through at least one release
cycle and the following prerequisites are met, GoReleaser tap publishing can
be enabled:

- [ ] `HOMEBREW_TAP_GITHUB_TOKEN` secret configured with write-only access to
      `he8um/homebrew-tap`.
- [ ] Token scope and rotation policy documented.
- [ ] GoReleaser `brews:` block reviewed for correctness.
- [ ] Dry-run of `goreleaser release --skip=publish` confirms formula output.
- [ ] Release pipeline trusted and publish step reviewed.

Until then, the `brews:` block in `.goreleaser.yml` remains commented out.
See [Homebrew Tap](homebrew-tap.md) for the full enabling checklist.

## References

- [Homebrew Tap](homebrew-tap.md)
- [Release Process](release-process.md)
- [v1.0.0 Release Assets](release-assets.md)
- [Installation](../installation.md)
- [Script source](../../scripts/update-homebrew-formula.sh)
