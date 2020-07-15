package reader

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
)

type NullReader struct {
	Reader
}

func init() {

	ctx := context.Background()
	err := RegisterReader(ctx, "null", NewNullReader)

	if err != nil {
		panic(err)
	}
}

func NewNullReader(ctx context.Context, uri string) (Reader, error) {

	r := &NullReader{}
	return r, nil
}

func (r *NullReader) Read(ctx context.Context, uri string) (io.ReadCloser, error) {
	br := bytes.NewReader([]byte(uri))
	return ioutil.NopCloser(br), nil
}

func (r *NullReader) URI(uri string) string {
	return uri
}
