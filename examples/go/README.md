# Moss Go Examples

Runnable examples for the Moss Go SDK.

## Examples

- [`basic/main.go`](./basic/main.go) creates an index, loads it, queries it, and deletes it.
- [`custom-embeddings/main.go`](./custom-embeddings/main.go) uses caller-provided vectors for documents and queries.

## Run

Set your Moss credentials:

```bash
export MOSS_PROJECT_ID=...
export MOSS_PROJECT_KEY=...
```

Install the native library for your platform (monorepo dev):

```bash
../../sdks/go/scripts/link_dev_lib.sh c-sdk-v0.9.0
```

Then run an example:

```bash
cd examples/go
go run ./basic
go run ./custom-embeddings
```

Published installs (`go get github.com/usemoss/moss/sdks/go/sdk`) bundle static
libraries automatically and do not need the script above.
