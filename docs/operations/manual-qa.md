# Manual QA Checklist

Use this checklist before a local pre-release validation pass. It is local-only:
do not publish releases, create tags, or push changes while running these steps.

See [Pre-Release Readiness](pre-release-readiness.md) for the current
`0.6.0-dev` readiness verdict and known release blockers.

The completed `0.6.0-dev` internal validation readiness pass is recorded in
[QA Results: 0.6.0-dev](qa-results-0.6.0-dev.md).

Run all commands from the repository root unless a step says otherwise.

## Prerequisites

- Go installed and available in `PATH`.
- Python 3 available for the local HTTP test server.
- GoReleaser v2 installed if running `goreleaser check` or `make release-check`.
- `golangci-lint`, `govulncheck`, and `gosec` installed for local lint and
  security checks. These are optional but recommended.
- A terminal capable of running the TUI.

## Clean Workspace

Confirm that the checkout starts clean and that the expected recent commits are
present:

```bash
git status
git log --oneline --decorate --graph --all -10
```

Expected:

- `git status` reports a clean working tree before QA begins.
- Recent history includes the latest pushed work for the TUI inspect flow.

## Quality Gates

Run the standard local checks:

```bash
go test ./...
go build ./...
go test -race ./internal/downloader
go test -race ./internal/tui
```

Run the optional tool-backed gates when the tools are installed:

```bash
make lint
govulncheck ./...
gosec ./...
make security
goreleaser check
```

For a release-candidate validation pass while Go 1.26.x can still resolve to
Go 1.26.3 in tooling, run:

```bash
make rc-check
```

If GoReleaser is installed and snapshot validation is desired, run:

```bash
make release-check
```

Expected:

- All required Go tests and builds pass.
- Optional lint, `gosec`, and GoReleaser checks pass when their tools are
  installed.
- `make rc-check` passes without running `govulncheck`.
- Local `make security` remains strict. Until local Go 1.26.x resolves to Go
  1.26.4 or newer, `govulncheck` may report known Go 1.26.3 standard-library
  vulnerabilities GO-2026-5039 and GO-2026-5037; these are temporarily advisory
  only in CI and should be reverted to blocking there once patched Go is
  available.
- `make release-check` writes only local snapshot artifacts under ignored build
  directories such as `dist/`; it must not publish.

## Local HTTP Server

Use a separate terminal for the local HTTP server. The helper script prepares
test files and starts a server on `127.0.0.1:8091`:

```bash
scripts/manual-qa-server.sh
```

Manual setup, if you do not use the helper script:

```bash
rm -rf /tmp/daryaft-qa-server /tmp/daryaft-qa-out
mkdir -p /tmp/daryaft-qa-server /tmp/daryaft-qa-out
printf 'Daryaft manual QA test file\n' > /tmp/daryaft-qa-server/file.txt
python3 - <<'PY'
from pathlib import Path
path = Path("/tmp/daryaft-qa-server/big.bin")
chunk = bytes(range(256)) * 4096
with path.open("wb") as f:
    for _ in range(24):
        f.write(chunk)
PY
cat > /tmp/daryaft-qa-server/urls.txt <<'EOF'
http://localhost:8091/file.txt
http://localhost:8091/big.bin
EOF
cd /tmp/daryaft-qa-server
python3 -m http.server 8091 --bind 127.0.0.1
```

Use another terminal for all Daryaft commands. If you use a different port,
replace `<port>` in the steps below.

## Default Output Smoke

```bash
go run . http://localhost:<port>/file.txt --dry-run
go run . http://localhost:<port>/file.txt --dry-run -o .
DARYAFT_DOWNLOAD_DIR=/tmp/daryaft-env-out go run . http://localhost:<port>/file.txt --dry-run
```

Expected:

- With no output flag, environment value, or config value, `Output:` is
  `~/Downloads` for the current user.
- With `-o .`, `Output:` is `.`.
- With `DARYAFT_DOWNLOAD_DIR=/tmp/daryaft-env-out`, `Output:` is
  `/tmp/daryaft-env-out`.
- In the TUI, the output directory input starts with the same effective output
  directory. Clearing it and pressing enter keeps that effective default; typing
  `.` uses the current directory explicitly.

## CLI Single Download

```bash
rm -rf /tmp/daryaft-qa-out
mkdir -p /tmp/daryaft-qa-out
go run . http://localhost:<port>/file.txt -o /tmp/daryaft-qa-out
ls -la /tmp/daryaft-qa-out
```

