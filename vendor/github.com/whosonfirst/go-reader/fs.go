package reader

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

func init() {

	ctx := context.Background()
	err := RegisterReader(ctx, "fs", initializeFSReader)

	if err != nil {
		panic(err)
	}
}

func initializeFSReader(ctx context.Context, uri string) (Reader, error) {

	r := NewFSReader()
	err := r.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return r, nil
}

type FSReader struct {
	Reader
	root string
}

func NewFSReader() Reader {
	r := FSReader{}
	return &r
}

func (r *FSReader) Open(ctx context.Context, uri string) error {

	u, err := url.Parse(uri)

	if err != nil {
		return err
	}

	root := u.Path
	info, err := os.Stat(root)

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("root is not a directory")
	}

	r.root = root
	return nil
}

func (r *FSReader) Read(ctx context.Context, path string) (io.ReadCloser, error) {

	abs_path := r.URI(path)

	_, err := os.Stat(abs_path)

	if err != nil {
		return nil, err
	}

	return os.Open(abs_path)
}

func (r *FSReader) URI(path string) string {
	return filepath.Join(r.root, path)
}
