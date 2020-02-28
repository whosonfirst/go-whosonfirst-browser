package writer

import (
	"context"
	"io"
	"io/ioutil"
)

func init() {

	ctx := context.Background()
	err := RegisterWriter(ctx, "null", initializeNullWriter)

	if err != nil {
		panic(err)
	}
}

func initializeNullWriter(ctx context.Context, uri string) (Writer, error) {

	wr := NewNullWriter()
	err := wr.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return wr, nil
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
