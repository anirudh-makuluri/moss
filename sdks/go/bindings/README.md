# Moss Go Bindings

This package wraps the native `libmoss` runtime for Go via CGO.

It mirrors the role of the other language bindings packages in this repository:

- native runtime access
- local index loading
- local query execution
- cloud-backed manage operations exposed through the native client

## Installation (consumers)

```bash
go get github.com/usemoss/moss/sdks/go/sdk
go run github.com/usemoss/moss/sdks/go/tools/install@latest
```

The install tool downloads a prebuilt static library for your platform from
[C SDK GitHub Releases](https://github.com/usemoss/moss/releases). No manual C
SDK download, `LD_LIBRARY_PATH`, or `-tags libmoss` is required.

Requirements:

- `CGO_ENABLED=1` (default on Linux and macOS)
- A C compiler (`gcc` on Linux, Xcode CLI tools on macOS, MinGW on Windows)

## Layout

```
bindings/
  include/libmoss.h          # committed C header
  libmoss.go                 # CGO wrapper (requires CGO)
  prebuilt_<os>_<arch>.go    # per-platform CGO linker flags
  generate.go                # //go:generate install hook
  lib/
    linux-amd64/             # libmoss.a (gitignored, downloaded at build time)
    linux-arm64/
    darwin-arm64/
    windows-amd64/
```

Native `.a` / `.lib` files are gitignored. The
[`tools/install`](../tools/install) command fetches them from GitHub Releases and
verifies SHA256 checksums.

## Local development

Install the native library for your current machine:

```bash
./sdks/go/scripts/link_dev_lib.sh c-sdk-v0.9.0
```

Or fetch all supported platforms:

```bash
./sdks/go/scripts/fetch-static-libs.sh c-sdk-v0.9.0
```

Or use `go generate` from this directory:

```bash
cd sdks/go/bindings
go generate .
```

Then build with CGO enabled:

```bash
CGO_ENABLED=1 go build .
```

## Publishing

Maintainers run the **Publish Go SDK** GitHub Actions workflow. It creates
source-only module tags (no binaries in git):

- `sdks/go/sdk/v0.1.0`
- `sdks/go/bindings/v0.1.0`

Consumers download native libraries at build time via `tools/install`.

## Build without CGO

When `CGO_ENABLED=0`, this package builds a stub that returns
`ErrBindingsUnavailable`. The public SDK can still run unit tests and cloud query
fallback tests without native libraries.
