package www

import (
	"github.com/aaronland/go-http-rewrite"
	_ "log"
	"net/http"
	"path/filepath"
	"strings"
)

func StaticFileSystem() (http.FileSystem, error) {
	fs := assetFS()
	return fs, nil
}

func StaticAssetsHandler() (http.Handler, error) {

	fs := assetFS()
	return http.FileServer(fs), nil
}

func StaticAssetsHandlerWithPrefix(prefix string) (http.Handler, error) {

	fs_handler, err := StaticAssetsHandler()

	if err != nil {
		return nil, err
	}

	prefix = strings.TrimRight(prefix, "/")

	if prefix == "" {
		return fs_handler, nil
	}

	rewrite_func := func(req *http.Request) (*http.Request, error) {
		req.URL.Path = strings.Replace(req.URL.Path, prefix, "", 1)
		return req, nil
	}

	rewrite_handler := rewrite.RewriteRequestHandler(fs_handler, rewrite_func)
	return rewrite_handler, nil
}

func AppendStaticAssetHandlers(mux *http.ServeMux) error {
	return AppendStaticAssetHandlersWithPrefix(mux, "")
}

func AppendStaticAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {

	asset_handler, err := StaticAssetsHandlerWithPrefix(prefix)

	if err != nil {
		return nil
	}

	for _, path := range AssetNames() {

		path := strings.Replace(path, "static", "", 1)

		if prefix != "" {
			path = appendPrefix(prefix, path)
		}

		mux.Handle(path, asset_handler)
	}

	return nil
}

func appendPrefix(prefix string, path string) string {

	prefix = strings.TrimRight(prefix, "/")

	if prefix != "" {
		path = strings.TrimLeft(path, "/")
		path = filepath.Join(prefix, path)
	}

	return path
}
