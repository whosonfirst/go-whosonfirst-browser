package reader

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type FSReader struct {
	Reader
	root string
}

func init() {

	ctx := context.Background()
	err := RegisterReader(ctx, "fs", NewFSReader)

	if err != nil {
		panic(err)
	}
}

func NewFSReader(ctx context.Context, uri string) (Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	root := u.Path
	info, err := os.Stat(root)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("root is not a directory")
	}

	r := &FSReader{
		root: root,
	}

	return r, nil
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
