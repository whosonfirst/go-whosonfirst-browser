package maps

import (
	gohttp "net/http"

	"github.com/aaronland/go-http-maps/provider"
	"github.com/aaronland/go-http-maps/static"
	aa_static "github.com/aaronland/go-http-static"
)

// MapsOptions provides a list of JavaScript and CSS link to include with HTML output.
type MapsOptions struct {
	JS             []string
	CSS            []string
	DataAttributes map[string]string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
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

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.CSS = opts.CSS
	static_opts.JS = opts.JS
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
}

// Append all the files in the net/http FS instance containing the embedded Maps assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *gohttp.ServeMux) error {
	return AppendAssetHandlersWithPrefix(mux, "")
}

// Append all the files in the net/http FS instance containing the embedded Maps assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
func AppendAssetHandlersWithPrefix(mux *gohttp.ServeMux, prefix string) error {

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
