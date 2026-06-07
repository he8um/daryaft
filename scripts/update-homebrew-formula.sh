#!/usr/bin/env bash
# update-homebrew-formula.sh — Update the Daryaft Homebrew formula in a local tap checkout.
#
# Usage:
#   scripts/update-homebrew-formula.sh --version VERSION --tap-dir PATH [--dry-run] [--repo OWNER/REPO]
#
# This script reads checksums.txt from a published GitHub release and updates
# Formula/daryaft.rb in a local tap checkout. It never pushes, commits, or
# creates releases.

set -euo pipefail

# ─── Constants ────────────────────────────────────────────────────────────────

SCRIPT_NAME="$(basename "$0")"
DEFAULT_REPO="he8um/daryaft"
FORMULA_REL="Formula/daryaft.rb"

ARM64_ARCHIVE="daryaft_darwin_arm64.tar.gz"
AMD64_ARCHIVE="daryaft_darwin_amd64.tar.gz"

# ─── Helpers ──────────────────────────────────────────────────────────────────

die() {
  echo "[${SCRIPT_NAME}] ERROR: $*" >&2
  exit 1
}

info() {
  echo "[${SCRIPT_NAME}] $*"
}

usage() {
  cat <<EOF
Usage: ${SCRIPT_NAME} --version VERSION --tap-dir PATH [OPTIONS]

Update Formula/daryaft.rb in a local tap checkout to a new release version.

Required:
  --version VERSION     Release version to update to, e.g. 1.2.0 or v1.2.0
  --tap-dir PATH        Path to a local clone of the Homebrew tap repository

Options:
  --dry-run             Show intended changes without modifying the formula
  --repo OWNER/REPO     GitHub repository (default: ${DEFAULT_REPO})
  --help                Show this help message

Examples:
  ${SCRIPT_NAME} --version 1.2.0 --tap-dir /tmp/homebrew-tap
  ${SCRIPT_NAME} --version 1.2.0 --tap-dir /tmp/homebrew-tap --dry-run
  ${SCRIPT_NAME} --version v1.2.0 --tap-dir /tmp/homebrew-tap

This script never pushes to the tap, never commits, and never creates releases.
After running, review the diff, validate with brew, and push manually.

Post-update validation commands:
  ruby -c Formula/daryaft.rb
  brew reinstall he8um/tap/daryaft
  daryaft version
  daryaft update --check
  daryaft doctor
EOF
}

# ─── Argument parsing ─────────────────────────────────────────────────────────

VERSION=""
TAP_DIR=""
DRY_RUN=false
REPO="${DEFAULT_REPO}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      VERSION="${2:-}"
      shift 2
      ;;
    --tap-dir)
      TAP_DIR="${2:-}"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    --repo)
      REPO="${2:-}"
      shift 2
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      die "Unknown argument: $1. Run with --help for usage."
      ;;
  esac
done

# ─── Input validation ─────────────────────────────────────────────────────────

[[ -n "${VERSION}" ]] || die "--version is required. Example: --version 1.2.0"
[[ -n "${TAP_DIR}" ]] || die "--tap-dir is required. Example: --tap-dir /tmp/homebrew-tap"

# Normalize: strip leading 'v' if present
VERSION="${VERSION#v}"

# Reject empty after stripping
[[ -n "${VERSION}" ]] || die "Version is empty after stripping 'v' prefix."

# Reject whitespace
[[ "${VERSION}" == "${VERSION//[[:space:]]/}" ]] || die "Version must not contain spaces: '${VERSION}'"

# Reject dev/pre-release suffixes
if echo "${VERSION}" | grep -qE '[-+]'; then
  die "Version '${VERSION}' looks like a pre-release or dev version. Only stable release versions are accepted."
fi

# Require semver-like form X.Y.Z
if ! echo "${VERSION}" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'; then
  die "Version '${VERSION}' is not in X.Y.Z format."
fi

TAG="v${VERSION}"

# ─── Tap directory validation ─────────────────────────────────────────────────

[[ -d "${TAP_DIR}" ]] || die "Tap directory does not exist: ${TAP_DIR}"
[[ -d "${TAP_DIR}/.git" ]] || die "Tap directory is not a git repository: ${TAP_DIR}"

FORMULA_PATH="${TAP_DIR}/${FORMULA_REL}"
[[ -f "${FORMULA_PATH}" ]] || die "Formula file not found: ${FORMULA_PATH}"

