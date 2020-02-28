package writer

import (
	"context"
	"io"
	"os"
)

func init() {

	ctx := context.Background()
	err := RegisterWriter(ctx, "stdout", initializeStdoutWriter)

	if err != nil {
		panic(err)
	}
}

func initializeStdoutWriter(ctx context.Context, uri string) (Writer, error) {

	wr := NewStdoutWriter()
	err := wr.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return wr, nil
}

type StdoutWriter struct {
	Writer
}

func NewStdoutWriter() Writer {
	wr := StdoutWriter{}
	return &wr
}

func (wr *StdoutWriter) Open(ctx context.Context, uri string) error {
	return nil
}

func (wr *StdoutWriter) Write(ctx context.Context, uri string, fh io.ReadCloser) error {
	_, err := io.Copy(os.Stdout, fh)
	return err
}

func (wr *StdoutWriter) URI(uri string) string {
	return uri
}
