package protomaps

import (
	"fmt"
	"github.com/aaronland/go-http-leaflet"
	"github.com/aaronland/go-http-rewrite"
	"github.com/sfomuseum/go-http-protomaps/static"
	"io/fs"
	_ "log"
	"net/http"
	"path/filepath"
	"strings"
)

// By default the go-http-protomaps package will also include and reference Leaflet.js resources using the aaronland/go-http-leaflet package. If you want or need to disable this behaviour set this variable to false.
var APPEND_LEAFLET_RESOURCES = true

// By default the go-http-protomaps package will also include and reference Leaflet.js assets using the aaronland/go-http-leaflet package. If you want or need to disable this behaviour set this variable to false.
var APPEND_LEAFLET_ASSETS = true

// ProtomapsOptions provides a list of JavaScript and CSS link to include with HTML output as well as a URL referencing a specific Protomaps PMTiles database to include a data attribute.
type ProtomapsOptions struct {
	// A list of relative JavaScript files to reference in one or more <script> tags
	JS []string
	// A list of relative CSS files to reference in one or more <link rel="stylesheet"> tags
	CSS []string
	// A URL for a specific PMTiles database to include as a 'data-protomaps-tile-url' attribute on the <body> tag.
	TileURL string
	// A leaflet.LeafletOptions struct
	LeafletOptions *leaflet.LeafletOptions
}

// Return a *ProtomapsOptions struct with default paths and URIs.
func DefaultProtomapsOptions() *ProtomapsOptions {

	leaflet_opts := leaflet.DefaultLeafletOptions()

	opts := &ProtomapsOptions{
		CSS: []string{},
		JS: []string{
			"/javascript/protomaps.js",
		},
		LeafletOptions: leaflet_opts,
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Protomaps JavaScript files and related assets.
func AppendResourcesHandler(next http.Handler, opts *ProtomapsOptions) http.Handler {

	if APPEND_LEAFLET_RESOURCES {
		next = leaflet.AppendResourcesHandler(next, opts.LeafletOptions)
	}

	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Protomaps JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next http.Handler, opts *ProtomapsOptions, prefix string) http.Handler {

	if APPEND_LEAFLET_RESOURCES {
		next = leaflet.AppendResourcesHandlerWithPrefix(next, opts.LeafletOptions, prefix)
	}

	js := opts.JS
	css := opts.CSS

	attrs := map[string]string{
		"protomaps-tile-url": opts.TileURL,
	}

	if prefix != "" {

		for i, path := range js {
			js[i] = appendPrefix(prefix, path)
		}

		for i, path := range css {
			css[i] = appendPrefix(prefix, path)
		}

		for k, path := range attrs {

			if strings.HasSuffix(k, "-url") && !strings.HasPrefix(path, "http") {
				attrs[k] = appendPrefix(prefix, path)
			}
		}
	}

	ext_opts := &rewrite.AppendResourcesOptions{
		JavaScript:     js,
		Stylesheets:    css,
		DataAttributes: attrs,
	}

	return rewrite.AppendResourcesHandler(next, ext_opts)
}

// AssetsHandler returns a net/http FS instance containing the embedded Protomaps assets that are included with this package.
func AssetsHandler() (http.Handler, error) {
	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs), nil
}

// AssetsHandler returns a net/http FS instance containing the embedded Protomaps assets that are included with this package ensuring that all URLs are stripped of prefix.
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

	if APPEND_LEAFLET_ASSETS {

		err := leaflet.AppendAssetHandlersWithPrefix(mux, prefix)

		if err != nil {
			return err
		}
	}

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

// FileHandlerFromPath will take a path and create a http.FileServer handler
// instance for the files in its root directory. The handler is returned with
// a relative URI for the filename in 'path' to be assigned to a net/http
// ServeMux instance.
func FileHandlerFromPath(path string, prefix string) (string, http.Handler, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return "", nil, fmt.Errorf("Failed to determine absolute path for '%s', %v", path, err)
	}

	fname := filepath.Base(abs_path)
	root := filepath.Dir(abs_path)

	tile_dir := http.Dir(root)
	tile_handler := http.FileServer(tile_dir)

	tile_url := fmt.Sprintf("/%s", fname)

	if prefix != "" {
		tile_handler = http.StripPrefix(prefix, tile_handler)
		tile_url = filepath.Join(prefix, fname)
	}

	return tile_url, tile_handler, nil
}

func appendPrefix(prefix string, path string) string {

	prefix = strings.TrimRight(prefix, "/")

	if prefix != "" {
		path = strings.TrimLeft(path, "/")
		path = filepath.Join(prefix, path)
	}

	return path
}
