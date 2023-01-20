package www

import (
	"fmt"
	"net/http"
	"os"
)

func DataHandler(root string) (http.Handler, error) {

	info, err := os.Stat(root)

	if err != nil {
		return nil, fmt.Errorf("Failed to stat %s, %w", root, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}

	fs := http.FileSystem(http.Dir(root))
	return http.FileServer(fs), nil
}
