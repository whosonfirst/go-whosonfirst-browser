package reader

import (
	"bytes"
	"context"
	"github.com/whosonfirst/go-ioutil"
	"io"
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

func (r *NullReader) Read(ctx context.Context, uri string) (io.ReadSeekCloser, error) {
	br := bytes.NewReader([]byte(uri))
	return ioutil.NewReadSeekCloser(br)
}

func (r *NullReader) ReaderURI(ctx context.Context, uri string) string {
	return uri
}
