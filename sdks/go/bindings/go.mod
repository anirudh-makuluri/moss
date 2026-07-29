module github.com/usemoss/moss/sdks/go/bindings

go 1.22.2

require (
	github.com/usemoss/moss/sdks/go/bindings/lib/darwin-arm64 v0.0.0
	github.com/usemoss/moss/sdks/go/bindings/lib/linux-amd64 v0.0.0
	github.com/usemoss/moss/sdks/go/bindings/lib/linux-arm64 v0.0.0
	github.com/usemoss/moss/sdks/go/bindings/lib/windows-amd64 v0.0.0
)

replace github.com/usemoss/moss/sdks/go/bindings/lib/darwin-arm64 => ./lib/darwin-arm64

replace github.com/usemoss/moss/sdks/go/bindings/lib/linux-amd64 => ./lib/linux-amd64

replace github.com/usemoss/moss/sdks/go/bindings/lib/linux-arm64 => ./lib/linux-arm64

replace github.com/usemoss/moss/sdks/go/bindings/lib/windows-amd64 => ./lib/windows-amd64