Expected:

- `file.txt` downloads successfully.
- `/tmp/daryaft-qa-out/file.txt` exists.
- No `file.txt.part` remains after success.

## CLI Batch Download

```bash
rm -rf /tmp/daryaft-qa-out
mkdir -p /tmp/daryaft-qa-out
go run . -f /tmp/daryaft-qa-server/urls.txt -o /tmp/daryaft-qa-out
ls -la /tmp/daryaft-qa-out
```

Expected:

- All files listed in `urls.txt` download.
- The batch summary shows the completed count.
- No successful download leaves a `.part` file behind.

## CLI Dry Run

For checksum dry-run, compute the expected SHA-256 first:

```bash
file_sha256="$(shasum -a 256 /tmp/daryaft-qa-server/file.txt | awk '{print $1}')"
```

```bash
rm -rf /tmp/daryaft-qa-out
mkdir -p /tmp/daryaft-qa-out
go run . http://localhost:<port>/file.txt -o /tmp/daryaft-qa-out --dry-run
go run . http://localhost:<port>/file.txt -o /tmp/daryaft-qa-out --checksum "sha256:${file_sha256}" --dry-run
ls -la /tmp/daryaft-qa-out
go run . -f /tmp/daryaft-qa-server/urls.txt -o /tmp/daryaft-qa-out --dry-run
ls -la /tmp/daryaft-qa-out
```

Expected:

- Each command prints a dry-run plan.
- The checksum dry-run shows `Checksum: sha256:<file_sha256>`.
- No files are written to `/tmp/daryaft-qa-out`.

## CLI Checksum

Compute expected checksums for the local fixture, then verify successful
single URL downloads:

```bash
expected_sha256="$(shasum -a 256 /tmp/daryaft-qa-server/file.txt | awk '{print $1}')"
expected_sha512="$(shasum -a 512 /tmp/daryaft-qa-server/file.txt | awk '{print $1}')"
rm -rf /tmp/daryaft-qa-out
mkdir -p /tmp/daryaft-qa-out
go run . http://localhost:<port>/file.txt -o /tmp/daryaft-qa-out --checksum "sha256:${expected_sha256}"
rm -f /tmp/daryaft-qa-out/file.txt
go run . http://localhost:<port>/file.txt -o /tmp/daryaft-qa-out --checksum "sha512:${expected_sha512}"
```

Expected:

- The command prints `Checksum verified: sha256`.
- The SHA-512 command prints `Checksum verified: sha512`.
- `/tmp/daryaft-qa-out/file.txt` exists.
- No `file.txt.part` remains after success.

Run a mismatch check:

```bash
rm -rf /tmp/daryaft-qa-out
mkdir -p /tmp/daryaft-qa-out
go run . http://localhost:<port>/file.txt -o /tmp/daryaft-qa-out --checksum "sha256:0000000000000000000000000000000000000000000000000000000000000000"
ls -la /tmp/daryaft-qa-out
```

Expected:

- The command exits non-zero.
- The error includes `checksum mismatch: expected`.
- The completed final file remains in `/tmp/daryaft-qa-out`.
- The mismatch does not delete the final file in this milestone.

Run invalid checksum and unsupported batch/file checks:

```bash
go run . http://localhost:<port>/file.txt -o /tmp/daryaft-qa-out --checksum "sha256:not-hex"
go run . http://localhost:<port>/file.txt http://localhost:<port>/big.bin -o /tmp/daryaft-qa-out --checksum "sha256:${expected_sha256}"
go run . -f /tmp/daryaft-qa-server/urls.txt -o /tmp/daryaft-qa-out --checksum "sha256:${expected_sha256}"
```

Expected:

- Invalid checksum exits non-zero before the network download starts.
- Batch and `--file` checksum commands exit non-zero.
- Batch and `--file` errors include
  `--checksum is currently supported only for single URL downloads`.
- TUI checksum behavior is covered separately in the TUI flow section.

## CLI Resume

Use the larger file so there is time to interrupt:

```bash
rm -rf /tmp/daryaft-qa-out
mkdir -p /tmp/daryaft-qa-out
go run . http://localhost:<port>/big.bin -o /tmp/daryaft-qa-out
```

Press `Ctrl+C` while the download is active, then check the partial state:

```bash
ls -la /tmp/daryaft-qa-out
```

Rerun the same command:

