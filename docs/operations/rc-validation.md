# Release-Candidate Validation

This guide describes how to validate Daryaft internal release-candidate tags.
It is for pre-1.0 release engineering only and does not publish releases.

## Current RC

`v0.6.0-rc.2` is the current internal validation release candidate for the
`0.6.0-dev` milestone. It is not a public stable release, does not imply a
public install channel, and must not be promoted as a stable user-facing
release. Public stable remains planned for `v1.0.0`.

`v0.6.0-rc.1` is the previous internal RC. It has been superseded by
`v0.6.0-rc.2`. See
[Daryaft v0.6.0-rc.1 Internal Release Candidate](release-notes-v0.6.0-rc.1.md)
for its historical notes.

The source default version remains `0.6.0-dev`. GoReleaser release and snapshot
builds inject build metadata with linker flags.

## Local Tag Validation

Fetch tags and inspect the release-candidate state without creating new tags:

```bash
git fetch --tags
git tag --list "v0.6.0-rc.*"
git describe --tags --always
git log --oneline --decorate --graph --all -12
```

To inspect the RC tag directly, use a detached checkout or inspect it without
changing branches:

```bash
git checkout v0.6.0-rc.2
git show --stat v0.6.0-rc.2
```

If you check out the tag, return to the working branch before making changes:

```bash
git checkout main
```

Do not create replacement tags during validation. If the expected tag is
missing or points at an unexpected commit, stop and report the mismatch.

## Validation Commands

Run the release-candidate quality checks from the repository root:

```bash
go run . version
go run . version --json
goreleaser check
make rc-check
make release-check
```

`make rc-check` includes blocking `govulncheck` and `gosec` checks. `make
security` also runs both as blocking checks locally.

`make release-check` runs:

```bash
goreleaser release --snapshot --clean --skip=publish
```

It must not publish a GitHub release, upload artifacts, create tags, or enable
package-manager publishing.

## Artifact Inspection

After `make release-check`, inspect generated artifacts under ignored local
directories:

```bash
find dist -maxdepth 2 -type f | sort
ls -lh dist
cat dist/checksums.txt
shasum -a 256 dist/checksums.txt
```

Expected snapshot archive names currently include:

- `dist/daryaft_darwin_amd64.tar.gz`
- `dist/daryaft_darwin_arm64.tar.gz`
- `dist/daryaft_linux_amd64.tar.gz`
- `dist/daryaft_linux_arm64.tar.gz`

Each archive should contain:

- `CHANGELOG.md`
- `LICENSE`
- `README.md`
- `daryaft`

To inspect embedded metadata from a compatible local artifact:

```bash
rm -rf /tmp/daryaft-rc-artifact
mkdir -p /tmp/daryaft-rc-artifact
tar -xzf dist/daryaft_darwin_arm64.tar.gz -C /tmp/daryaft-rc-artifact
/tmp/daryaft-rc-artifact/daryaft version
/tmp/daryaft-rc-artifact/daryaft version --json
```

On non-Darwin ARM64 hosts, choose an archive compatible with the local system or
inspect metadata through GoReleaser metadata files.

Confirm generated artifacts are ignored and not staged:

```bash
git status --short --ignored dist bin
git status --short
```

## GitHub Actions Tag Behavior

The current test workflow runs on `push` and `pull_request`. Because the `push`
trigger is not filtered to branches only, GitHub Actions runs for tag pushes as
well as branch pushes. This means an internal RC tag push such as
`v0.6.0-rc.2` runs the test, GoReleaser config, lint, and security jobs.
GitHub Actions for `v0.6.0-rc.2` passed.

The workflow does not publish releases and does not run a tag-triggered release
job.

## Optional GitHub Pre-Release (Manual Only)

If you want to publish an internal GitHub pre-release for sharing, use:

```bash
gh release create v0.6.0-rc.2 \
  --title "Daryaft v0.6.0-rc.2" \
  --notes-file docs/operations/release-notes-v0.6.0-rc.2.md \
  --prerelease \
  --verify-tag
```

Do not run this unless intentionally publishing a GitHub pre-release. Do not
enable package-manager artifact publishing from an RC tag. This command is
provided for reference only — it is not run automatically.

## Confirm No Release Was Published

During local validation:

- `make release-check` output should include `--skip=publish`.
- GoReleaser should report that announce, publish, and validate are skipped.
- No GitHub Release should be created.
- No new tag should be created.
- `dist/` and `bin/` should remain ignored local artifacts.

If any command attempts to publish, upload artifacts, create a release, or
modify tags, stop and report it.

## Finding Record

Record each validation pass with:

- RC tag and commit.
- Date.
- Commands run.
- GoReleaser artifacts generated.
- Version metadata observations.
- Quality gate results.
- Known notes.
- Blockers or non-blockers.
- Confirmation that no release was published and no tag was created.

Use [QA Results: 0.6.0-dev](qa-results-0.6.0-dev.md) as the current validation
record style.

## Known Notes

- The previous Go `1.26.3` standard-library advisory gap (GO-2026-5039 and
  GO-2026-5037) is resolved by using Go `1.26.4` or newer.
- `govulncheck` and `gosec` are both blocking in CI and in `make rc-check`.
- Real-terminal interactive TUI QA passed for `v0.6.0-rc.2`.
- Windows is not officially tested or supported yet.
- Batch checksum semantics are unsupported.
- Checksum file discovery and signed checksum verification are not implemented.

## Next Decision

After an RC validation pass, decide whether to:

- Continue internal validation on the same RC.
- Fix findings and create another internal RC tag.
- Continue toward public release readiness work without publishing.

Do not publish a public stable release from this RC. Public stable remains
planned for `v1.0.0`.
