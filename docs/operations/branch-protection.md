# Branch Protection

These are recommended GitHub branch protection settings for `main`. They are
not enforced by repository files; configure them in GitHub repository settings.

Daryaft is pre-1.0. Public stable install and release channels remain planned
for `v1.0.0`.

## Recommended Settings

- Require a pull request before merging when the project accepts external
  contributors.
- Require branches to be up to date before merging.
- Restrict force pushes.
- Require linear history if the maintainers choose that workflow.
- Do not require signed commits yet unless the project intentionally adopts
  that policy.

## Required Status Checks

Require these checks before merging to `main`:

- `Go test/build (ubuntu-latest)`
- `Go test/build (macos-latest)`
- `goreleaser-check`
- `lint`
- `security`

The Go test/build matrix verifies module tidiness, `go test ./...`,
`go build ./...`, and `go test -race ./internal/tui` on Linux and macOS.

The `goreleaser-check` job validates `.goreleaser.yml` with `goreleaser check`
only. It does not publish releases, create tags, run snapshot builds, or use
publishing secrets.

The `lint` job runs blocking `golangci-lint run` using `.golangci.yml`. The
`security` job runs `govulncheck ./...` and `gosec ./...`; `gosec` remains
blocking. `govulncheck` is temporarily advisory in CI because Go 1.26.x on
hosted tooling can still resolve to Go 1.26.3 and report standard-library
vulnerabilities GO-2026-5039 and GO-2026-5037, both fixed in Go 1.26.4. Restore
`govulncheck` to blocking once CI can use Go 1.26.4 or newer. These jobs do not
publish releases, create tags, run snapshot builds, or use publishing secrets.
They may use a newer Go toolchain than the test/build matrix only to install
current quality tools; runtime compatibility remains governed by `go.mod`.

## Release Safety

Branch protection should not be used as a replacement for release policy. Before
`v1.0.0`, do not publish public releases, create release tags, or enable
Homebrew, deb, rpm, or Arch publishing.
