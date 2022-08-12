package writer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
)

// IOWRITER_TARGET_KEY is the key used to store an `io.Writer` instance in a `context.Context` instance.
const IOWRITER_TARGET_KEY string = "github.com/whosonfirst/go-writer#io_writer"

// IOWriter is a struct that implements the `Writer` interface for writing documents to an `io.Writer` instance.
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

// NewIOWriter returns a new `IOWriter` instance for writing documents to the current working directory
// configured by 'uri' in the form of:
//
//	io://
//
// In order to assign the actual `io.Writer` instance to use you will need to call the `SetIOWriterWithContext`
// method and pass the resultant `context.Context` instance to the `Write` method.
func NewIOWriter(ctx context.Context, uri string) (Writer, error) {
	wr := &IOWriter{}
	return wr, nil
}

// Write copies the content of 'fh' to 'path'. It is assumed that 'ctx' contains a valid `io.Writer` instance
// that has been assigned by the `SetIOWriterWithContext` method.
func (wr *IOWriter) Write(ctx context.Context, path string, fh io.ReadSeeker) (int64, error) {

	target, err := GetIOWriterFromContext(ctx)

	if err != nil {
		return 0, fmt.Errorf("Failed to get io.Writer instance from context, %w", err)
	}

	return io.Copy(target, fh)
}

// WriterURI returns the final URI for path.
func (wr *IOWriter) WriterURI(ctx context.Context, path string) string {
	return path
}

// Flush publish any outstanding data.
func (wr *IOWriter) Flush(ctx context.Context) error {
	return nil
}

// Close closes the underlying writer mechanism.
func (wr *IOWriter) Close(ctx context.Context) error {
	return nil
}

// SetLogger assigns 'logger' to 'wr'.
func (wr *IOWriter) SetLogger(ctx context.Context, logger *log.Logger) error {
	return nil
}

// SetIOWriterWithContext returns a new `context.Context` instance with 'wr' assigned to the
// `IOWRITER_TARGET_KEY` value.
func SetIOWriterWithContext(ctx context.Context, wr io.Writer) (context.Context, error) {

	ctx = context.WithValue(ctx, IOWRITER_TARGET_KEY, wr)
	return ctx, nil
}

// GetIOWriterFromContext returns the `io.Writer` instance associated with the `IOWRITER_TARGET_KEY`
// value in 'ctx'.
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
