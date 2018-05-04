package reader

import (
	"github.com/whosonfirst/go-whosonfirst-readwrite/bytes"
	"io"
)

type NullReader struct {
	Reader
}

func NewNullReader() (Reader, error) {

	r := NullReader{}
	return &r, nil
}

func (r *NullReader) Read(uri string) (io.ReadCloser, error) {

	return bytes.ReadCloserFromBytes([]byte(uri))
}

func (r *NullReader) URI(uri string) string {
     return uri
}
