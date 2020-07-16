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

func (wr *StdoutWriter) Write(ctx context.Context, uri string, fh io.ReadCloser) error {
	_, err := io.Copy(os.Stdout, fh)
	return err
}

func (wr *StdoutWriter) URI(uri string) string {
	return uri
}
