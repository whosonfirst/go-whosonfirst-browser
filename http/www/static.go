package www

import (
	aa_static "github.com/aaronland/go-http-static"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/static"
	"net/http"
)

// BrowserOptions provides a list of JavaScript and CSS link to include with HTML output.
type BrowserOptions struct {
	JS             []string
	CSS            []string
	DataAttributes map[string]string
}

// Return a *BrowserOptions struct with default paths and URIs.
func DefaultBrowserOptions() *BrowserOptions {

	opts := &BrowserOptions{
		CSS:            []string{},
		JS:             []string{},
		DataAttributes: make(map[string]string),
	}

	return opts
}

func AppendStaticAssetHandlers(mux *http.ServeMux) error {
	return aa_static.AppendStaticAssetHandlers(mux, static.FS)
}

func AppendStaticAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {
	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Browser JavaScript and CSS files and related assets.
func AppendResourcesHandler(next http.Handler, opts *BrowserOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Browser JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next http.Handler, opts *BrowserOptions, prefix string) http.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.CSS = opts.CSS
	static_opts.JS = opts.JS
	static_opts.DataAttributes = opts.DataAttributes

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
}
