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
	r := NewFSRader()
	Register("local", r)	// deprecated but backwards compatibility and all that
	Register("fs", r)	
}

type FSRader struct {
	Reader
	root string
}

func NewFSRader() Reader {

	r := FSRader{}
	return &r
}

func (r *FSRader) Open(ctx context.Context, uri string) error {

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

func (r *FSRader) Read(ctx context.Context, path string) (io.ReadCloser, error) {

	abs_path := r.URI(path)

	_, err := os.Stat(abs_path)

	if err != nil {
		return nil, err
	}

	return os.Open(abs_path)
}

func (r *FSRader) URI(path string) string {
	return filepath.Join(r.root, path)
}
