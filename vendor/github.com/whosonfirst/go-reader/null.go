package reader

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
)

func init() {
	r := NewNullReader()
	Register("null", r)
}

type NullReader struct {
	Reader
}

func NewNullReader() Reader {

	r := NullReader{}
	return &r
}

func (r *NullReader) Open(ctx context.Context, uri string) error {
	return nil
}

func (r *NullReader) Read(ctx context.Context, uri string) (io.ReadCloser, error) {
	br := bytes.NewReader([]byte(uri))
	return ioutil.NopCloser(br), nil
}

func (r *NullReader) URI(uri string) string {
	return uri
}
