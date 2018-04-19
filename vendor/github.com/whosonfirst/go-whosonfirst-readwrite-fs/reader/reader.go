package reader

import (
	"errors"
	wof_reader "github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"io"
	"os"
	"path/filepath"
)

type FSReader struct {
	wof_reader.Reader
	root string
}

func NewFSReader(root string) (wof_reader.Reader, error) {

	info, err := os.Stat(root)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("root is not a directory")
	}

	r := FSReader{
		root: root,
	}

	return &r, nil
}

func (r *FSReader) Read(path string) (io.ReadCloser, error) {

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
