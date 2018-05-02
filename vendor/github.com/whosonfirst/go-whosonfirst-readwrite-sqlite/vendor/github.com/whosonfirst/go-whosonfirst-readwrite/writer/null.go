package writer

import (
	"io"
	_ "log"
)

type NullWriter struct {
	Writer
}

func NewNullWriter() (Writer, error) {

	w := NullWriter{}
	return &w, nil
}

func (w *NullWriter) Write(path string, fh io.ReadCloser) error {
	// maybe drain fh here?
	return nil
}

func (w *NullWriter) URI(path string) string {
     return "/dev/null"
}
