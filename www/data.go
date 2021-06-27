package www

import (
	"errors"
	"net/http"
	"os"
)

func DataHandler(root string) (http.Handler, error) {

	info, err := os.Stat(root)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("Not a directory")
	}

	fs := http.FileSystem(http.Dir(root))
	return http.FileServer(fs), nil
}
