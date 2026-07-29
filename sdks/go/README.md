# Moss Go SDK

The Go SDK follows the same two-layer design as the other Moss SDKs:

- `sdks/go/sdk/` contains the public Go SDK
- `sdks/go/bindings/` wraps the native `libmoss` runtime via CGO

## Install

```bash
go get github.com/usemoss/moss/sdks/go/sdk
```

Published module versions bundle static `libmoss` per platform. You need CGO and a
C compiler, but not a manual C SDK download.

## Local development

```bash
./sdks/go/scripts/link_dev_lib.sh c-sdk-v0.9.0
cd sdks/go/sdk
CGO_ENABLED=1 go test ./...
```

Unit tests run without native libraries when `CGO_ENABLED=0`.

## Publishing

See [`bindings/README.md`](./bindings/README.md) and
[`.github/workflows/publish-go-sdk.yml`](../../.github/workflows/publish-go-sdk.yml).

The public SDK module lives under [`sdk/`](./sdk/), and the native bindings module
lives under [`bindings/`](./bindings/).
