# v1.0.0 Release Asset Strategy

## Decision

Binary assets **will be attached** to the `v1.0.0` GitHub release.

The `v0.6.0-rc.2` pre-release shipped as source/tag-only, which is acceptable
for an internal RC. The stable v1.0.0 release is a public milestone and must
include compiled binaries so users can download and run Daryaft without a local
Go toolchain.

## Asset List

GoReleaser generates the following artifacts from the `.goreleaser.yml`
configuration on a clean tag build (`goreleaser release`):

| File | Description |
|------|-------------|
| `daryaft_linux_amd64.tar.gz` | Linux x86-64 binary archive |
| `daryaft_linux_arm64.tar.gz` | Linux ARM64 binary archive |
| `daryaft_darwin_amd64.tar.gz` | macOS x86-64 binary archive |
| `daryaft_darwin_arm64.tar.gz` | macOS ARM64 (Apple Silicon) binary archive |
| `checksums.txt` | SHA-256 checksums for all archives |

Each archive contains: `daryaft` binary, `CHANGELOG.md`, `LICENSE`, `README.md`.

## What Is Not Included at v1.0.0

- Windows binaries — Windows is not officially supported at v1.0.0.
- Package manager publishing (Homebrew tap, deb, rpm, Arch) — post-1.0 work.
  The nfpms and brews blocks in `.goreleaser.yml` remain commented out.
- Signed checksums — not implemented at v1.0.0.
- Automatic asset upload pipeline — the release workflow does not have a
  tag-triggered release job; asset upload is a manual step at v1.0.0.

## GoReleaser Configuration

No `.goreleaser.yml` changes are required for v1.0.0 binary assets.

The current configuration already:
- Builds linux/amd64, linux/arm64, darwin/amd64, darwin/arm64.
- Creates `tar.gz` archives with the correct name template.
- Generates `checksums.txt`.
- Injects version metadata via ldflags.
- Keeps package-manager publishing commented out with a clear note.

The `snapshot.version_template` (`0.6.0-dev-SNAPSHOT-{{ .ShortCommit }}`) is
used only for local snapshot builds (`make release-check`). A non-snapshot
`goreleaser release` on the `v1.0.0` tag will produce binaries reporting
version `1.0.0`.

## Build Process

To produce the v1.0.0 release artifacts:

```bash
# Verify the tag exists and CI is green on the release commit
git fetch --tags
git describe --tags

# Validate configuration
goreleaser check

# Build release artifacts (local verification before upload)
goreleaser release --clean --skip=publish

# Inspect artifacts
find dist -maxdepth 2 -type f | sort
cat dist/checksums.txt
```

Then upload `dist/*.tar.gz` and `dist/checksums.txt` to the GitHub release.

Do not use `--snapshot` for the release build. The `--snapshot` flag modifies
the version string; only a real `goreleaser release` on a clean tag injects the
correct version.

## Validation Commands

After uploading assets to the GitHub release, verify checksums:

```bash
# Download and verify a specific archive
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.0.0/daryaft_darwin_arm64.tar.gz
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.0.0/checksums.txt
shasum -a 256 --check checksums.txt
```

Spot-check the binary version:

```bash
tar -xzf daryaft_darwin_arm64.tar.gz
./daryaft version
./daryaft version --json
```

Expected output: `version: 1.0.0`, `built_by: goreleaser`.

## Notes

- Package-manager publishing (Homebrew, deb, rpm, Arch) is post-1.0 work.
  Do not enable the commented-out sections of `.goreleaser.yml` at v1.0.0.
- The CI workflow does not have a tag-triggered `goreleaser release` job.
  Asset upload at v1.0.0 is a manual step.
- Automated release pipeline is post-1.0 work.

## References

- [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md)
- [v0.6.0-rc.2 Release Status](release-status-v0.6.0-rc.2.md)
- [Release Process](release-process.md)
- [Release-Candidate Validation](rc-validation.md)
