#!/usr/bin/env bash
set -euo pipefail

# Updates bindings and sdk go.mod files to pin platform lib module versions.
#
# Usage:
#   ./sdks/go/scripts/bump-module-versions.sh v0.9.0

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <go-module-version>" >&2
  echo "Example: $0 v0.9.0" >&2
  exit 1
fi

VERSION="$1"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BINDINGS_DIR="${ROOT_DIR}/bindings"
SDK_DIR="${ROOT_DIR}/sdk"

PLATFORMS=(
  linux-amd64
  linux-arm64
  darwin-arm64
  windows-amd64
)

(
  cd "${BINDINGS_DIR}"
  for platform in "${PLATFORMS[@]}"; do
    go mod edit -require="github.com/usemoss/moss/sdks/go/bindings/lib/${platform}@${VERSION}"
    go mod edit -dropreplace="github.com/usemoss/moss/sdks/go/bindings/lib/${platform}" 2>/dev/null || true
  done
)

(
  cd "${SDK_DIR}"
  go mod edit -require="github.com/usemoss/moss/sdks/go/bindings@${VERSION}"
  go mod edit -dropreplace=github.com/usemoss/moss/sdks/go/bindings 2>/dev/null || true
)

echo "Pinned Go module versions to ${VERSION}"
echo "Publish platform lib tags first, then run publish-sdk-module-tags.sh"
