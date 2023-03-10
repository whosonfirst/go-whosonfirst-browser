package protomaps

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/aaronland/go-http-leaflet"
	aa_static "github.com/aaronland/go-http-static"
	"github.com/sfomuseum/go-http-protomaps/static"
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
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
}

// Return a *ProtomapsOptions struct with default paths and URIs.
func DefaultProtomapsOptions() *ProtomapsOptions {

	leaflet_opts := leaflet.DefaultLeafletOptions()

	opts := &ProtomapsOptions{
		CSS: []string{},
		JS: []string{
			"/javascript/protomaps.min.js",
		},
		LeafletOptions: leaflet_opts,
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Protomaps JavaScript files and related assets.
func AppendResourcesHandler(next http.Handler, opts *ProtomapsOptions) http.Handler {

	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Protomaps JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next http.Handler, opts *ProtomapsOptions, prefix string) http.Handler {

	if APPEND_LEAFLET_RESOURCES {
		opts.LeafletOptions.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF
		next = leaflet.AppendResourcesHandlerWithPrefix(next, opts.LeafletOptions, prefix)
	}

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.CSS = opts.CSS
	static_opts.JS = opts.JS
	static_opts.DataAttributes["protomaps-tile-url"] = opts.TileURL
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
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
			return fmt.Errorf("Failed to append Leaflet assets, %w", err)
		}
	}

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
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
