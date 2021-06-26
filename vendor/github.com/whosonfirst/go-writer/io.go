package writer

import (
	"context"
	"errors"
	"io"
)

const IOWRITER_TARGET_KEY string = "github.com/whosonfirst/go-writer#io_writer"

type IOWriter struct {
	Writer
}

func init() {

	ctx := context.Background()
	err := RegisterWriter(ctx, "io", NewIOWriter)

	if err != nil {
		panic(err)
	}
}

func NewIOWriter(ctx context.Context, uri string) (Writer, error) {

	wr := &IOWriter{}
	return wr, nil
}

func (wr *IOWriter) Write(ctx context.Context, uri string, fh io.ReadSeeker) (int64, error) {

	target, err := GetIOWriterFromContext(ctx)

	if err != nil {
		return 0, err
	}

	return io.Copy(target, fh)
}

func (wr *IOWriter) WriterURI(ctx context.Context, uri string) string {
	return uri
}

func (wr *IOWriter) Close(ctx context.Context) error {
	return nil
}

func SetIOWriterWithContext(ctx context.Context, wr io.Writer) (context.Context, error) {

	ctx = context.WithValue(ctx, IOWRITER_TARGET_KEY, wr)
	return ctx, nil
}

func GetIOWriterFromContext(ctx context.Context) (io.Writer, error) {

	v := ctx.Value(IOWRITER_TARGET_KEY)

	if v == nil {
		return nil, errors.New("Missing writer")
	}

	var target io.Writer

	switch v.(type) {
	case io.Writer:
		target = v.(io.Writer)
	default:
		return nil, errors.New("Invalid writer")
	}

	return target, nil
}
