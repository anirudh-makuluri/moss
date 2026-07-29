#!/usr/bin/env bash
set -euo pipefail

# Publishes source-only bindings + sdk module tags.
#
# Usage:
#   ./sdks/go/scripts/publish-sdk-module-tags.sh v0.1.2

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
git add sdks/go/bindings/include/libmoss.h sdks/go/bindings/generate.go
git add sdks/go/tools/install

if git diff --cached --quiet; then
  echo "Nothing to publish for bindings/sdk ${VERSION}" >&2
  exit 1
fi

git commit -m "chore(go): publish bindings and sdk ${VERSION}"

for tag in \
  "sdks/go/bindings/${VERSION}" \
  "sdks/go/sdk/${VERSION}" \
  "sdks/go/tools/install/${VERSION}"; do
  git tag -f "${tag}"
  git push "${REMOTE}" "${tag}"
  echo "Published ${tag}"
done
