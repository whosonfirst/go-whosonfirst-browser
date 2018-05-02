package writer

import (
       "fmt"
	"io"
	_ "log"
	"os"
)

type StdoutWriter struct {
	Writer
}

func NewStdoutWriter() (Writer, error) {

	w := StdoutWriter{}
	return &w, nil
}

func (w *StdoutWriter) Write(path string, fh io.ReadCloser) error {
	_, err := io.Copy(os.Stdout, fh)
	return err
}

func (w *StdoutWriter) URI(path string) string {
     return fmt.Sprintf("stdout://%s", path)
}
