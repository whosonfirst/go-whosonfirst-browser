package writer

import (
	"errors"
	"github.com/facebookgo/atomicfile"
	wof_writer "github.com/whosonfirst/go-whosonfirst-readwrite/writer"
	"io"
	"os"
	"path/filepath"
)

type FSWriter struct {
	wof_writer.Writer
	root      string
	dir_mode  os.FileMode
	file_mode os.FileMode
}

func NewFSWriter(root string) (wof_writer.Writer, error) {

	info, err := os.Stat(root)

	if err == nil && !info.IsDir() {
		return nil, errors.New("Target is not a directory")
	}

	w := FSWriter{
		root:      root,
		dir_mode:  0755,
		file_mode: 0644,
	}

	return &w, nil
}

func (w *FSWriter) Write(path string, fh io.ReadCloser) error {

	abs_path := w.URI(path)

	abs_root := filepath.Dir(abs_path)

	_, err := os.Stat(abs_root)

	if os.IsNotExist(err) {

		err = os.MkdirAll(abs_root, w.dir_mode)

		if err != nil {
			return err
		}
	}

	out, err := atomicfile.New(abs_path, w.file_mode)

	if err != nil {
		return err
	}

	_, err = io.Copy(out, fh)

	if err != nil {
		out.Abort()
		return err
	}

	return out.Close()
}

func (w *FSWriter) URI(path string) string {
	return filepath.Join(w.root, path)
}
