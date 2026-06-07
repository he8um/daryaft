#!/usr/bin/env bash
# release-preflight.sh — Read-only guardrail that validates a target release
# version before a human creates a tag.
#
# Usage:
#   scripts/release-preflight.sh 1.5.0
#   scripts/release-preflight.sh v1.5.0
#   scripts/release-preflight.sh 1.6.0 --allow-skip
#
# This script NEVER creates tags, publishes releases, uploads assets, modifies
# the Homebrew tap, runs git push, or changes any file.

set -euo pipefail

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BOLD='\033[1m'
RESET='\033[0m'

pass()  { echo -e "  ${GREEN}✓${RESET} $*"; }
fail()  { echo -e "  ${RED}✗${RESET} $*"; FAILURES=$((FAILURES + 1)); }
warn()  { echo -e "  ${YELLOW}!${RESET} $*"; WARNINGS=$((WARNINGS + 1)); }
info()  { echo -e "    $*"; }

FAILURES=0
WARNINGS=0

# ---------------------------------------------------------------------------
# Argument parsing
# ---------------------------------------------------------------------------

TARGET_RAW=""
ALLOW_SKIP=0

for arg in "$@"; do
  case "$arg" in
    --allow-skip) ALLOW_SKIP=1 ;;
    --*)
      echo "Unknown flag: $arg" >&2
      echo "Usage: $0 <version> [--allow-skip]" >&2
      exit 1
      ;;
    *) TARGET_RAW="$arg" ;;
  esac
done

if [[ -z "$TARGET_RAW" ]]; then
  echo "Error: version argument is required." >&2
  echo "Usage: $0 <version> [--allow-skip]" >&2
  exit 1
fi

# Strip leading 'v'
TARGET="${TARGET_RAW#v}"

# ---------------------------------------------------------------------------
# Version format validation
# ---------------------------------------------------------------------------

# Must be X.Y.Z with numeric components only
if ! [[ "$TARGET" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo -e "${RED}Error:${RESET} '$TARGET_RAW' is not a valid release version (must be X.Y.Z, e.g. 1.5.0)." >&2
  exit 1
fi

# Reject dev/prerelease suffixes (the regex above already ensures this, but be explicit)
if [[ "$TARGET_RAW" == *"-dev"* || "$TARGET_RAW" == *"-alpha"* || "$TARGET_RAW" == *"-beta"* || "$TARGET_RAW" == *"-rc"* ]]; then
  echo -e "${RED}Error:${RESET} '$TARGET_RAW' looks like a prerelease or dev version. Only stable X.Y.Z versions are accepted." >&2
  exit 1
fi

TARGET_TAG="v${TARGET}"

IFS='.' read -r T_MAJOR T_MINOR T_PATCH <<< "$TARGET"

# ---------------------------------------------------------------------------
# Header
# ---------------------------------------------------------------------------

echo ""
echo -e "${BOLD}Release preflight${RESET}"
echo ""
echo -e "  Target: ${BOLD}${TARGET_TAG}${RESET}"
echo ""

# ---------------------------------------------------------------------------
# Check: inside a git repository
# ---------------------------------------------------------------------------

if ! git rev-parse --git-dir > /dev/null 2>&1; then
  echo -e "${RED}Error:${RESET} Not inside a git repository." >&2
  exit 1
fi

# ---------------------------------------------------------------------------
# Check: working tree clean
# ---------------------------------------------------------------------------

if [[ -n "$(git status --porcelain)" ]]; then
  fail "Working tree is dirty. Commit or stash all changes before release."
else
  pass "Working tree is clean."
fi

# ---------------------------------------------------------------------------
# Check: on main branch
# ---------------------------------------------------------------------------

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "$CURRENT_BRANCH" != "main" ]]; then
  fail "Current branch is '$CURRENT_BRANCH', expected 'main'."
else
  pass "Branch: main"
fi

# ---------------------------------------------------------------------------
# Check: not behind origin/main
# ---------------------------------------------------------------------------

git fetch origin main --quiet 2>/dev/null || true
LOCAL_SHA="$(git rev-parse HEAD)"
REMOTE_SHA="$(git rev-parse origin/main 2>/dev/null || echo '')"

if [[ -z "$REMOTE_SHA" ]]; then
  warn "Could not resolve origin/main. Skipping behind-check."
