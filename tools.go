// +build tools
// The build tag above ensures that the imports in this file will not be
// included in the binaries using this package.
//
// To add dependencies, add the tool to the go mod
// `go get -u golang.org/x/lint/golint`
// Then, add it as an import to this file.
//
// To install the tools specified in this file, run the following:
// `go generate tools.go`
// or directly from tools.go:
// `cat tools.go | grep '^\s_ ' | sed 's/^.*_ "\(.*\)"/\1/g' | xargs -tI % go install %`

//go:generate go install github.com/goreleaser/goreleaser
//go:generate go install golang.org/x/tools/cmd/goimports
package goversion

import (
	_ "github.com/goreleaser/goreleaser"
	_ "golang.org/x/tools/cmd/goimports"
)