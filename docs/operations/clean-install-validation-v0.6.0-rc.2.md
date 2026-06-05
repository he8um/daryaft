# Clean Install Validation: v0.6.0-rc.2

Date: 2026-06-06

Tag validated: `v0.6.0-rc.2`

Verdict: **PASS WITH NOTES**

## Validation Sources

| Source | Method |
|--------|--------|
| Clean clone | `git clone https://github.com/he8um/daryaft.git /tmp/daryaft-clean-validation` then `git checkout v0.6.0-rc.2` |
| Local build | `make build-local` from clean checkout |
| GoReleaser snapshot | `make release-check` from clean checkout, extracted `daryaft_darwin_arm64.tar.gz` |

## Environment

- Platform: darwin arm64
- Go: go1.26.4
- Temp HOME: `/tmp/daryaft-clean-home`
- HTTP test server: Python `http.server` on `127.0.0.1:18092`
- Test file: `file.txt` (23 bytes, `hello clean validation\n`)
- Test file checksums:
  - SHA-256: `1d0b79cdd94bbbc3fa6efca0ab89c8feaf4a98b8abb4ab28055445849cf712c8`
  - SHA-512: `103cec09c95d1c7ed7b7ff05ccbaebeb41577443edc6a547dc128972decd0c2c7b58b432599e85686f7e54531f9b620869662d674710605f016d2161af7f923d`

## Results

### Clean Checkout

| Check | Result |
|-------|--------|
| Clone from GitHub | PASS |
| `git checkout v0.6.0-rc.2` | PASS — detached HEAD at `5059625` |
| `git describe --tags --always` | `v0.6.0-rc.2` |
| `go version` | `go1.26.4 darwin/arm64` |

### Build and Tests (from clean clone)

| Check | Result |
|-------|--------|
| `go run . version` | PASS — version: `0.6.0-dev`, built_by: source (expected; no ldflags) |
| `go run . version --json` | PASS — valid JSON |
| `go test ./...` | PASS — all packages pass |
| `go build ./...` | PASS |

### Local Build Binary (`make build-local`)

| Check | Result |
|-------|--------|
| Build succeeds | PASS — `bin/daryaft` 13 MB |
| `./bin/daryaft version` | PASS — commit: `5059625`, built_by: make, date injected |
| `./bin/daryaft version --json` | PASS — valid JSON |
| `./bin/daryaft doctor` | PASS WITH NOTE — 1 warning: Downloads dir missing in temp HOME (expected; benign) |
| `./bin/daryaft doctor --json` | PASS — valid JSON, `"ok": true`, `"warnings": 1`, `"failures": 0` |
| `./bin/daryaft doctor --strict` (no Downloads) | Non-zero exit — expected; Downloads dir does not exist in temp HOME |
| `./bin/daryaft doctor --strict` (Downloads created) | PASS — exit 0 once `~/Downloads` exists |

Note: `doctor --strict` exits non-zero in a fresh temp HOME because the Downloads
directory does not exist yet. This is correct behavior — it warns the user to create
it. Once the directory exists (as it would on a real user machine), `--strict` exits 0.

### GoReleaser Snapshot Artifacts (`make release-check`)

| Check | Result |
|-------|--------|
| `make release-check` | PASS — snapshot build succeeded, no publish |
| Archives present | `daryaft_darwin_amd64.tar.gz`, `daryaft_darwin_arm64.tar.gz`, `daryaft_linux_amd64.tar.gz`, `daryaft_linux_arm64.tar.gz` |
| `checksums.txt` present | PASS |
| Extract `daryaft_darwin_arm64.tar.gz` | PASS — contains `CHANGELOG.md`, `LICENSE`, `README.md`, `daryaft` |
| Artifact `daryaft version` | PASS — version: `0.6.0-dev-SNAPSHOT-5059625`, built_by: goreleaser, commit injected |
| Artifact `daryaft version --json` | PASS — valid JSON |
| Artifact `daryaft doctor` | PASS — same expected warning as above |

Note: Snapshot artifact version shows `0.6.0-dev-SNAPSHOT-5059625`. A `goreleaser release`
build on a clean tag (not snapshot) would show the tag version. This is expected behavior
for local validation.

### Inspect

