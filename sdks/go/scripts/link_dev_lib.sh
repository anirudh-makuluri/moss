#!/usr/bin/env bash
set -euo pipefail

# Installs the native library for the current machine only (local dev helper).
#
# Usage:
#   ./sdks/go/scripts/link_dev_lib.sh [c-sdk-tag]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RELEASE_TAG="${1:-c-sdk-v0.9.0}"
LIBMOSS_VERSION="${RELEASE_TAG#c-sdk-}"
REPO="${MOSS_RELEASE_REPO:-usemoss/moss}"
BASE_URL="https://github.com/${REPO}/releases/download/${RELEASE_TAG}"

case "$(uname -s)/$(uname -m)" in
  Linux/x86_64)
    PLATFORM="linux-amd64"
    TRIPLE="x86_64-unknown-linux-gnu"
    LIB_NAME="libmoss.a"
    SRC_LIB="lib/libmoss.a"
    ;;
  Linux/aarch64|Linux/arm64)
    PLATFORM="linux-arm64"
    TRIPLE="aarch64-unknown-linux-gnu"
    LIB_NAME="libmoss.a"
    SRC_LIB="lib/libmoss.a"
    ;;
  Darwin/arm64)
    PLATFORM="darwin-arm64"
    TRIPLE="aarch64-apple-darwin"
    LIB_NAME="libmoss.a"
    SRC_LIB="lib/libmoss.a"
    ;;
  MINGW*|MSYS*|CYGWIN*)
    PLATFORM="windows-amd64"
    TRIPLE="x86_64-pc-windows-msvc"
    LIB_NAME="moss.lib"
    SRC_LIB="lib/moss.lib"
    ;;
  *)
    echo "Unsupported dev platform: $(uname -s)/$(uname -m)" >&2
    exit 1
    ;;
esac

ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BINDINGS_DIR="${ROOT_DIR}/bindings"
DEST="${BINDINGS_DIR}/lib/${PLATFORM}"
ARCHIVE="libmoss-${LIBMOSS_VERSION}-${TRIPLE}.tar.gz"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "${TMPDIR}"' EXIT

mkdir -p "${DEST}" "${BINDINGS_DIR}/include"

echo "Downloading ${ARCHIVE}"
curl -fsSL "${BASE_URL}/${ARCHIVE}" -o "${TMPDIR}/${ARCHIVE}"
tar -xzf "${TMPDIR}/${ARCHIVE}" -C "${TMPDIR}"

EXTRACTED="${TMPDIR}/libmoss-${LIBMOSS_VERSION}-${TRIPLE}"
cp "${EXTRACTED}/include/libmoss.h" "${BINDINGS_DIR}/include/libmoss.h"
cp "${EXTRACTED}/${SRC_LIB}" "${DEST}/${LIB_NAME}"

echo "Installed ${PLATFORM} into ${DEST}/${LIB_NAME}"
