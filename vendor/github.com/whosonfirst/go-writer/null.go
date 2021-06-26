package writer

import (
	"context"
	"io"
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

func (wr *NullWriter) Write(ctx context.Context, uri string, fh io.ReadSeeker) (int64, error) {
	return io.Copy(io.Discard, fh)
}

func (wr *NullWriter) WriterURI(ctx context.Context, uri string) string {
	return uri
}

func (wr *NullWriter) Close(ctx context.Context) error {
	return nil
}
