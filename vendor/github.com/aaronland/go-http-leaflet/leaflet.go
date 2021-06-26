package leaflet

import (
	"fmt"
	"github.com/aaronland/go-http-leaflet/static"
	"github.com/aaronland/go-http-rewrite"
	"io/fs"
	_ "log"
	"net/http"
	"path/filepath"
	"strings"
)

// LeafletOptions provides a list of JavaScript and CSS link to include with HTML output.
type LeafletOptions struct {
	JS  []string
	CSS []string
}

// Append the Javascript and CSS URLs for the Leaflet.Fullscreen plugin.
func (opts *LeafletOptions) EnableFullscreen() {
	opts.CSS = append(opts.CSS, "/css/leaflet.fullscreen.css")
	opts.JS = append(opts.JS, "/javascript/leaflet.fullscreen.min.js")
}

// Append the Javascript and CSS URLs for the Leaflet.Hash plugin.
func (opts *LeafletOptions) EnableHash() {
	opts.JS = append(opts.JS, "/javascript/leaflet-hash.js")
}

// Append the Javascript and CSS URLs for the Leaflet.Draw plugin.
func (opts *LeafletOptions) EnableDraw() {
	opts.CSS = append(opts.CSS, "/css/leaflet.draw.css")
	opts.JS = append(opts.JS, "/javascript/leaflet.draw.js")
}

// Return a *LeafletOptions struct with default paths and URIs.
func DefaultLeafletOptions() *LeafletOptions {

	opts := &LeafletOptions{
		CSS: []string{
			"/css/leaflet.css",
		},
		JS: []string{
			"/javascript/leaflet.js",
		},
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Leaflet JavaScript files and related assets.
func AppendResourcesHandler(next http.Handler, opts *LeafletOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Leaflet JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next http.Handler, opts *LeafletOptions, prefix string) http.Handler {

	js := opts.JS
	css := opts.CSS

	if prefix != "" {

		for i, path := range js {
			js[i] = appendPrefix(prefix, path)
		}

		for i, path := range css {
			css[i] = appendPrefix(prefix, path)
		}
	}

	ext_opts := &rewrite.AppendResourcesOptions{
		JavaScript:  js,
		Stylesheets: css,
	}

	return rewrite.AppendResourcesHandler(next, ext_opts)
}

// AssetsHandler returns a net/http FS instance containing the embedded Leaflet assets that are included with this package.
func AssetsHandler() (http.Handler, error) {

	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs), nil
}

// AssetsHandler returns a net/http FS instance containing the embedded Leaflet assets that are included with this package ensuring that all URLs are stripped of prefix.
func AssetsHandlerWithPrefix(prefix string) (http.Handler, error) {

	fs_handler, err := AssetsHandler()

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

// Append all the files in the net/http FS instance containing the embedded Protomaps assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux) error {
	return AppendAssetHandlersWithPrefix(mux, "")
}

// Append all the files in the net/http FS instance containing the embedded Protomaps assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
func AppendAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {

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
