package custom

import (
	"io/fs"
)

type CustomValidationFunc func([]byte) error

type CustomValidationWasm struct {
	FS       fs.FS
	Path     string
	Function string
}
