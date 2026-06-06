# Self-Update Roadmap

Daryaft's self-update story is being built incrementally to avoid unsafe
in-place binary replacement before the release infrastructure is trusted.

## Current State (v1.1.0)

`daryaft update --check` is implemented and live.

- Queries the GitHub Releases API.
- Compares the current version to the latest stable release.
- Reports status with install-channel-aware upgrade instructions.
- Is **read-only**: does not download, install, or replace the current binary.

```bash
daryaft update --check
daryaft update --check --json
daryaft update --check --include-prerelease
```

`daryaft update` without `--check` exits non-zero and directs the user to
`daryaft update --check`.

## Install-Channel Awareness

| Channel | Update command suggested |
|---------|--------------------------|
| `homebrew` | `brew update && brew upgrade daryaft` |
| `goreleaser` | GitHub Releases URL |
| `source` | GitHub Releases URL |
| `unknown` | GitHub Releases URL |

Homebrew users should always upgrade through Homebrew rather than replacing
the binary manually, as Homebrew manages symlinks, cellar cleanup, and rollback.

## Planned: Auto-Update

`daryaft update` (without `--check`) is a placeholder for a future auto-update
command. It is not yet implemented.

Before auto-update can be implemented safely, the following must be in place:

1. **Stable release channel**: the Homebrew tap (`he8um/homebrew-tap`) is live.
   Direct binary download is also available. Both channels are working for v1.0.0.

2. **Binary replacement safety**: replacing a running binary in-place requires
   OS-specific care (write to a temp path, rename, validate signature/checksum).
   The replacement must not run untrusted code.

3. **Checksum verification**: downloaded binaries must be verified against
   `checksums.txt` from the official GitHub release before replacing the active
   binary.

4. **Signature verification**: signed checksums or release asset signatures
   should be verified before execution. Not yet implemented.

5. **Channel-aware update path**: Homebrew-installed binaries should not be
   replaced directly — they should upgrade through Homebrew to avoid breaking
   the Cellar. Auto-update for Homebrew installations may therefore just invoke
   `brew upgrade daryaft` rather than downloading a binary.

6. **Rollback**: if the new binary fails a smoke test, the old binary should be
   restored.

7. **CI and platform testing**: auto-update must be tested on Linux and macOS
   and must handle permission errors, disk-full conditions, and concurrent
   invocations safely.

## Non-Goals

- Auto-update will **not** require a GitHub auth token.
- Auto-update will **not** modify system paths outside the install prefix.
- Auto-update will **not** be invoked silently or on a background schedule.
- Auto-update will **not** replace the binary without user confirmation, at
  least initially.

## References

- [Homebrew Tap](../operations/homebrew-tap.md)
- [v1.0.0 Release Assets](../operations/release-assets.md)
- [Post-1.0 Feature Packs](post-1-feature-packs.md)
- [Versioning Policy](versioning-policy.md)
