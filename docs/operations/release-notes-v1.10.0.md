# Daryaft v1.10.0

Daryaft v1.10.0 is a stable release focused on safe configuration persistence improvements.

## Highlights

- Added `user_agent` configuration support for default download user-agent.
- Added `DARYAFT_USER_AGENT` environment override.
- Added `timeout` configuration support for overall download request timeout.
- Added `DARYAFT_TIMEOUT` environment override.
- Added `--timeout <duration>` download flag.
- Added global `--config <path>` flag for explicit configuration file selection.
- Added tests for config/env/CLI precedence.
- Updated configuration, command reference, usage, manual QA, roadmap, and README documentation.
- Preserved the security boundary: credentials and secret-bearing values are not persisted in config.

## Configuration

Daryaft already supports YAML configuration at the platform-specific config path:

```text
<UserConfigDir>/daryaft/config.yaml
```

Examples:

```text
macOS: ~/Library/Application Support/daryaft/config.yaml
Linux: ~/.config/daryaft/config.yaml
```

This release adds two safe non-secret keys:

```yaml
user_agent: ""
timeout: ""
```

`user_agent` sets the default download user-agent when `--user-agent` is not provided.

`timeout` sets the overall HTTP request timeout when `--timeout` is not provided. It uses Go duration syntax such as:

```text
30s
2m
1m30s
```

## Environment overrides

New environment variables:

```text
DARYAFT_USER_AGENT
DARYAFT_TIMEOUT
```

Precedence remains:

```text
CLI flags > environment variables > config file > built-in defaults
```

## New flags

Global config path override:

```bash
daryaft --config ./config.yaml config show
daryaft --config ./config.yaml download https://example.com/file.zip
```

Download timeout override:

```bash
daryaft download https://example.com/file.zip --timeout 30s
```

## Security model

Daryaft v1.10.0 intentionally does not persist secrets.

The following remain unsupported in config:

* usernames
* passwords
* tokens
* API keys
* Authorization headers
* Proxy-Authorization headers
* cookies
* proxy URLs
* arbitrary headers

Use existing CLI flags or environment variables for per-invocation credentials where supported.

## Validation

* Invalid timeout values are rejected.
* Zero or negative timeout values are rejected.
* Invalid user-agent values are rejected.
* Missing default config file remains non-error.
* Missing explicit `--config <path>` is an error.
* Unknown config fields remain rejected by strict YAML parsing.

## Documentation

Updated or added:

* `docs/configuration.md`
* `docs/command-reference.md`
* `docs/usage.md`
* `docs/operations/manual-qa.md`
* `docs/roadmap/v1.10.0-config-persistence-safe-core.md`
* `docs/index.md`
* `README.md`
* `CHANGELOG.md`

## Known limitations

* Credentials are intentionally not persisted.
* Auth headers are intentionally not persisted.
* Proxy persistence is intentionally deferred.
* Arbitrary header persistence is intentionally deferred.
* Config profiles are not implemented.
* TUI config editor is not implemented.
* Auto-update is still not implemented.
* GoReleaser Homebrew publishing remains disabled.
* Windows is not officially supported.

## Upgrade

Homebrew:

```bash
brew update && brew upgrade daryaft
```

GitHub binary users:

Download the matching archive from GitHub Releases.

Source users:

```bash
git pull
go build .
```
