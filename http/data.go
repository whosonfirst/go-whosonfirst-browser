package http

import (
	"errors"
	gohttp "net/http"
	"os"
)

func DataHandler(root string) (gohttp.Handler, error) {

	info, err := os.Stat(root)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("Not a directory")
	}

	fs := gohttp.FileSystem(gohttp.Dir(root))
	return gohttp.FileServer(fs), nil
}
