package writer

import (
	"errors"
	wof_writer "github.com/whosonfirst/go-whosonfirst-readwrite/writer"
	"io"
)

type HTTPWriter struct {
	wof_writer.Writer
}

func NewHTTPWriter(root string) (wof_writer.Writer, error) {

	wr := HTTPWriter{}

	return &wr, nil
}

func (wr *HTTPWriter) Write(path string, fh io.ReadCloser) error {
	return errors.New("Please write me")
}

func (wr *HTTPWriter) URI(path string) string{
     return ""
}
