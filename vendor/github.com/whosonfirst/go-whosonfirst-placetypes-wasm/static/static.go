// package static provides an `embed.FS` containing JavaScript and WebAssembly binaries used by the go-whosonfirst-export-wasm tools and methods.
package static

import (
	"embed"
)

//go:embed javascript/* wasm/*
var FS embed.FS
