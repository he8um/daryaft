# Pre-1.0 Roadmap

Pre-1.0 milestones are for building the stable base. These should be developed locally first and only promoted publicly when safe.

## v0.1.0 Core MVP

Features:

1. `daryaft <url>`
2. Basic file download
3. TUI progress bar
4. Speed, ETA, downloaded/total display
5. Footer attribution

## v0.2.0 Input and File Handling

Features:

1. `-f, --file urls.txt`
2. `-o, --output` directory
3. `--name` filename override
4. Filename detection from `Content-Disposition` and URL

## v0.3.0 Resume and Retry

Features:

1. `.part` files
2. HTTP Range resume
3. Retry with exponential backoff
4. User-friendly network errors

## v0.4.0 TUI Foundation

Features:

1. Interactive mode when running `daryaft` without arguments
2. Simple queue screen
3. Pause/resume in UI
4. Lightweight spinner and retry animations

## v0.5.0 Update System

Features:

1. `daryaft update`
2. `daryaft update --check`
3. Changelog preview
4. Update size preview
5. Checksum verification

## v0.6.0 Packaging Validation

Features:

1. GoReleaser build configuration
2. Archive artifacts
3. `.deb`, `.rpm`, Arch package generation for validation
4. Publishing gates that block public package channels before v1.0.0

## v0.7.0 Documentation and Help

Features:

1. README complete
2. `docs/` complete
3. GitHub Wiki content prepared
4. `-h` help complete for all commands

## v0.8.0 Stability Candidate

Features:

1. Unit tests for downloader core
2. Tests for filename detection
3. Tests for resume logic
4. Tests for update asset selection

## v0.9.0 Release Candidate

Features:

1. UI polish
2. Cross-platform validation
3. Package validation
4. Final docs audit
5. v1.0.0 release checklist pass
