# Moss Go Bindings

This package wraps the native `libmoss` runtime for Go via CGO.

It mirrors the role of the other language bindings packages in this repository:

- native runtime access
- local index loading
- local query execution
- cloud-backed manage operations exposed through the native client

## Status

The real bindings implementation is compiled only with the `libmoss` build tag.
Without that tag, this package builds a stub that returns a clear
`ErrBindingsUnavailable` error.

## Bundled native runtime

The bindings module vendors prebuilt `libmoss` artifacts from:

- <https://github.com/usemoss/moss/releases/tag/c-sdk-v0.9.0>

Supported bundled targets:

- `linux/amd64`
- `linux/arm64`
- `darwin/arm64`
- `windows/amd64`

Build with the `libmoss` tag:

```bash
go test -tags libmoss ./...
```

No external `CGO_CFLAGS`, `CGO_LDFLAGS`, or local C SDK install is needed for
the supported targets. Without `-tags libmoss`, the package still builds the
stub implementation for cloud fallback and unsupported native environments.

The bundled files live under [`internal/native`](./internal/native). Linux and
macOS builds link against the bundled shared library with an rpath pointing back
to the module cache. Windows builds link against the bundled import library and
ship the matching DLL in the same native artifact directory.
