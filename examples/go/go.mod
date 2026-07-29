module github.com/usemoss/moss/examples/go

go 1.22.2

require github.com/usemoss/moss/sdks/go/sdk v0.0.0

require (
	github.com/usemoss/moss/sdks/go/bindings v0.0.0 // indirect
	github.com/usemoss/moss/sdks/go/bindings/lib/darwin-arm64 v0.0.0 // indirect
	github.com/usemoss/moss/sdks/go/bindings/lib/linux-amd64 v0.0.0 // indirect
	github.com/usemoss/moss/sdks/go/bindings/lib/linux-arm64 v0.0.0 // indirect
	github.com/usemoss/moss/sdks/go/bindings/lib/windows-amd64 v0.0.0 // indirect
)

replace github.com/usemoss/moss/sdks/go/sdk => ../../sdks/go/sdk

replace github.com/usemoss/moss/sdks/go/bindings => ../../sdks/go/bindings

replace github.com/usemoss/moss/sdks/go/bindings/lib/darwin-arm64 => ../../sdks/go/bindings/lib/darwin-arm64

replace github.com/usemoss/moss/sdks/go/bindings/lib/linux-amd64 => ../../sdks/go/bindings/lib/linux-amd64

replace github.com/usemoss/moss/sdks/go/bindings/lib/linux-arm64 => ../../sdks/go/bindings/lib/linux-arm64

replace github.com/usemoss/moss/sdks/go/bindings/lib/windows-amd64 => ../../sdks/go/bindings/lib/windows-amd64
