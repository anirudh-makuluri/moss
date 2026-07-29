#!/usr/bin/env bash
set -euo pipefail

# Publishes one platform lib module tag at a time to avoid huge git pushes.
#
# Usage:
#   ./sdks/go/scripts/publish-platform-lib-tag.sh linux-amd64 v0.1.0

if [[ $# -lt 2 ]]; then
  echo "Usage: $0 <platform> <version> [remote]" >&2
  echo "Platforms: linux-amd64 linux-arm64 darwin-arm64 windows-amd64" >&2
  exit 1
fi

PLATFORM="$1"
VERSION="$2"
REMOTE="${3:-origin}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_ROOT="$(cd "${ROOT_DIR}/../.." && pwd)"
TAG="sdks/go/bindings/lib/${PLATFORM}/${VERSION}"
LIB_DIR="${REPO_ROOT}/sdks/go/bindings/lib/${PLATFORM}"

cd "${REPO_ROOT}"

if [[ ! -d "${LIB_DIR}" ]]; then
  echo "Missing platform directory: ${LIB_DIR}" >&2
  exit 1
fi

if [[ "${PLATFORM}" == "windows-amd64" ]]; then
  if [[ ! -f "${LIB_DIR}/moss.lib" ]]; then
    echo "Missing ${LIB_DIR}/moss.lib — run fetch-static-libs.sh first" >&2
    exit 1
  fi
else
  if [[ ! -f "${LIB_DIR}/libmoss.a" ]]; then
    echo "Missing ${LIB_DIR}/libmoss.a — run fetch-static-libs.sh first" >&2
    exit 1
  fi
fi

git add "${LIB_DIR}/go.mod" "${LIB_DIR}/prebuilt.go" "${LIB_DIR}/.gitignore"
git add -f "${LIB_DIR}/libmoss.a" "${LIB_DIR}/moss.lib" 2>/dev/null || true

if git diff --cached --quiet; then
  echo "Nothing to publish for ${PLATFORM}" >&2
  exit 1
fi

git commit -m "chore(go): publish ${PLATFORM} ${VERSION}"
git tag -f "${TAG}"
git push "${REMOTE}" "${TAG}"

echo "Published ${TAG}"
