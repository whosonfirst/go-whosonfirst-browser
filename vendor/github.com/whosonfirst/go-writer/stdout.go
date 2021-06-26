package writer

import (
	"context"
	"io"
	"os"
)

type StdoutWriter struct {
	Writer
}

func init() {

	ctx := context.Background()
	err := RegisterWriter(ctx, "stdout", NewStdoutWriter)

	if err != nil {
		panic(err)
	}
}

func NewStdoutWriter(ctx context.Context, uri string) (Writer, error) {

	wr := &StdoutWriter{}
	return wr, nil
}

func (wr *StdoutWriter) Write(ctx context.Context, uri string, fh io.ReadSeeker) (int64, error) {
	return io.Copy(os.Stdout, fh)
}

func (wr *StdoutWriter) WriterURI(ctx context.Context, uri string) string {
	return uri
}

func (wr *StdoutWriter) Close(ctx context.Context) error {
	return nil
}
