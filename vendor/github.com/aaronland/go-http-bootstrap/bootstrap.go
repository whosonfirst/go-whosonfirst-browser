package bootstrap

import (
	"net/http"

	"github.com/aaronland/go-http-bootstrap/static"
	aa_static "github.com/aaronland/go-http-static"
)

// BootstrapOptions provides a list of JavaScript and CSS link to include with HTML output.
type BootstrapOptions struct {
	// A list of relative Bootstrap Javascript URLs to append as resources in HTML output.
	JS []string
	// A list of relative Bootstrap CSS URLs to append as resources in HTML output.
	CSS []string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
}

// Return a *BootstrapOptions struct with default paths and URIs.
func DefaultBootstrapOptions() *BootstrapOptions {

	opts := &BootstrapOptions{
		CSS: []string{"/css/bootstrap.min.css"},
		JS:  make([]string, 0),
	}

	return opts
}

func (opts *BootstrapOptions) EnableJavascript() {
	opts.JS = append(opts.JS, "/javascript/bootstrap.bundle.min.js")
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Bootstrap JavaScript files and related assets.
func AppendResourcesHandler(next http.Handler, opts *BootstrapOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Bootstrap JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next http.Handler, opts *BootstrapOptions, prefix string) http.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.CSS = opts.CSS
	static_opts.JS = opts.JS
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
}

// Append all the files in the net/http FS instance containing the embedded Bootstrap assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux) error {

	return aa_static.AppendStaticAssetHandlers(mux, static.FS)
}

// Append all the files in the net/http FS instance containing the embedded Bootstrap assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
func AppendAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
