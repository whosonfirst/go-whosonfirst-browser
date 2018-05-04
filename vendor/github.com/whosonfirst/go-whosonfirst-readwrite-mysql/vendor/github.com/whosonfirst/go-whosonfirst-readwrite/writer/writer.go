package writer

import (
	"io"
)

type Writer interface {
	Write(string, io.ReadCloser) error
	URI(string) string
}
