package reader

import (
	"io"
)

type Reader interface {
	Read(string) (io.ReadCloser, error)
	URI(string) string
}