elif [[ "$LOCAL_SHA" != "$REMOTE_SHA" ]]; then
  BEHIND="$(git rev-list HEAD..origin/main --count 2>/dev/null || echo '?')"
  fail "Local main is behind origin/main by $BEHIND commit(s). Run: git pull --ff-only"
else
  pass "Local main is up to date with origin/main."
fi

# ---------------------------------------------------------------------------
# Check: source development version matches target
# ---------------------------------------------------------------------------

# Read directly from pkg/version/version.go to avoid needing go run
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
VERSION_FILE="$REPO_ROOT/pkg/version/version.go"

if [[ ! -f "$VERSION_FILE" ]]; then
  fail "Cannot find $VERSION_FILE to read source version."
  SOURCE_DEV_VERSION="(unknown)"
else
  SOURCE_DEV_VERSION="$(grep -E 'Version\s*=' "$VERSION_FILE" | head -1 | grep -oE '"[^"]+"' | tr -d '"')"
fi

EXPECTED_DEV="${TARGET}-dev"
echo -e "  Source dev version: ${BOLD}${SOURCE_DEV_VERSION}${RESET}"

if [[ "$SOURCE_DEV_VERSION" == "$EXPECTED_DEV" ]]; then
  pass "Source dev version matches target ($EXPECTED_DEV)."
elif [[ "$SOURCE_DEV_VERSION" == "$TARGET" ]]; then
  warn "Source dev version is already '$TARGET' (without -dev suffix). Was the dev bump skipped?"
else
  fail "Source dev version is '$SOURCE_DEV_VERSION', expected '${EXPECTED_DEV}'. Update pkg/version/version.go before releasing."
fi

# ---------------------------------------------------------------------------
# Check: latest stable tag and skip detection
# ---------------------------------------------------------------------------

# Collect semver tags (vX.Y.Z), ignore rc/alpha/beta
ALL_STABLE_TAGS="$(git tag --list "v*.*.*" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -t. -k1,1V -k2,2n -k3,3n || true)"
LATEST_STABLE_TAG="$(echo "$ALL_STABLE_TAGS" | tail -1)"

if [[ -z "$LATEST_STABLE_TAG" ]]; then
  warn "No stable tags found. This would be the first release."
  LATEST_STABLE_TAG="(none)"
  echo -e "  Latest stable tag: ${BOLD}(none)${RESET}"
else
  echo -e "  Latest stable tag: ${BOLD}${LATEST_STABLE_TAG}${RESET}"
fi

# Skip detection: compare target to expected next version
if [[ "$LATEST_STABLE_TAG" != "(none)" ]]; then
  LATEST_VER="${LATEST_STABLE_TAG#v}"
  IFS='.' read -r L_MAJOR L_MINOR L_PATCH <<< "$LATEST_VER"

  # Expected next: bump minor (same major), patch stays 0; or bump patch
  EXPECTED_NEXT_MINOR="${L_MAJOR}.$((L_MINOR + 1)).0"
  EXPECTED_NEXT_PATCH="${L_MAJOR}.${L_MINOR}.$((L_PATCH + 1))"

  if [[ "$TARGET" == "$EXPECTED_NEXT_MINOR" || "$TARGET" == "$EXPECTED_NEXT_PATCH" ]]; then
    pass "Target $TARGET_TAG is the expected next version after $LATEST_STABLE_TAG."
  elif [[ "$T_MAJOR" -gt "$L_MAJOR" ]]; then
    if [[ "$ALLOW_SKIP" -eq 1 ]]; then
      warn "Target $TARGET_TAG is a major-version jump from $LATEST_STABLE_TAG (--allow-skip set)."
    else
      fail "Target $TARGET_TAG is a major-version jump from $LATEST_STABLE_TAG. Use --allow-skip if intentional."
    fi
  else
    # Target is ahead by more than one minor/patch step
    if [[ "$ALLOW_SKIP" -eq 1 ]]; then
      warn "Target $TARGET_TAG skips past expected next version ($EXPECTED_NEXT_MINOR or $EXPECTED_NEXT_PATCH) from $LATEST_STABLE_TAG. --allow-skip set — proceeding."
    else
      fail "Target $TARGET_TAG skips past expected next version ($EXPECTED_NEXT_MINOR or $EXPECTED_NEXT_PATCH) from $LATEST_STABLE_TAG. Use --allow-skip if this skip is intentional."
    fi
  fi
fi

# ---------------------------------------------------------------------------
# Check: release notes exist
# ---------------------------------------------------------------------------

RELEASE_NOTES="$REPO_ROOT/docs/operations/release-notes-${TARGET_TAG}.md"
echo -e "  Release notes: ${BOLD}${RELEASE_NOTES#$REPO_ROOT/}${RESET}"

if [[ ! -f "$RELEASE_NOTES" ]]; then
  fail "Release notes not found: $RELEASE_NOTES"
  info "Create docs/operations/release-notes-${TARGET_TAG}.md before tagging."
else
  # Check that notes mention the correct version in the title
  if grep -q "# Daryaft ${TARGET_TAG}" "$RELEASE_NOTES" || grep -q "# Daryaft ${TARGET}" "$RELEASE_NOTES"; then
    pass "Release notes found and mention Daryaft $TARGET_TAG."
  else
    fail "Release notes found but do not contain '# Daryaft ${TARGET_TAG}' in the title."
  fi
fi

# ---------------------------------------------------------------------------
# Check: CHANGELOG.md has a target section
# ---------------------------------------------------------------------------

CHANGELOG="$REPO_ROOT/CHANGELOG.md"
echo -e "  Changelog: ${BOLD}CHANGELOG.md${RESET}"

if [[ ! -f "$CHANGELOG" ]]; then
  fail "CHANGELOG.md not found."
else
  # Accept either [X.Y.Z] or [vX.Y.Z] as a section header
  if grep -qE "^## \[(v?${TARGET})\]" "$CHANGELOG"; then
    pass "CHANGELOG.md has a [$TARGET] section."
  else
    fail "CHANGELOG.md does not have a [$TARGET] or [v$TARGET] section. Prepare the changelog entry before tagging."
    info "Tip: promote the [Unreleased] section to [$TARGET] - YYYY-MM-DD."
  fi
fi

# ---------------------------------------------------------------------------
# Check: local tag absence
# ---------------------------------------------------------------------------

if git rev-parse "$TARGET_TAG" > /dev/null 2>&1; then
  fail "Local tag $TARGET_TAG already exists. Delete it before re-running: git tag -d $TARGET_TAG"
else
  pass "Local tag $TARGET_TAG is absent."
fi

# ---------------------------------------------------------------------------
# Check: remote tag absence
# ---------------------------------------------------------------------------

REMOTE_TAG_SHA="$(git ls-remote --tags origin "refs/tags/${TARGET_TAG}" 2>/dev/null | awk '{print $1}')"
if [[ -n "$REMOTE_TAG_SHA" ]]; then
  fail "Remote tag $TARGET_TAG already exists on origin."
else
  pass "Remote tag $TARGET_TAG is absent."
fi

# ---------------------------------------------------------------------------
# Check: GitHub release absence (requires gh)
# ---------------------------------------------------------------------------

if command -v gh > /dev/null 2>&1; then
  GH_RELEASE_OUT="$(gh release view "$TARGET_TAG" --json tagName 2>&1 || true)"
  if echo "$GH_RELEASE_OUT" | grep -q "release not found"; then
    pass "GitHub release $TARGET_TAG does not exist."
  elif echo "$GH_RELEASE_OUT" | grep -q '"tagName"'; then
    fail "GitHub release $TARGET_TAG already exists."
  else
    warn "Could not determine GitHub release status: $GH_RELEASE_OUT"
  fi
else
  warn "gh CLI not found. Skipping GitHub release check."
fi

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------

echo ""
echo -e "  ${BOLD}Checks complete${RESET}"
echo -e "  Failures : ${FAILURES}"
echo -e "  Warnings : ${WARNINGS}"
echo ""

if [[ "$FAILURES" -eq 0 ]]; then
  echo -e "${GREEN}${BOLD}  Decision: PASS${RESET}"
  echo ""
  echo -e "  Proceed to quality gates, then:"
  echo -e "    git tag -a $TARGET_TAG -m \"Daryaft $TARGET_TAG\""
  echo -e "    git push origin $TARGET_TAG"
  echo ""
  exit 0
else
  echo -e "${RED}${BOLD}  Decision: FAIL${RESET} — $FAILURES check(s) failed."
  echo ""
  echo -e "  Resolve the failures above before creating tag $TARGET_TAG."
  echo ""
  exit 1
fi
