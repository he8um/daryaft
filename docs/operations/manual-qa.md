# Manual QA Checklist

Use this checklist before a local pre-release validation pass. It is local-only:
do not publish releases, create tags, or push changes while running these steps.

See [Pre-Release Readiness](pre-release-readiness.md) for the current
`0.6.0-dev` readiness verdict and known release blockers.

The completed `0.6.0-dev` internal validation readiness pass is recorded in
[QA Results: 0.6.0-dev](qa-results-0.6.0-dev.md).
For RC tag and artifact validation, use
[Release-Candidate Validation](rc-validation.md) and the
[Daryaft v0.6.0-rc.2 Internal Release Candidate](release-notes-v0.6.0-rc.2.md)
notes. `v0.6.0-rc.1` is superseded.

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

For a release-candidate validation pass, run:

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
- `make rc-check` passes including blocking `govulncheck` and `gosec` checks.
- `make security` passes with Go `1.26.4` or newer; `govulncheck` is blocking
  in CI. The previous Go `1.26.3` standard-library advisory gap is resolved.
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

For batch/per-target checksums, use `--checksum-file`:

```bash
a_sum="$(shasum -a 256 /tmp/daryaft-qa-server/file.txt | awk '{print $1}')"
b_sum="$(shasum -a 256 /tmp/daryaft-qa-server/big.bin | awk '{print $1}')"
cat > /tmp/daryaft-qa-checksums.txt <<EOF
sha256:$a_sum http://localhost:<port>/file.txt
sha256:$b_sum http://localhost:<port>/big.bin
EOF
go run . http://localhost:<port>/file.txt http://localhost:<port>/big.bin \
  -o /tmp/daryaft-qa-out --checksum-file /tmp/daryaft-qa-checksums.txt
```

Expected:

- Exit 0 with a `Checksum verified: 2` line in the batch summary.
- A manifest with a wrong digest fails that item and exits non-zero, leaving the
  file in place.
- A manifest missing a target, with an extra URL, or with a duplicate/malformed
  line exits non-zero before the network download starts, with a
  `checksum file:` error (line-numbered for malformed lines).
- `--checksum` together with `--checksum-file` exits non-zero with
  `--checksum and --checksum-file cannot be used together`.
- See [Checksum Verification QA](checksum-verification-qa.md) for the full
  batch checksum checklist.

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

## HTTP Request Customization Smoke Tests

These are quick sanity checks. Full coverage is in
[HTTP Customization QA](http-customization-qa.md).

```bash
go run . download --help | grep -E '^\s+--proxy|--header|--user-agent|--username|--password'
go run . inspect --help | grep -E '^\s+--proxy|--header|--user-agent|--username|--password'
```

Expected: all five flags appear in each command's help output.

```bash
go run . download http://localhost:<port>/file.txt --header "NoColon" --dry-run
go run . inspect http://localhost:<port>/file.txt --header "NoColon"
```

Expected: clear error about invalid header format, no download or inspection.

```bash
go run . download http://localhost:<port>/file.txt \
  --username alice \
  --password topsecret \
  --dry-run
```

Expected: dry-run output does not contain `topsecret`; contains `[REDACTED]`.

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

Input UX Polish (v1.12.0):

- Open the URL input screen. Verify the prompt says `Enter a download URL
  (https:// or http://)`. Verify a defaults preview line appears showing the
  configured save directory, retries, and resume value.
- Press Enter without typing. Verify an inline error appears mentioning
  `https://` as a guidance example.
- Type `ftp://example.com/file.zip` and press Enter. Verify an inline error
  appears mentioning `scheme must be http or https`. Verify no plan screen
  appears.
- Type a character after the error. Verify the error clears immediately.
- Open the file input screen. Verify the prompt mentions absolute path and
  one-URL-per-line format. Verify a defaults preview line appears.
- Press Enter without typing. Verify an inline error appears mentioning `.txt`
  as guidance.
- Open the Help screen. Verify it mentions the `c` shortcut for Settings.

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
- Empty URL or empty file path shows guidance, not a raw Go error.
- FTP and other non-HTTP schemes are rejected with a clear inline message.
- Defaults preview is visible on both URL and file input screens.

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

## v1.6.0 Download + TUI UX Polish QA

1. Single-URL success: verify `Completed:` message shows size and elapsed time.
2. Single-URL failure: verify `Failed:` prefix appears.
3. Resume: verify `Resuming:` prefix appears when server returns 206.
4. Restart: verify `Restarting:` prefix appears when server does not support Range.
5. Batch summary: verify `Not started` (not `Skipped`) for unstarted items.
6. TUI batch: run a multi-URL batch, verify queue list with ✓/✗/→ markers.
7. TUI no-color (`--no-color`): verify ASCII markers `[ok]`, `[!]`, `[>]`.
8. TUI post-run hint: verify hint says `new download` after batch completes.

## v1.9.0 Download Reliability QA