# Check for unexpected local modifications before editing
TAP_STATUS="$(git -C "${TAP_DIR}" status --porcelain 2>/dev/null)"
if [[ -n "${TAP_STATUS}" ]]; then
  die "Tap working tree has uncommitted changes. Commit or reset before running this script.
Dirty files:
${TAP_STATUS}"
fi

# ─── Fetch checksums from GitHub release ─────────────────────────────────────

CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${TAG}/checksums.txt"
info "Fetching checksums from: ${CHECKSUMS_URL}"

if ! command -v curl >/dev/null 2>&1; then
  die "curl is required but not found."
fi

CHECKSUMS_TMPFILE="$(mktemp)"
trap 'rm -f "${CHECKSUMS_TMPFILE}"' EXIT

HTTP_STATUS="$(curl -fsSL --write-out "%{http_code}" --output "${CHECKSUMS_TMPFILE}" "${CHECKSUMS_URL}" 2>&1)" || {
  die "Failed to download checksums.txt. Verify that the GitHub release ${TAG} exists and has assets.
URL: ${CHECKSUMS_URL}"
}

CHECKSUMS="$(cat "${CHECKSUMS_TMPFILE}")"
[[ -n "${CHECKSUMS}" ]] || die "Downloaded checksums.txt is empty."

info "Checksums fetched successfully."

# ─── Extract required SHA-256 values ─────────────────────────────────────────

SHA256_ARM64="$(echo "${CHECKSUMS}" | awk -v f="${ARM64_ARCHIVE}" '$2==f{print $1}')"
SHA256_AMD64="$(echo "${CHECKSUMS}" | awk -v f="${AMD64_ARCHIVE}" '$2==f{print $1}')"

[[ -n "${SHA256_ARM64}" ]] || die "SHA-256 for ${ARM64_ARCHIVE} not found in checksums.txt."
[[ -n "${SHA256_AMD64}" ]] || die "SHA-256 for ${AMD64_ARCHIVE} not found in checksums.txt."

# Basic sanity: SHA-256 is 64 hex chars
if ! echo "${SHA256_ARM64}" | grep -qE '^[0-9a-f]{64}$'; then
  die "SHA-256 for ${ARM64_ARCHIVE} does not look valid: '${SHA256_ARM64}'"
fi
if ! echo "${SHA256_AMD64}" | grep -qE '^[0-9a-f]{64}$'; then
  die "SHA-256 for ${AMD64_ARCHIVE} does not look valid: '${SHA256_AMD64}'"
fi

info "arm64 SHA-256: ${SHA256_ARM64}"
info "amd64 SHA-256: ${SHA256_AMD64}"

# ─── Detect current formula version ──────────────────────────────────────────

CURRENT_VERSION="$(grep -E '^\s+version "' "${FORMULA_PATH}" | head -1 | sed 's/.*version "\(.*\)".*/\1/')"
info "Current formula version: ${CURRENT_VERSION:-unknown}"
info "Target version:          ${VERSION}"

if [[ "${CURRENT_VERSION}" == "${VERSION}" ]]; then
  info "Formula is already at version ${VERSION}. Checking if checksums also match..."
  CURRENT_ARM64="$(grep -A2 'Hardware::CPU.arm?' "${FORMULA_PATH}" | grep sha256 | head -1 | sed 's/.*sha256 "\(.*\)".*/\1/')"
  CURRENT_AMD64="$(grep -A2 'else' "${FORMULA_PATH}" | grep sha256 | head -1 | sed 's/.*sha256 "\(.*\)".*/\1/')"
  if [[ "${CURRENT_ARM64}" == "${SHA256_ARM64}" && "${CURRENT_AMD64}" == "${SHA256_AMD64}" ]]; then
    info "Formula is already up to date for version ${VERSION}. No changes needed."
    exit 0
  else
    info "Version matches but checksums differ — updating checksums."
  fi
fi

# ─── Verify formula structure is recognizable ─────────────────────────────────

grep -q 'version "' "${FORMULA_PATH}" || die "Cannot find 'version' field in formula. Layout may have changed."
grep -q 'Hardware::CPU.arm?' "${FORMULA_PATH}" || die "Cannot find arm? CPU guard in formula. Layout may have changed."
grep -q 'daryaft_darwin_arm64.tar.gz' "${FORMULA_PATH}" || die "Cannot find arm64 archive reference in formula."
grep -q 'daryaft_darwin_amd64.tar.gz' "${FORMULA_PATH}" || die "Cannot find amd64 archive reference in formula."