| Check | Result |
|-------|--------|
| `inspect http://127.0.0.1:18092/file.txt` | PASS — URL, filename, content-length, content-type, last-modified all reported |
| `inspect --json` | PASS — valid JSON, all fields present |

### Single Download

| Check | Result |
|-------|--------|
| Download `file.txt` to `-o /tmp/daryaft-clean-output` | PASS |
| File content matches source | PASS — `hello clean validation` |
| No `.part` file after success | PASS |

### Batch Download

| Check | Result |
|-------|--------|
| Batch from `urls.txt` (file.txt + big.bin) | PASS — Total: 2, Completed: 2, Failed: 0 |
| Both files present after batch | PASS |
| No `.part` files after success | PASS |

### Checksum Verification

| Check | Result |
|-------|--------|
| Valid SHA-256 | PASS — `Checksum verified: sha256`, exit 0 |
| Valid SHA-512 | PASS — `Checksum verified: sha512`, exit 0 |
| Mismatch SHA-256 | PASS — error message with expected vs got, exit 1 |
| Final file preserved after mismatch | PASS — `file.txt` present in output dir |

### Dry-Run and Default Output

| Check | Result |
|-------|--------|
| Default output (temp HOME) | PASS — `Output: /tmp/daryaft-clean-home/Downloads` |
| `-o .` | PASS — `Output: .` |
| `DARYAFT_DOWNLOAD_DIR` env override | PASS — `Output: /tmp/daryaft-clean-env-out` |

### Config

| Check | Result |
|-------|--------|
| `config show` (fresh) | PASS — defaults shown |
| `config get download_dir` | PASS — empty (not set) |
| `config set retries 5` | PASS — `Updated config: retries=5` |
| `config get retries` | PASS — `5` |
| `config set theme mono` | PASS — `Updated config: theme=mono` |
| `config get theme` | PASS — `mono` |
| `config set retries 21` (invalid) | PASS — rejected with error, exit 1 |
| `config set theme invalid-theme` (invalid) | PASS — rejected with error, exit 1 |

### Doctor (see also above)

| Check | Result |
|-------|--------|
| `doctor` human output | PASS — readable, structured |
| `doctor --json` parseable | PASS — valid JSON |
| `doctor --strict` with Downloads present | PASS — exit 0 |

### Shell Completion

| Check | Result |
|-------|--------|
| `completion bash` | PASS — 22,482 bytes |
| `completion zsh` | PASS — 7,784 bytes |
| `completion fish` | PASS — 9,785 bytes |
| `completion powershell` | PASS — 10,854 bytes |

### TUI Smoke

Interactive TUI was already validated in a real terminal during `v0.6.0-rc.2`
development validation. Automated/non-interactive TUI smoke is not run here
because a real PTY is not available in this environment.

TUI flows verified in the earlier real-terminal pass:
- Home screen renders.
- Inspect URL flow works.
- Single URL download works.
- Checksum input is present.
- Batch download does not request checksum.
- `q` exits.

### Cleanup

All temporary directories and files removed:
`/tmp/daryaft-clean-*`, `/tmp/_daryaft-clean`

## Notes

- Source and local-build runs show `version: 0.6.0-dev` — expected; ldflags inject
  commit and date from `make build-local` but the version string remains `0.6.0-dev`
  until a tagged release build.
- GoReleaser snapshot artifact shows `0.6.0-dev-SNAPSHOT-5059625` — expected.
  A non-snapshot `goreleaser release` on the `v0.6.0-rc.2` tag would show `v0.6.0-rc.2`.
  A stable `v1.0.0` release build would show `1.0.0`.
- The GitHub pre-release at `v0.6.0-rc.2` contains no binary assets (source/tag-only).
  This validation confirms the source and local-build paths work correctly.
- `doctor --strict` exits non-zero when the Downloads directory does not exist in a
  fresh environment. This is correct behavior and resolves automatically once the user
  has a Downloads directory (universal on macOS and most Linux desktops).

## References

- [v0.6.0-rc.2 Release Status](release-status-v0.6.0-rc.2.md)
- [Release-Candidate Validation](rc-validation.md)
- [Release Readiness: v1.0](../roadmap/release-readiness-v1.0.md)
