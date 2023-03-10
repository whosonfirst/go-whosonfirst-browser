package leaflet

import (
	"github.com/aaronland/go-http-leaflet/static"
	aa_static "github.com/aaronland/go-http-static"
	_ "log"
	"net/http"
)

// LeafletOptions provides a list of JavaScript and CSS link to include with HTML output.
type LeafletOptions struct {
	JS             []string
	CSS            []string
	DataAttributes map[string]string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
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

// Append the Javascript and CSS URLs for the leaflet-geoman plugin.
// https://github.com/geoman-io/leaflet-geoman/
func (opts *LeafletOptions) EnableDraw() {
	opts.CSS = append(opts.CSS, "/css/leaflet-geoman.css")
	opts.JS = append(opts.JS, "/javascript/leaflet-geoman.min.js")
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
		DataAttributes: make(map[string]string),
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Leaflet JavaScript and CSS files and related assets.
func AppendResourcesHandler(next http.Handler, opts *LeafletOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Leaflet JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next http.Handler, opts *LeafletOptions, prefix string) http.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.CSS = opts.CSS
	static_opts.JS = opts.JS
	static_opts.DataAttributes = opts.DataAttributes
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
}

// Append all the files in the net/http FS instance containing the embedded Leaflet assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux) error {

	return aa_static.AppendStaticAssetHandlers(mux, static.FS)
}

// Append all the files in the net/http FS instance containing the embedded Leaflet assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
func AppendAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