# ─── Build new URLs ───────────────────────────────────────────────────────────

BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"
URL_ARM64="${BASE_URL}/${ARM64_ARCHIVE}"
URL_AMD64="${BASE_URL}/${AMD64_ARCHIVE}"

# ─── Apply or preview changes ─────────────────────────────────────────────────

# Build the new formula content via sed substitutions.
# We update:
#   version "X.Y.Z"
#   the arm64 url line
#   the sha256 line after the arm64 url
#   the amd64 url line (in the else branch)
#   the sha256 line after the amd64 url

UPDATED_FORMULA="$(
  perl -0777 -pe "
    # Replace version field
    s{(version \")[^\"]+(\")}{\${1}${VERSION}\${2}};

    # Replace arm64 url
    s{(url \"https://github\.com/[^/]+/daryaft/releases/download/)[^/]+/(daryaft_darwin_arm64\.tar\.gz\")}{\${1}${TAG}/\${2}};

    # Replace arm64 sha256 (the sha256 that follows the arm64 url line)
    s{(daryaft_darwin_arm64\.tar\.gz\"[^\n]*\n\s+sha256 \")[^\"]+(\")}{
      my \$pre = \$1;
      \$pre . \"${SHA256_ARM64}\" . \$2
    }e;

    # Replace amd64 url
    s{(url \"https://github\.com/[^/]+/daryaft/releases/download/)[^/]+/(daryaft_darwin_amd64\.tar\.gz\")}{\${1}${TAG}/\${2}};

    # Replace amd64 sha256 (the sha256 that follows the amd64 url line)
    s{(daryaft_darwin_amd64\.tar\.gz\"[^\n]*\n\s+sha256 \")[^\"]+(\")}{
      my \$pre = \$1;
      \$pre . \"${SHA256_AMD64}\" . \$2
    }e;
  " "${FORMULA_PATH}"
)"

if "${DRY_RUN}"; then
  info "--- DRY RUN --- No files will be modified."
  echo ""
  echo "=== Intended changes to ${FORMULA_REL} ==="
  TMPFILE_ORIG="$(mktemp)"
  TMPFILE_NEW="$(mktemp)"
  trap 'rm -f "${CHECKSUMS_TMPFILE}" "${TMPFILE_ORIG}" "${TMPFILE_NEW}"' EXIT
  cp "${FORMULA_PATH}" "${TMPFILE_ORIG}"
  echo "${UPDATED_FORMULA}" > "${TMPFILE_NEW}"
  diff "${TMPFILE_ORIG}" "${TMPFILE_NEW}" || true
  echo ""
  info "Dry run complete. Run without --dry-run to apply changes."
  exit 0
fi

# Write updated formula
echo "${UPDATED_FORMULA}" > "${FORMULA_PATH}"

info "Formula updated: ${FORMULA_PATH}"

# ─── Syntax validation ────────────────────────────────────────────────────────

if command -v ruby >/dev/null 2>&1; then
  info "Running ruby -c syntax check..."
  ruby -c "${FORMULA_PATH}" && info "Ruby syntax OK."
else
  info "ruby not found — skipping syntax check."
fi

# ─── Summary ─────────────────────────────────────────────────────────────────

echo ""
echo "=== Formula update complete ==="
echo "  Version : ${VERSION}"
echo "  arm64   : ${SHA256_ARM64}"
echo "  amd64   : ${SHA256_AMD64}"
echo ""
echo "=== Diff ==="
git -C "${TAP_DIR}" diff "${FORMULA_REL}" || true
echo ""
echo "=== Next steps (manual) ==="
echo "  1. Review diff above."
echo "  2. Run: ruby -c ${FORMULA_PATH}"
echo "  3. Run: brew reinstall he8um/tap/daryaft"
echo "  4. Run: daryaft version && daryaft update --check && daryaft doctor"
echo "  5. Run in tap dir:"
echo "       git add Formula/daryaft.rb"
echo "       git commit -m \"Update daryaft to v${VERSION}\""
echo "       git push"
echo ""
echo "  This script does not push or commit."
