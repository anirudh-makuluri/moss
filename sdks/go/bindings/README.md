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
```

Published releases bundle prebuilt static libraries per platform. No manual C SDK
download, `LD_LIBRARY_PATH`, or `-tags libmoss` is required.

Requirements:

- `CGO_ENABLED=1` (default on Linux and macOS)
- A C compiler (`gcc` on Linux, Xcode CLI tools on macOS, MinGW on Windows)

## Layout

```
bindings/
  include/libmoss.h          # committed C header
  libmoss.go                 # CGO wrapper (requires CGO)
  prebuilt_<os>_<arch>.go    # selects the platform linker module
  lib/
    linux-amd64/             # separate Go module with libmoss.a (release tags)
    linux-arm64/
    darwin-arm64/
    windows-amd64/
```

On `main`, native `.a` / `.lib` files are gitignored. Release CI downloads them
from [C SDK GitHub Releases](https://github.com/usemoss/moss/releases) and
publishes versioned Go module tags.

## Local development

Install the native library for your current machine:

```bash
./sdks/go/scripts/link_dev_lib.sh c-sdk-v0.9.0
```

Or fetch all supported platforms (release prep):

```bash
./sdks/go/scripts/fetch-static-libs.sh c-sdk-v0.9.0
```

Then build with CGO enabled:

```bash
cd sdks/go/bindings
CGO_ENABLED=1 go build .
```

## Publishing

Maintainers run the **Publish Go SDK** GitHub Actions workflow. It:

1. Downloads static libraries from a C SDK release
2. Pins module versions in `go.mod` files
3. Creates a release commit and pushes module tags (without updating `main`)

Tags follow the monorepo convention:

- `sdks/go/sdk/v0.9.0`
- `sdks/go/bindings/v0.9.0`
- `sdks/go/bindings/lib/linux-amd64/v0.9.0`
- …

## Build without CGO

When `CGO_ENABLED=0`, this package builds a stub that returns
`ErrBindingsUnavailable`. The public SDK can still run unit tests and cloud query
fallback tests without native libraries.
