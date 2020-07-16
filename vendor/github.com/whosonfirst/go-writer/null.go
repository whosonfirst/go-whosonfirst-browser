package writer

import (
	"context"
	"io"
	"io/ioutil"
)

type NullWriter struct {
	Writer
}

func init() {

	ctx := context.Background()
	err := RegisterWriter(ctx, "null", NewNullWriter)

	if err != nil {
		panic(err)
	}
}

func NewNullWriter(ctx context.Context, uri string) (Writer, error) {

	wr := &NullWriter{}
	return wr, nil
}

func (wr *NullWriter) Write(ctx context.Context, uri string, fh io.ReadCloser) error {
	_, err := io.Copy(ioutil.Discard, fh)
	return err
}

func (wr *NullWriter) URI(uri string) string {
	return uri
}
