package custom

import (
	"io/fs"
)

type CustomValidationFunc func([]byte) error

type CustomValidationWasm struct {
	FS   fs.FS
	Path string
	// You might be wondering how the JavaScript will know what function exported
	// by this WASM binary will invoke. The answer is that it doesn't. Or rather
	// in whosonfirst.browser.create.js there is an explicit check for a variable
	// named "whosonfirst_validate_feature_custom" of type "function" which means
	// it is expected that the code used to compile your GOOS=JS WASM binary exports
	// its custom validation functionality at that address (in addition to any
	// others you choose to define). For example:
	// https://github.com/sfomuseum/go-sfomuseum-validate-wasm/blob/main/cmd/validate_feature/main.go
}