```bash
go run . http://localhost:<port>/big.bin -o /tmp/daryaft-qa-out
```

Expected:

- After interruption, `.part` and `.part.daryaft.json` files exist.
- Rerunning the same command resumes when the server supports ranges.
- If the server does not support ranges, Daryaft restarts safely.
- After success, `big.bin` exists and partial files are removed.

## CLI Safe Cancellation

Start a large or slow enough download:

```bash
rm -rf /tmp/daryaft-qa-out
mkdir -p /tmp/daryaft-qa-out
go run . http://localhost:<port>/big.bin -o /tmp/daryaft-qa-out
```

Press `Ctrl+C` during the transfer.

Expected:

- Daryaft prints a cancellation message.
- The command exits non-zero.
- The partial download is kept.
- The final `big.bin` file is not created.

## CLI Inspect

```bash
go run . inspect http://localhost:<port>/file.txt
go run . inspect http://localhost:<port>/file.txt --json
```

Optional JSON parse check:

```bash
go run . inspect http://localhost:<port>/file.txt --json | python3 -m json.tool
```

Expected:

- Metadata is printed.
- No files are written.
- JSON output is valid JSON.

## TUI Flows

Start the TUI:

```bash
go run .
```

Test these flows:

- Download from URL: `http://localhost:<port>/file.txt`.
- Download from URL with empty checksum.
- Download from URL with valid checksum:
  `sha256:<file_sha256>`.
- Download from URL with invalid checksum:
  `sha256:not-hex`.
- Download from `.txt` file: `/tmp/daryaft-qa-server/urls.txt`.
- Inspect URL: `http://localhost:<port>/file.txt`.
- Press `q` during an active download.
- Press `q` during an active inspect.
- Use Backspace in input fields, including on an empty input.
- Resize the terminal while screens are active, if practical.

Expected:

- No panic.
- Screens navigate correctly.
- Empty checksum reaches the plan with `Checksum: none`.
- Valid checksum reaches the plan, then downloads successfully.
- Invalid checksum stays on the checksum input screen and does not start a
  download.
- The `.txt` batch flow does not ask for checksum.
- Inspect does not write files.
- A cancelled download keeps the partial file.

## Config And Environment

These commands update the local Daryaft config file. Before testing, record the
path so you can restore it if needed:

```bash
go run . config path
go run . config show
```

Run:

```bash
go run . config set retries 5
go run . config set theme mono
go run . config show
go run . config set retries 21
go run . config set theme invalid-theme
DARYAFT_RETRIES=4 go run . http://localhost:<port>/file.txt --dry-run
```

Expected:

- `config show` reflects `retries: 5` and `theme: mono`.
- Invalid `retries` values greater than `20` are rejected.
- Invalid theme values are rejected.
- The `DARYAFT_RETRIES=4` dry-run uses the environment override without writing
  files.

## Doctor

```bash
go run . doctor
go run . doctor --json
go run . doctor --json | python3 -m json.tool
go run . doctor --strict
```

Expected:

- Plain output is readable.
- JSON output is parseable.
- `--strict` returns non-zero when warnings or failures are present and returns
  zero when no warnings or failures are present.

## Completion

```bash
go run . completion bash > /tmp/daryaft-qa-completion.bash
go run . completion zsh > /tmp/daryaft-qa-completion.zsh
go run . completion fish > /tmp/daryaft-qa-completion.fish
go run . completion powershell > /tmp/daryaft-qa-completion.ps1
wc -c /tmp/daryaft-qa-completion.*
```

Expected:

- Each completion command exits successfully.
- Each generated completion file is non-empty.

## Release Readiness Local Only

```bash
goreleaser check
make release-check
```

Expected:

- `goreleaser check` passes.
- `make release-check` creates local snapshot artifacts only.
- No GitHub release is published.
- No Git tags are created.

## Cleanup

Stop the local HTTP server with `Ctrl+C`, then remove temporary QA artifacts:

```bash
rm -rf /tmp/daryaft-qa-server /tmp/daryaft-qa-out
rm -f /tmp/daryaft-qa-completion.bash
rm -f /tmp/daryaft-qa-completion.zsh
rm -f /tmp/daryaft-qa-completion.fish
rm -f /tmp/daryaft-qa-completion.ps1
git status
```

Expected:

- The HTTP server is stopped.
- Temporary test directories and completion files are removed.
- `git status` is clean except for intentional documentation or script changes
  from this checklist update.
