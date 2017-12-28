package reader

import (
	"bytes"
	"io"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type NullReader struct {
	Reader
}

func NewNullReader() (Reader, error) {

	r := NullReader{}
	return &r, nil
}

func (r *NullReader) Read(uri string) (io.ReadCloser, error) {

	buf := bytes.NewReader([]byte(uri))
	return nopCloser{buf}, nil
}
