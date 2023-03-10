package tangramjs

import (
	"fmt"
	"net/http"

	"github.com/aaronland/go-http-leaflet"
	aa_static "github.com/aaronland/go-http-static"
	"github.com/aaronland/go-http-tangramjs/static"
)

// NEXTZEN_MVT_ENDPOINT is the default endpoint for Nextzen vector tiles
const NEXTZEN_MVT_ENDPOINT string = "https://tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt"

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
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
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

	if APPEND_LEAFLET_RESOURCES {
		opts.LeafletOptions.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF
		next = leaflet.AppendResourcesHandlerWithPrefix(next, opts.LeafletOptions, prefix)
	}

	attrs := map[string]string{
		"nextzen-api-key":   opts.NextzenOptions.APIKey,
		"nextzen-style-url": opts.NextzenOptions.StyleURL,
		"nextzen-tile-url":  opts.NextzenOptions.TileURL,
	}

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.CSS = opts.CSS
	static_opts.JS = opts.JS
	static_opts.DataAttributes = attrs
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
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
			return fmt.Errorf("Failed to append Leaflet assets, %w", err)
		}
	}

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
