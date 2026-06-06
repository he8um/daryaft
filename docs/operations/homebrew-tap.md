# Homebrew Tap

Daryaft `v1.0.0` is the first stable release. A Homebrew tap is the first
planned package-manager install channel after `v1.0.0`.

This document covers the tap repository plan, formula strategy, validation
commands, and the steps required before GoReleaser Homebrew publishing can be
enabled.

## Status

- `v1.0.0` is the current stable release. GitHub binary assets are available.
- The `he8um/homebrew-tap` repository **does not yet exist**.
- GoReleaser Homebrew publishing is **not yet enabled** — the `brews:` block in
  `.goreleaser.yml` remains commented out.
- A formula draft is available for local review at
  [docs/operations/homebrew-formula-draft/daryaft.rb](homebrew-formula-draft/daryaft.rb).
  It is a reference file only — not installed or published.

## Tap Repository

| Field | Value |
|-------|-------|
| Tap owner | `he8um` |
| Tap repository | `he8um/homebrew-tap` |
| Formula path | `Formula/daryaft.rb` |
| `brew tap` command | `brew tap he8um/tap` |
| `brew install` command | `brew install he8um/tap/daryaft` |

The tap repository must be created before any of the above commands work.

### Creating the Tap Repository

When ready, the maintainer can create the tap with:

```bash
gh repo create he8um/homebrew-tap --public \
  --description "Homebrew tap for Daryaft" \
  --clone=false
```

Do not create the tap repository automatically. The formula must be validated
locally before it is pushed.

## Formula Strategy

- Formula installs from pre-built GitHub release archives, not from source.
- Formula selects the archive by architecture at install time.
- Formula pins the exact version and SHA-256 checksums.
- Formula does not use `latest` or floating URLs.

### v1.0.0 Asset SHA-256 Checksums

Extracted from `checksums.txt` attached to the
[v1.0.0 GitHub release](https://github.com/he8um/daryaft/releases/tag/v1.0.0):

| Archive | SHA-256 |
|---------|---------|
| `daryaft_darwin_arm64.tar.gz` | `5874045f452016dd2ccb61e347ab584deb50678b6d64f60d430bdf10fdcb1be3` |
| `daryaft_darwin_amd64.tar.gz` | `711f6e5cffe77c2d87534119443c1837946e59d02d318b0fc4d9f8c52faa3eca` |

### Formula Draft

See [docs/operations/homebrew-formula-draft/daryaft.rb](homebrew-formula-draft/daryaft.rb).

The draft formula selects URL and SHA-256 by CPU architecture using the
standard Homebrew `Hardware::CPU.arm?` guard, installs `daryaft` into `bin`,
and tests with `system "#{bin}/daryaft", "version"`.

## Linuxbrew

Linux archive assets (`daryaft_linux_amd64.tar.gz`,
`daryaft_linux_arm64.tar.gz`) are shipped with the `v1.0.0` release. Adding
Linuxbrew support to the formula requires an `on_linux` block with the
corresponding URL and SHA-256 values.

Linuxbrew coverage is planned for a future formula update after the macOS
formula is validated. It is not required for the initial tap launch.

## Security Requirements

- Formula must pin exact version and SHA-256 checksums. Never use `latest`.
- SHA-256 values must be taken from the `checksums.txt` attached to the
  official GitHub release, not computed from a re-downloaded archive.
- Formula must reference official GitHub release archives from
  `https://github.com/he8um/daryaft/releases/download/`.
- On each new release, update `version`, both `url` entries, and both `sha256`
  entries.

## Validation Commands

After placing the formula at `Formula/daryaft.rb` in the tap repository, run:

```bash
# Syntax check
ruby -c Formula/daryaft.rb

# Homebrew audit (requires local Homebrew)
brew audit --strict --new-formula Formula/daryaft.rb

# Install from local formula file
brew install --verbose Formula/daryaft.rb

# Verify binary
daryaft version
daryaft doctor

# Run formula test
brew test daryaft
```

For end-to-end tap validation after the tap repository is live:

```bash
brew tap he8um/tap
brew install he8um/tap/daryaft
daryaft version
daryaft doctor
brew test he8um/tap/daryaft
```

For upgrades after future releases:

```bash
brew update
brew upgrade daryaft
```

## GoReleaser Homebrew Publishing

GoReleaser can automate formula updates via a `brews:` block in
`.goreleaser.yml`. This is not yet enabled.

### Prerequisite Checklist Before Enabling

- [ ] `he8um/homebrew-tap` repository exists and is public.
- [ ] `Formula/daryaft.rb` is in the tap repository and manually validated.
- [ ] `brew install`, `daryaft version`, `daryaft doctor`, and `brew test`
      all pass locally.
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

## Post-1.0 Homebrew Roadmap

1. Create `he8um/homebrew-tap`.
2. Add `Formula/daryaft.rb` using the draft from this repository.
3. Validate formula locally with `ruby -c`, `brew audit`, `brew install`,
   `daryaft version`, `daryaft doctor`, `brew test`.
4. Push formula to tap manually for the first install.
5. Set up `HOMEBREW_TAP_GITHUB_TOKEN` CI secret.
6. Enable GoReleaser `brews:` block after manual validation.
7. Add Linuxbrew `on_linux` block to formula.
8. Document upgrade path in `docs/installation.md`.

## References

- [v1.0.0 Release Assets](release-assets.md)
- [v1.0.0 Release Notes](release-notes-v1.0.0.md)
- [Release Process](release-process.md)
- [Installation](../installation.md)
- [Post-1.0 Feature Packs](../roadmap/post-1-feature-packs.md)
- [Versioning Policy](../roadmap/versioning-policy.md)
- [Formula Draft](homebrew-formula-draft/daryaft.rb)
