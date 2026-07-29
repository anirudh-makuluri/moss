#!/usr/bin/env bash
set -euo pipefail

# Publishes bindings + sdk module tags (no native binaries in this commit).
#
# Usage:
#   ./sdks/go/scripts/publish-sdk-module-tags.sh v0.1.0

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 <version> [remote]" >&2
  exit 1
fi

VERSION="$1"
REMOTE="${2:-origin}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_ROOT="$(cd "${ROOT_DIR}/../.." && pwd)"

cd "${REPO_ROOT}"

"${ROOT_DIR}/scripts/bump-module-versions.sh" "${VERSION}"

git add sdks/go/bindings/go.mod sdks/go/sdk/go.mod sdks/go/bindings/version.go
git add sdks/go/bindings/include/libmoss.h

if git diff --cached --quiet; then
  echo "Nothing to publish for bindings/sdk ${VERSION}" >&2
  exit 1
fi

git commit -m "chore(go): publish bindings and sdk ${VERSION}"

for tag in "sdks/go/bindings/${VERSION}" "sdks/go/sdk/${VERSION}"; do
  git tag -f "${tag}"
  git push "${REMOTE}" "${tag}"
  echo "Published ${tag}"
done
