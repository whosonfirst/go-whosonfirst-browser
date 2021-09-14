package reader

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type FileReader struct {
	Reader
	root string
}

func init() {

	ctx := context.Background()

	err := RegisterReader(ctx, "fs", NewFileReader) // Deprecated

	if err != nil {
		panic(err)
	}

}

func NewFileReader(ctx context.Context, uri string) (Reader, error) {

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

	r := &FileReader{
		root: root,
	}

	return r, nil
}

func (r *FileReader) Read(ctx context.Context, path string) (io.ReadSeekCloser, error) {

	abs_path := r.ReaderURI(ctx, path)

	_, err := os.Stat(abs_path)

	if err != nil {
		return nil, err
	}

	return os.Open(abs_path)
}

func (r *FileReader) ReaderURI(ctx context.Context, path string) string {
	return filepath.Join(r.root, path)
}
