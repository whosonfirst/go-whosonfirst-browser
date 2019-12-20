package writer

import (
	"context"
	"io"
	"io/ioutil"
)

func init() {
	wr := NewNullWriter()
	Register("null", wr)
}

type NullWriter struct {
	Writer
}

func NewNullWriter() Writer {

	wr := NullWriter{}
	return &wr
}

func (wr *NullWriter) Open(ctx context.Context, uri string) error {
	return nil
}

func (wr *NullWriter) Write(ctx context.Context, uri string, fh io.ReadCloser) error {
	_, err := io.Copy(ioutil.Discard, fh)
	return err
}

func (wr *NullWriter) URI(uri string) string {
	return uri
}
