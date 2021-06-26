package writer

import (
	"context"
	"errors"
	"github.com/natefinch/atomic"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type FileWriter struct {
	Writer
	root      string
	dir_mode  os.FileMode
	file_mode os.FileMode
}

func init() {

	ctx := context.Background()

	schemes := []string{
		"fs",
		"file",
	}

	for _, scheme := range schemes {

		err := RegisterWriter(ctx, scheme, NewFileWriter)

		if err != nil {
			panic(err)
		}
	}
}

func NewFileWriter(ctx context.Context, uri string) (Writer, error) {

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

	// check for dir/file mode query parameters here

	wr := &FileWriter{
		dir_mode:  0755,
		file_mode: 0644,
		root:      root,
	}

	return wr, nil
}

func (wr *FileWriter) Write(ctx context.Context, path string, fh io.ReadSeeker) (int64, error) {

	abs_path := wr.WriterURI(ctx, path)
	abs_root := filepath.Dir(abs_path)

	tmp_file, err := os.CreateTemp("", filepath.Base(abs_path))

	if err != nil {
		return 0, err
	}

	tmp_path := tmp_file.Name()
	defer os.Remove(tmp_path)

	b, err := io.Copy(tmp_file, fh)

	if err != nil {
		return 0, err
	}

	err = tmp_file.Close()

	if err != nil {
		return 0, err
	}

	err = os.Chmod(tmp_path, wr.file_mode)

	if err != nil {
		return 0, err
	}

	_, err = os.Stat(abs_root)

	if os.IsNotExist(err) {

		err = os.MkdirAll(abs_root, wr.dir_mode)

		if err != nil {
			return 0, err
		}
	}

	err = atomic.ReplaceFile(tmp_path, abs_path)

	if err != nil {
		return 0, err
	}

	return b, nil
}

func (wr *FileWriter) WriterURI(ctx context.Context, path string) string {
	return filepath.Join(wr.root, path)
}

func (wr *FileWriter) Close(ctx context.Context) error {
	return nil
}
