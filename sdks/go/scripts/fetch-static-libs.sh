#!/usr/bin/env bash
set -euo pipefail

# Downloads static libmoss archives from GitHub Releases and installs them into
# sdks/go/bindings/lib/<platform>/ for local development or release publishing.
#
# Usage:
#   ./sdks/go/scripts/fetch-static-libs.sh [c-sdk-tag]
#
# Example:
#   ./sdks/go/scripts/fetch-static-libs.sh c-sdk-v0.9.0

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BINDINGS_DIR="${ROOT_DIR}/bindings"
RELEASE_TAG="${1:-c-sdk-v0.9.0}"
LIBMOSS_VERSION="${RELEASE_TAG#c-sdk-}"
REPO="${MOSS_RELEASE_REPO:-usemoss/moss}"
BASE_URL="https://github.com/${REPO}/releases/download/${RELEASE_TAG}"

declare -A PLATFORM_TRIPLES=(
  ["linux-amd64"]="x86_64-unknown-linux-gnu"
  ["linux-arm64"]="aarch64-unknown-linux-gnu"
  ["darwin-arm64"]="aarch64-apple-darwin"
  ["windows-amd64"]="x86_64-pc-windows-msvc"
)

fetch_platform() {
  local platform="$1"
  local triple="$2"
  local archive="libmoss-${LIBMOSS_VERSION}-${triple}.tar.gz"
  local url="${BASE_URL}/${archive}"
  local tmpdir
  tmpdir="$(mktemp -d)"

  echo "Fetching ${platform} from ${url}"
  curl -fsSL "${url}" -o "${tmpdir}/${archive}"
  tar -xzf "${tmpdir}/${archive}" -C "${tmpdir}"

  local extracted="${tmpdir}/libmoss-${LIBMOSS_VERSION}-${triple}"
  local dest="${BINDINGS_DIR}/lib/${platform}"
  mkdir -p "${dest}"

  cp "${extracted}/include/libmoss.h" "${BINDINGS_DIR}/include/libmoss.h"

  case "${platform}" in
    windows-amd64)
      cp "${extracted}/lib/moss.lib" "${dest}/moss.lib"
      ;;
    *)
      cp "${extracted}/lib/libmoss.a" "${dest}/libmoss.a"
      ;;
  esac

  rm -rf "${tmpdir}"

  echo "Installed ${platform} native library into ${dest}"
}

main() {
  if ! command -v curl >/dev/null 2>&1; then
    echo "curl is required" >&2
    exit 1
  fi

  mkdir -p "${BINDINGS_DIR}/include"

  for platform in "${!PLATFORM_TRIPLES[@]}"; do
    fetch_platform "${platform}" "${PLATFORM_TRIPLES[$platform]}"
  done

  echo "Done. Header: ${BINDINGS_DIR}/include/libmoss.h"
}

main "$@"
