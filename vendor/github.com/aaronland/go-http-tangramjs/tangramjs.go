package tangramjs

import (
	"fmt"
	"github.com/aaronland/go-http-leaflet"
	"github.com/aaronland/go-http-rewrite"
	"github.com/aaronland/go-http-tangramjs/static"
	"io/fs"
	_ "log"
	"net/http"
	"path/filepath"
	"strings"
)

// NEXTZEN_MVT_ENDPOINT is the default endpoint for Nextzen vector tiles
const NEXTZEN_MVT_ENDPOINT string = "https://{s}.tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt"

// By default the go-http-tangramjs package will also include and reference Leaflet.js resources using the aaronland/go-http-leaflet package. If you want or need to disable this behaviour set this variable to false.
var APPEND_LEAFLET_RESOURCES = true

// By default the go-http-tangramjs package will also include and reference Leaflet.js assets using the aaronland/go-http-leaflet package. If you want or need to disable this behaviour set this variable to false.
var APPEND_LEAFLET_ASSETS = true

// NextzenOptions provides configuration variables for Nextzen map tiles.
type NextzenOptions struct {
	// A valid Nextzen developer API key
	APIKey string
	// The URL for a valid Tangram.js style.
	StyleURL string
	// The URL template to use for fetching Nextzen map tiles.
	TileURL string
}

// TangramJSOptions provides a list of JavaScript and CSS link to include with HTML output as well as options for Nextzen tiles and Leaflet.js.
type TangramJSOptions struct {
	// A list of Tangram.js Javascript files to append to HTML resources.
	JS []string
	// A list of Tangram.js CSS files to append to HTML resources.
	CSS []string
	// A NextzenOptions instance.
	NextzenOptions *NextzenOptions
	// A leaflet.LeafletOptions instance.
	LeafletOptions *leaflet.LeafletOptions
}

// Return a *NextzenOptions struct with default values.
func DefaultNextzenOptions() *NextzenOptions {

	opts := &NextzenOptions{
		APIKey:   "",
		StyleURL: "",
		TileURL:  NEXTZEN_MVT_ENDPOINT,
	}

	return opts
}

// Return a *TangramJSOptions struct with default values.
func DefaultTangramJSOptions() *TangramJSOptions {

	leaflet_opts := leaflet.DefaultLeafletOptions()
	nextzen_opts := DefaultNextzenOptions()

	opts := &TangramJSOptions{
		CSS: []string{},
		JS: []string{
			"/javascript/tangram.min.js",
		},
		LeafletOptions: leaflet_opts,
		NextzenOptions: nextzen_opts,
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Tangram.js files and related assets.
func AppendResourcesHandler(next http.Handler, opts *TangramJSOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Tangram.js files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next http.Handler, opts *TangramJSOptions, prefix string) http.Handler {

	js := opts.JS
	css := opts.CSS

	attrs := map[string]string{
		"nextzen-api-key":   opts.NextzenOptions.APIKey,
		"nextzen-style-url": opts.NextzenOptions.StyleURL,
		"nextzen-tile-url":  opts.NextzenOptions.TileURL,
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

	append_opts := &rewrite.AppendResourcesOptions{
		JavaScript:     js,
		Stylesheets:    css,
		DataAttributes: attrs,
	}

	if APPEND_LEAFLET_RESOURCES {
		next = leaflet.AppendResourcesHandlerWithPrefix(next, opts.LeafletOptions, prefix)
	}

	return rewrite.AppendResourcesHandler(next, append_opts)
}

// AssetsHandler returns a net/http FS instance containing the embedded Tangram.js assets that are included with this package.
func AssetsHandler() (http.Handler, error) {

	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs), nil
}

// AssetsHandler returns a net/http FS instance containing the embedded Tangram.js assets that are included with this package ensuring that all URLs are stripped of prefix.
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

// Append all the files in the net/http FS instance containing the embedded Tangram.js assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux) error {
	return AppendAssetHandlersWithPrefix(mux, "")
}

// Append all the files in the net/http FS instance containing the embedded Tangram.js assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
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

func appendPrefix(prefix string, path string) string {

	prefix = strings.TrimRight(prefix, "/")

	if prefix != "" {
		path = strings.TrimLeft(path, "/")
		path = filepath.Join(prefix, path)
	}

	return path
}
