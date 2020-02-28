package reader

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
)

func init() {

	ctx := context.Background()
	err := RegisterReader(ctx, "null", initializeNullReader)

	if err != nil {
		panic(err)
	}
}

func initializeNullReader(ctx context.Context, uri string) (Reader, error) {

	r := NewNullReader()
	err := r.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return r, nil
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
