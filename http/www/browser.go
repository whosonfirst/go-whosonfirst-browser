package www

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	aa_static "github.com/aaronland/go-http-static"
	"github.com/sfomuseum/go-http-rollup"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/static"
)

// BrowserOptions provides a list of JavaScript and CSS link to include with HTML output.
type BrowserOptions struct {
	JS             []string
	CSS            []string
	DataAttributes map[string]string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
	RollupAssets          bool
	Prefix                string
	Logger                *log.Logger
}

// Return a *BrowserOptions struct with default paths and URIs.
func DefaultBrowserOptions() *BrowserOptions {

	logger := log.New(io.Discard, "", 0)

	opts := &BrowserOptions{
		CSS: []string{
			"/css/whosonfirst.www.css",
			"/css/whosonfirst.common.css",
			"/css/whosonfirst.browser.css",
		},
		JS: []string{
			"/javascript/localforage.min.js",
			"/javascript/whosonfirst.www.js",
			"/javascript/whosonfirst.render.js",
			"/javascript/whosonfirst.properties.js",
			"/javascript/whosonfirst.cache.js",
			"/javascript/whosonfirst.uri.js",
			"/javascript/whosonfirst.net.js",
			"/javascript/whosonfirst.namify.js",
			"/javascript/whosonfirst.geojson.js",
			"/javascript/whosonfirst.leaflet.utils.js",
			"/javascript/whosonfirst.leaflet.styles.js",
			"/javascript/whosonfirst.leaflet.handlers.js",
			"/javascript/whosonfirst.browser.common.js",
			"/javascript/whosonfirst.browser.feedback.js",
			"/javascript/whosonfirst.browser.maps.js",
		},
		DataAttributes: make(map[string]string),
		Logger:         logger,
	}

	return opts
}

func (opts *BrowserOptions) WithIdHandlerResources() *BrowserOptions {

	css := []string{
		"/css/whosonfirst.browser.id.css",
	}

	js := []string{
		"/javascript/whosonfirst.browser.id.js",
		"/javascript/whosonfirst.browser.id.init.js",
	}

	new_opts := opts.Clone()

	for _, uri := range css {
		new_opts.CSS = append(new_opts.CSS, uri)
	}

	for _, uri := range js {
		new_opts.JS = append(new_opts.JS, uri)
	}

	return new_opts
}

func (opts *BrowserOptions) WithCreateHandlerResources() *BrowserOptions {

	css := []string{
		"/css/whosonfirst.browser.edit.css",
	}

	js := []string{
		"/javascript/whosonfirst.browser.api.js",
		"/javascript/whosonfirst.browser.leaflet.js",
		"/javascript/whosonfirst.webcomponent.existentialflag.js",
		"/javascript/whosonfirst.webcomponent.placetype.js",
		"/javascript/whosonfirst.browser.create.js",
		"/javascript/whosonfirst.browser.create.init.js",
	}

	new_opts := opts.Clone()

	for _, uri := range css {
		new_opts.CSS = append(new_opts.CSS, uri)
	}

	for _, uri := range js {
		new_opts.JS = append(new_opts.JS, uri)
	}

	return new_opts
}

func (opts *BrowserOptions) WithGeometryHandlerResources() *BrowserOptions {

	css := []string{
		"/css/whosonfirst.browser.edit.css",
		"/css/whosonfirst.browser.edit.geometry.css",
	}

	js := []string{
		"/javascript/whosonfirst.browser.api.js",
		"/javascript/whosonfirst.browser.leaflet.js",
		"/javascript/whosonfirst.browser.geometry.js",
		"/javascript/whosonfirst.browser.geometry.init.js",
	}

	new_opts := opts.Clone()

	for _, uri := range css {
		new_opts.CSS = append(new_opts.CSS, uri)
	}

	for _, uri := range js {
		new_opts.JS = append(new_opts.JS, uri)
	}

	return new_opts
}

func (opts *BrowserOptions) Clone() *BrowserOptions {

	css := make([]string, len(opts.CSS))
	js := make([]string, len(opts.JS))
	attrs := make(map[string]string)

	for idx, uri := range opts.CSS {
		css[idx] = uri
	}

	for idx, uri := range opts.JS {
		js[idx] = uri
	}

	for k, v := range opts.DataAttributes {
		attrs[k] = v
	}

	new_opts := &BrowserOptions{
		Logger:                opts.Logger,
		AppendJavaScriptAtEOF: opts.AppendJavaScriptAtEOF,
		RollupAssets:          opts.RollupAssets,
		Prefix:                opts.Prefix,
		CSS:                   css,
		JS:                    js,
		DataAttributes:        attrs,
	}

	return new_opts
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Browser JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandler(next http.Handler, opts *BrowserOptions) http.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.DataAttributes = opts.DataAttributes
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	if opts.RollupAssets {

		static_opts.CSS = []string{
			"/css/browser.rollup.css",
		}

		static_opts.JS = []string{
			"/javascript/browser.rollup.js",
		}

	} else {

		static_opts.CSS = opts.CSS
		static_opts.JS = opts.JS
	}

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, opts.Prefix)
}

// Append all the files in the net/http FS instance containing the embedded Browser assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux, opts *BrowserOptions) error {

	if !opts.RollupAssets {
		return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, opts.Prefix)
	}

	js_paths := make([]string, len(opts.JS))
	css_paths := make([]string, len(opts.CSS))

	for idx, path := range opts.JS {
		path = strings.TrimLeft(path, "/")
		js_paths[idx] = path
	}

	for idx, path := range opts.CSS {
		path = strings.TrimLeft(path, "/")
		css_paths[idx] = path
	}

	rollup_js_paths := map[string][]string{
		"browser.rollup.js": js_paths,
	}

	rollup_js_opts := &rollup.RollupJSHandlerOptions{
		FS:     static.FS,
		Paths:  rollup_js_paths,
		Logger: opts.Logger,
	}

	rollup_js_handler, err := rollup.RollupJSHandler(rollup_js_opts)

	if err != nil {
		return fmt.Errorf("Failed to create rollup JS handler, %w", err)
	}

	rollup_js_uri := "/javascript/browser.rollup.js"

	if opts.Prefix != "" {

		u, err := url.JoinPath(opts.Prefix, rollup_js_uri)

		if err != nil {
			return fmt.Errorf("Failed to append prefix to %s, %w", rollup_js_uri, err)
		}

		rollup_js_uri = u
	}

	mux.Handle(rollup_js_uri, rollup_js_handler)

	// CSS

	rollup_css_paths := map[string][]string{
		"browser.rollup.css": css_paths,
	}

	rollup_css_opts := &rollup.RollupCSSHandlerOptions{
		FS:     static.FS,
		Paths:  rollup_css_paths,
		Logger: opts.Logger,
	}

	rollup_css_handler, err := rollup.RollupCSSHandler(rollup_css_opts)

	if err != nil {
		return fmt.Errorf("Failed to create rollup CSS handler, %w", err)
	}

	rollup_css_uri := "/css/browser.rollup.css"

	if opts.Prefix != "" {

		u, err := url.JoinPath(opts.Prefix, rollup_css_uri)

		if err != nil {
			return fmt.Errorf("Failed to append prefix to %s, %w", rollup_css_uri, err)
		}

		rollup_css_uri = u
	}

	mux.Handle(rollup_css_uri, rollup_css_handler)
	return nil
}
