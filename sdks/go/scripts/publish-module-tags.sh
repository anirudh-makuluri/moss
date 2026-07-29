#!/usr/bin/env bash
set -euo pipefail

# Creates a release commit with native libraries and pushes module tags.
# The commit is tag-only and is not merged back into main.
#
# Usage:
#   ./sdks/go/scripts/publish-module-tags.sh v0.9.0 [remote]

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 <go-module-version> [git-remote]" >&2
  exit 1
fi

VERSION="$1"
REMOTE="${2:-origin}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_ROOT="$(cd "${ROOT_DIR}/../.." && pwd)"

TAGS=(
  "sdks/go/sdk/${VERSION}"
  "sdks/go/bindings/${VERSION}"
  "sdks/go/bindings/lib/linux-amd64/${VERSION}"
  "sdks/go/bindings/lib/linux-arm64/${VERSION}"
  "sdks/go/bindings/lib/darwin-arm64/${VERSION}"
  "sdks/go/bindings/lib/windows-amd64/${VERSION}"
)

cd "${REPO_ROOT}"

if [[ -z "$(git status --porcelain sdks/go)" ]]; then
  echo "No changes under sdks/go to publish" >&2
  exit 1
fi

git add sdks/go/bindings/include/libmoss.h
git add sdks/go/bindings/lib/*/libmoss.a sdks/go/bindings/lib/*/moss.lib 2>/dev/null || true
git add sdks/go/bindings/go.mod sdks/go/bindings/go.sum
git add sdks/go/sdk/go.mod sdks/go/sdk/go.sum
git add sdks/go/bindings/version.go

if git diff --cached --quiet; then
  echo "Nothing staged for release commit" >&2
  exit 1
fi

git commit -m "chore(go): publish ${VERSION} native bindings"

for tag in "${TAGS[@]}"; do
  git tag -f "${tag}"
done

git push "${REMOTE}" "${TAGS[@]}"

echo "Published tags:"
printf '  %s\n' "${TAGS[@]}"
