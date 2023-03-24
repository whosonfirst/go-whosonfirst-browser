package www

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	aa_static "github.com/aaronland/go-http-static"
	aa_log "github.com/aaronland/go-log/v2"
	"github.com/sfomuseum/go-http-rollup"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/static"
)

type Foo struct {
	JS  []string
	CSS []string
}

var foo = map[string]Foo{

	// Core JS/CSS for all HTML pages
	"whosonfirst.browser": Foo{
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
		CSS: []string{
			"/css/whosonfirst.www.css",
			"/css/whosonfirst.common.css",
			"/css/whosonfirst.browser.css",
		},
	},
	// JS/CSS assets for the /create (edit) endpoint
	"whosonfirst.browser.create": Foo{
		JS: []string{
			"/javascript/whosonfirst.browser.api.js",
			"/javascript/whosonfirst.browser.leaflet.js",
			"/javascript/whosonfirst.webcomponent.existentialflag.js",
			"/javascript/whosonfirst.webcomponent.placetype.js",
			"/javascript/whosonfirst.browser.create.js",
			"/javascript/whosonfirst.browser.create.init.js",
		},
		CSS: []string{
			"/css/whosonfirst.browser.edit.css",
		},
	},

	// JS/CSS assets for the /geometry (edit) endpoint
	"whosonfirst.browser.geometry": Foo{
		JS: []string{
			"/javascript/whosonfirst.browser.api.js",
			"/javascript/whosonfirst.browser.leaflet.js",
			"/javascript/whosonfirst.browser.geometry.js",
			"/javascript/whosonfirst.browser.geometry.init.js",
		},
		CSS: []string{
			"/css/whosonfirst.browser.edit.css",
			"/css/whosonfirst.browser.edit.geometry.css",
		},
	},

	// JS/CSS assets for the /id endpoint
	"whosonfirst.browser.id": Foo{
		CSS: []string{
			"/css/whosonfirst.browser.id.css",
		},
		JS: []string{
			"/javascript/whosonfirst.browser.id.js",
			"/javascript/whosonfirst.browser.id.init.js",
		},
	},
}

// BrowserOptions provides a list of JavaScript and CSS link to include with HTML output.
type BrowserOptions struct {
	JS             []string
	CSS            []string
	DataAttributes map[string]string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
	RollupAssets          bool
	assets                []string
	JSRollupURI           []string
	CSSRollupURI          []string
	Prefix                string
	Logger                *log.Logger
}

// Return a *BrowserOptions struct with default paths and URIs.
func DefaultBrowserOptions() *BrowserOptions {

	logger := log.New(io.Discard, "", 0)

	opts := &BrowserOptions{
		assets: []string{
			"whosonfirst.browser",
		},
		DataAttributes: make(map[string]string),
		Logger:         logger,
	}

	return opts
}

func (opts *BrowserOptions) Clone() *BrowserOptions {

	assets := make([]string, len(opts.assets))
	attrs := make(map[string]string)

	for idx, label := range opts.assets {
		assets[idx] = label
	}

	for k, v := range opts.DataAttributes {
		attrs[k] = v
	}

	new_opts := &BrowserOptions{
		Logger:                opts.Logger,
		AppendJavaScriptAtEOF: opts.AppendJavaScriptAtEOF,
		RollupAssets:          opts.RollupAssets,
		Prefix:                opts.Prefix,
		DataAttributes:        attrs,
		assets:                assets,
	}

	return new_opts
}

func (opts *BrowserOptions) WithIdHandlerResources() *BrowserOptions {

	new_opts := opts.Clone()

	new_opts.assets = append(new_opts.assets, "whosonfirst.browser.id")
	return new_opts
}

func (opts *BrowserOptions) WithIdHandlerAssets() *BrowserOptions {

	new_opts := opts.Clone()

	new_opts.assets = []string{
		"whosonfirst.browser.id",
	}

	return new_opts
}

func (opts *BrowserOptions) WithCreateHandlerResources() *BrowserOptions {

	new_opts := opts.Clone()

	new_opts.assets = append(new_opts.assets, "whosonfirst.browser.create")
	return new_opts
}

func (opts *BrowserOptions) WithCreateHandlerAssets() *BrowserOptions {

	new_opts := opts.Clone()

	new_opts.assets = []string{
		"whosonfirst.browser.create",
	}

	return new_opts
}

func (opts *BrowserOptions) WithGeometryHandlerResources() *BrowserOptions {

	new_opts := opts.Clone()

	new_opts.assets = append(new_opts.assets, "whosonfirst.browser.geometry")
	return new_opts
}

func (opts *BrowserOptions) WithGeometryHandlerAssets() *BrowserOptions {

	new_opts := opts.Clone()

	new_opts.assets = []string{
		"whosonfirst.browser.geometry",
	}

	return new_opts
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Browser JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandler(next http.Handler, opts *BrowserOptions) http.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	static_opts.DataAttributes = opts.DataAttributes

	if opts.RollupAssets {

		static_opts.CSS = make([]string, len(opts.assets))
		static_opts.JS = make([]string, len(opts.assets))

		for idx, label := range opts.assets {

			if len(foo[label].CSS) > 0 {
				css_uri := fmt.Sprintf("/css/%s.rollup.css", label)
				static_opts.CSS[idx] = css_uri
			}

			if len(foo[label].JS) > 0 {
				js_uri := fmt.Sprintf("/javascript/%s.rollup.js", label)
				static_opts.JS[idx] = js_uri
			}
		}

	} else {

		static_opts.CSS = make([]string, 0)
		static_opts.JS = make([]string, 0)

		for _, label := range opts.assets {

			for _, uri := range foo[label].CSS {
				static_opts.CSS = append(static_opts.CSS, uri)
			}

			for _, uri := range foo[label].JS {
				static_opts.JS = append(static_opts.JS, uri)
			}

		}
	}

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, opts.Prefix)
}

// Append all the files in the net/http FS instance containing the embedded Browser assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux, opts *BrowserOptions) error {

	if !opts.RollupAssets {
		return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, opts.Prefix)
	}

	for _, label := range opts.assets {

		rollup_js_uri := fmt.Sprintf("/javascript/%s.rollup.js", label)
		rollup_css_uri := fmt.Sprintf("/css/%s.rollup.css", label)

		js_label := filepath.Base(rollup_js_uri)
		css_label := filepath.Base(rollup_css_uri)

		js_paths := make([]string, len(foo[label].JS))
		css_paths := make([]string, len(foo[label].CSS))

		for idx, path := range foo[label].JS {
			path = strings.TrimLeft(path, "/")

			aa_log.Debug(opts.Logger, "Add %s to JS rollup %s", path, rollup_js_uri)
			js_paths[idx] = path
		}

		for idx, path := range foo[label].CSS {
			path = strings.TrimLeft(path, "/")

			aa_log.Debug(opts.Logger, "Add %s to CSS rollup %s", path, rollup_css_uri)
			css_paths[idx] = path
		}

		rollup_js_paths := map[string][]string{
			js_label: js_paths,
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
			css_label: css_paths,
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

		if opts.Prefix != "" {

			u, err := url.JoinPath(opts.Prefix, rollup_css_uri)

			if err != nil {
				return fmt.Errorf("Failed to append prefix to %s, %w", rollup_css_uri, err)
			}

			rollup_css_uri = u
		}

		mux.Handle(rollup_css_uri, rollup_css_handler)
	}

	return nil
}
