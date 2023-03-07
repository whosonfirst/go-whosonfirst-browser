package maps

import (
	"fmt"
	"io/fs"
	_ "log"
	gohttp "net/http"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-http-maps/provider"
	"github.com/aaronland/go-http-maps/static"
	"github.com/aaronland/go-http-rewrite"	
)

// MapsOptions provides a list of JavaScript and CSS link to include with HTML output.
type MapsOptions struct {
	JS             []string
	CSS            []string
	DataAttributes map[string]string
}

// Return a *MapsOptions struct with default paths and URIs.
func DefaultMapsOptions() *MapsOptions {

	opts := &MapsOptions{
		CSS: []string{
			"/css/aaronland.maps.css",
		},
		JS: []string{
			"/javascript/aaronland.maps.js",
		},
		DataAttributes: make(map[string]string),
	}

	return opts
}

func AppendResourcesHandlerWithPrefixAndProvider(next gohttp.Handler, map_provider provider.Provider, maps_opts *MapsOptions, prefix string) gohttp.Handler {
	next = map_provider.AppendResourcesHandlerWithPrefix(next, prefix)
	next = AppendResourcesHandlerWithPrefix(next, maps_opts, prefix)
	return next
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Maps JavaScript and CSS files and related assets.
func AppendResourcesHandler(next gohttp.Handler, opts *MapsOptions) gohttp.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Maps JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next gohttp.Handler, opts *MapsOptions, prefix string) gohttp.Handler {

	js := make([]string, len(opts.JS))
	css := make([]string, len(opts.CSS))

	for i, path := range opts.JS {
		js[i] = appendPrefix(prefix, path)
	}

	for i, path := range opts.CSS {
		css[i] = appendPrefix(prefix, path)
	}

	ext_opts := &rewrite.AppendResourcesOptions{
		JavaScript:     js,
		Stylesheets:    css,
		DataAttributes: opts.DataAttributes,
	}

	return rewrite.AppendResourcesHandler(next, ext_opts)
}

// AssetsHandler returns a net/http FS instance containing the embedded Maps assets that are included with this package.
func AssetsHandler() (gohttp.Handler, error) {

	http_fs := gohttp.FS(static.FS)
	return gohttp.FileServer(http_fs), nil
}

// AssetsHandler returns a net/http FS instance containing the embedded Maps assets that are included with this package ensuring that all URLs are stripped of prefix.
func AssetsHandlerWithPrefix(prefix string) (gohttp.Handler, error) {

	fs_handler, err := AssetsHandler()

	if err != nil {
		return nil, err
	}

	prefix = strings.TrimRight(prefix, "/")

	if prefix == "" {
		return fs_handler, nil
	}

	rewrite_func := func(req *gohttp.Request) (*gohttp.Request, error) {
		req.URL.Path = strings.Replace(req.URL.Path, prefix, "", 1)
		return req, nil
	}

	rewrite_handler := rewrite.RewriteRequestHandler(fs_handler, rewrite_func)
	return rewrite_handler, nil
}

// Append all the files in the net/http FS instance containing the embedded Maps assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *gohttp.ServeMux) error {
	return AppendAssetHandlersWithPrefix(mux, "")
}

// Append all the files in the net/http FS instance containing the embedded Maps assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
func AppendAssetHandlersWithPrefix(mux *gohttp.ServeMux, prefix string) error {

	asset_handler, err := AssetsHandlerWithPrefix(prefix)

	if err != nil {
		return nil
	}

	walk_func := func(path string, info fs.DirEntry, err error) error {

		if path == "." {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if prefix != "" {
			path = appendPrefix(prefix, path)
		}

		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}

		// log.Println("APPEND", path)

		mux.Handle(path, asset_handler)
		return nil
	}

	return fs.WalkDir(static.FS, ".", walk_func)
}

func appendPrefix(prefix string, path string) string {

	prefix = strings.TrimRight(prefix, "/")

	if prefix != "" {
		path = strings.TrimLeft(path, "/")
		path = filepath.Join(prefix, path)
	}

	return path
}
