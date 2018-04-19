package bytes

import (
	gobytes "bytes"
	"io"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func ReadCloserFromBytes(b []byte) (io.ReadCloser, error) {
	body := gobytes.NewReader(b)
	return nopCloser{body}, nil
}
