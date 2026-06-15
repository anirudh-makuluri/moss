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

Native runtime operations require the `libmoss` build tag:

```bash
go run -tags libmoss ./basic
go run -tags libmoss ./custom-embeddings
```
