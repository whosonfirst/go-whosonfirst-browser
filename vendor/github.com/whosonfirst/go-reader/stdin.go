package reader

import (
	"context"
	"io"
	"os"

	"github.com/whosonfirst/go-ioutil"
)

type StdinReader struct {
	Reader
}

func init() {

	ctx := context.Background()
	err := RegisterReader(ctx, "stdin", NewStdinReader)

	if err != nil {
		panic(err)
	}
}

func NewStdinReader(ctx context.Context, uri string) (Reader, error) {

	r := &StdinReader{}
	return r, nil
}

func (r *StdinReader) Read(ctx context.Context, uri string) (io.ReadSeekCloser, error) {
	return ioutil.NewReadSeekCloser(os.Stdin)
}

func (r *StdinReader) ReaderURI(ctx context.Context, uri string) string {
	return "-"
}
