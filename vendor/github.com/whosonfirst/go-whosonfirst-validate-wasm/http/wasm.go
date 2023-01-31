package http

import (
	gohttp "net/http"

	aa_static "github.com/aaronland/go-http-static"	
	"github.com/whosonfirst/go-whosonfirst-validate-wasm/static"	
)

// WASMOptions provides a list of JavaScript and CSS link to include with HTML output.
type WASMOptions struct {
	JS  []string
	CSS []string
}

// Append the Javascript and CSS URLs for the wasm_exec.js script.
func (opts *WASMOptions) EnableWASMExec() {
	opts.JS = append(opts.JS, "/javascript/wasm_exec.js")
}

// Return a *WASMOptions struct with default paths and URIs.
func DefaultWASMOptions() *WASMOptions {

	opts := &WASMOptions{
		CSS: []string{},
		JS: []string{
			"/javascript/whosonfirst.validate.feature.js",
		},
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load WASM JavaScript and CSS files and related assets.
func AppendResourcesHandler(next gohttp.Handler, opts *WASMOptions) gohttp.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load WASM JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next gohttp.Handler, opts *WASMOptions, prefix string) gohttp.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.JS = opts.JS
	static_opts.CSS = opts.CSS
	
	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
}

// Append all the files in the net/http FS instance containing the embedded WASM assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *gohttp.ServeMux) error {

	return aa_static.AppendStaticAssetHandlers(mux, static.FS)
}

// Append all the files in the net/http FS instance containing the embedded WASM assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
func AppendAssetHandlersWithPrefix(mux *gohttp.ServeMux, prefix string) error {

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
