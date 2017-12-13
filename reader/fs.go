package reader

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type FSReader struct {
	Reader
	root string
}

func NewFSReader(root string) (Reader, error) {

	info, err := os.Stat(root)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("root is not a directory")
	}

	s := FSReader{
		root: root,
	}

	return &s, nil
}

func (s *FSReader) Read(uri string) (io.ReadCloser, error) {

	path := filepath.Join(s.root, uri)

	_, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	return os.Open(path)
}