These checks exercise the retry/resume reliability hardening. The deterministic
behavior is covered by `internal/downloader` tests; this is a quick manual
sanity pass.

1. 408 retry then success: point Daryaft at a test endpoint that returns
   `408 Request Timeout` once and then `200 OK`. Verify a `Retrying` line
   appears and the download then completes.
2. Partial larger than remote restarts safely: create a `.part` file larger
   than the remote file, then download with resume enabled. Verify Daryaft
   prints `Restarting: Partial file is larger than remote file; restarting
   download` and the final file matches the remote content exactly.
3. Missing/corrupt sidecar restarts safely: create a `.part` file with no
   `.part.daryaft.json` (or a corrupt one), then download with resume enabled.
   Verify no panic and the final file is correct.
4. Cancellation exits non-zero: cancel an active CLI download with Ctrl+C and
   verify the process exits with code `1`, keeps the `.part` file, and does not
   rename to the final target.

## v1.10.0 Config Safe Core Gap-Closing QA

These checks exercise the new `user_agent`, `timeout`, and `--config` additions.

1. **`--config` explicit path loads correctly**: run
   `daryaft --config /tmp/test.yaml config init` then
   `daryaft --config /tmp/test.yaml config set retries 9` then
   `daryaft --config /tmp/test.yaml config get retries`. Verify output is `9`.
   Run `daryaft config get retries` (without `--config`). Verify output is `3`
   (default), confirming the files are independent.

2. **`--config` missing file errors**: run
   `daryaft --config /tmp/does-not-exist.yaml config show`.
   Verify output is `config file not found: /tmp/does-not-exist.yaml` and exit
   code is non-zero.

3. **`user_agent` config default**: run
   `daryaft --config /tmp/test.yaml config set user_agent "QABot/1.0"`.
   Verify `daryaft --config /tmp/test.yaml config get user_agent` returns
   `QABot/1.0`. Run `daryaft --config /tmp/test.yaml config show` and verify
   `user_agent: QABot/1.0` appears in the YAML output.

4. **`DARYAFT_USER_AGENT` overrides config**: with `user_agent: QABot/1.0` in
   config, run `DARYAFT_USER_AGENT=EnvBot/2.0 daryaft --config /tmp/test.yaml config get user_agent`.
   Verify output is `EnvBot/2.0`.

5. **`timeout` config and `--timeout` flag**: run
   `daryaft --config /tmp/test.yaml config set timeout 45s`.
   Verify `daryaft --config /tmp/test.yaml config get timeout` returns `45s`.
   Run a real download with `--timeout 5s` and verify it completes without
   error on a fast server.

6. **Invalid timeout rejected**: run `daryaft config set timeout abc` and
   `daryaft config set timeout 0s`. Verify both return a clear error with
   `timeout` in the message.

7. **Secret fields rejected by strict YAML**: write a config file containing
   `username: test`. Run `daryaft --config /path/to/that.yaml config show`.
   Verify Daryaft returns a `parse config` error (strict `KnownFields` enforcement).

8. **Cleanup**: run `daryaft --config /tmp/test.yaml config reset` to restore
   defaults at the test path, then delete `/tmp/test.yaml`.

## v1.11.0 — TUI Settings Screen

1. **Settings screen — no config file**: start `daryaft` with a temp home dir
   where no config file exists. Open Settings from the home menu. Verify:
   - `Config loaded: no (using defaults)` is shown.
   - All safe config keys are present (`download_dir`, `retries`, `resume`, etc.).
   - No password, token, authorization, cookie, or proxy field appears.

2. **Settings screen — with config file**: run `daryaft config init` to create a
   config file, then start `daryaft`. Open Settings. Verify:
   - `Config loaded: yes` is shown.
   - The correct config path appears.

3. **Settings screen — explicit `--config`**: run `daryaft --config /tmp/qa.yaml config init`,
   then `daryaft --config /tmp/qa.yaml`. Open Settings. Verify `/tmp/qa.yaml` is shown as the
   config path.

4. **Settings screen — user_agent and timeout**: run
   `daryaft config set user_agent QABot/1.11` and
   `daryaft config set timeout 30s`, then start `daryaft`. Open Settings. Verify
   `user_agent: QABot/1.11` and `timeout: 30s` appear.

5. **Settings screen — defaults markers**: with empty `user_agent` and `timeout`,
   verify `user_agent: (default)` and `timeout: (none)`.

6. **`c` shortcut**: from the TUI home screen, press `c`. Verify Settings screen opens.
   Press `esc` and verify home screen returns.

7. **`c` inside input screens**: enter the URL input screen, type `c`. Verify it appends
   `c` to the input field and does not open the Settings screen.

8. **No secrets**: verify Settings does not show any of: `password`, `token`,
   `authorization`, `cookie`, `proxy_authorization`, or values of
   `DARYAFT_USERNAME` / `DARYAFT_PASSWORD`.

9. **Cleanup**: delete `/tmp/qa.yaml` after QA.
