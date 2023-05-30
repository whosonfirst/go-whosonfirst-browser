package browser

import (
	_ "github.com/aaronland/go-http-server-tsnet"
	_ "github.com/aaronland/gocloud-blob-s3"
	_ "github.com/whosonfirst/go-reader-findingaid"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
)

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"

	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-server/handler"
	aa_log "github.com/aaronland/go-log/v2"
	wasm_exec "github.com/sfomuseum/go-http-wasm/v2"
	"github.com/sfomuseum/go-template/html"
	"github.com/sfomuseum/go-template/text"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/api"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/templates/javascript"
	wasm_placetypes "github.com/whosonfirst/go-whosonfirst-placetypes-wasm/http"
	wasm_validate "github.com/whosonfirst/go-whosonfirst-validate-wasm/http"
)

func Run(ctx context.Context, logger *log.Logger) error {

	fs, err := DefaultFlagSet(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create default flagset, %w", err)
	}

	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	cfg, err := ConfigFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive config from flagset, %w", err)
	}

	if config.Verbose {
		aa_log.SetMinLevelWithPrefix(aa_log.DEBUG_PREFIX)
	} else {
		aa_log.SetMinLevelWithPrefix(aa_log.INFO_PREFIX)
	}

	// Set up uris_table here

	t1 := time.Now()

	route_handlers := make(map[string]handler.RouteHandlerFunc)

	// Start setting up handlers

	// aa_log.Debug(logger, "Handle ping endpoint at %s\n", uris_table.Ping)
	route_handlers[uris_table.Ping] = pingHandlerFunc

	if capabilities.PNG {

		route_handlers[uris_table.PNG] = pngHandlerFunc

		for _, alt_path := range uris_table.PNGAlt {
			route_handlers[alt_path] = pngHandlerFunc
		}
	}

	if capabilities.SVG {

		route_handlers[uris_table.SVG] = svgHandlerFunc

		for _, alt_path := range uris_table.SVGAlt {
			route_handlers[alt_path] = svgHandlerFunc
		}
	}

	if capabilities.SPR {

		route_handlers[uris_table.SPR] = sprHandlerFunc

		for _, alt_path := range uris_table.SPRAlt {
			route_handlers[alt_path] = sprHandlerFunc
		}
	}

	if capabilities.GeoJSON {

		route_handlers[uris_table.GeoJSON] = geojsonHandlerFunc

		for _, alt_path := range uris_table.GeoJSONAlt {
			route_handlers[alt_path] = sprHandlerFunc
		}
	}

	if capabilities.GeoJSONLD {

		route_handlers[uris_table.GeoJSONLD] = geojsonLDHandlerFunc

		for _, alt_path := range uris_table.GeoJSONLDAlt {
			route_handlers[alt_path] = geojsonLDHandlerFunc
		}
	}

	if capabilities.NavPlace {

		route_handlers[uris_table.NavPlace] = navPlaceHandlerFunc

		for _, alt_path := range uris_table.NavPlaceAlt {
			route_handlers[alt_path] = navPlaceHandlerFunc
		}
	}

	if capabilities.Select {

		route_handlers[uris_table.Select] = selectHandlerFunc

		for _, alt_path := range uris_table.SelectAlt {
			route_handlers[alt_path] = selectHandlerFunc
		}
	}

	if capabilities.WebFinger {

		route_handlers[uris_table.Select] = webFingerHandlerFunc

		for _, alt_path := range uris_table.WebFingerAlt {
			route_handlers[alt_path] = webFingerHandlerFunc
		}
	}

	if capabilities.SearchAPI {

		route_handlers[uris_table.SearchAPI] = apiSearchHandlerFunc
	}

	if capabilities.PointInPolygonAPI {
		// To do: Need to sort out what's necessary to create *spatial_app.SpatialApplication
		// https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/app/app.go#L20
	}

	// Common code for HTML handler (public and/or edit handlers)

	if settings.HasHTMLCapabilities() {

		setupStaticOnce.Do(setupStatic)

		if setupStaticError != nil {
			return fmt.Errorf("Failed to configure static setup, %w", setupStaticError)
		}

		err = bootstrap.AppendAssetHandlers(mux, bootstrap_opts)

		if err != nil {
			return fmt.Errorf("Failed to append Bootstrap asset handlers, %w", err)
		}

		err = www.AppendAssetHandlers(mux, www_opts)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %w", err)
		}

		err = settings.CustomChrome.AppendStaticAssetHandlersWithPrefix(mux, uris_table.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append custom asset handlers, %w", err)
		}

		// Final map stuff

		maps_opts = maps.DefaultMapsOptions()
		maps_opts.AppendJavaScriptAtEOF = settings.JavaScriptAtEOF
		maps_opts.RollupAssets = capabilities.RollupAssets
		maps_opts.Prefix = uris_table.URIPrefix
		maps_opts.Logger = logger

		err = maps.AppendAssetHandlers(mux, maps_opts)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %v")
		}

		err = settings.MapProvider.AppendAssetHandlers(mux)

		if err != nil {
			return fmt.Errorf("Failed to append provider asset handlers, %v", err)
		}

		// Null handler for annoying things like favicons

		null_handler := www.NewNullHandler()

		favicon_path := filepath.Join(uris_table.Id, "favicon.ico")
		mux.Handle(favicon_path, null_handler)
	}

	// Public HTML handlers

	// To do: Consider hooks to require auth?

	if capabilities.Index {

		route_handlers[uris_table.Index] = wwwIndexHandlerFunc
	}

	if capabilities.Id {

		route_handlers[uris_table.Id] = wwwIdHandlerFunc

		// TBD...

		if capabilities.RollupAssets {
			err = www.AppendAssetHandlers(mux, www_opts.WithIdHandlerAssets())

			if err != nil {
				return fmt.Errorf("Failed to append asset handler for ID handler, %w", err)
			}
		}
	}

	if capabilities.Search {
		route_handlers[uris_table.Search] = wwwSearchHandlerFunc
	}

	// Edit/write HTML handlers

	if capabilities.EditUI {

		handler_routes[uris_table.EditGeometry] = wwwEditGeometryHandlerFunc

		// TBD...

		if capabilities.RollupAssets {

			err = www.AppendAssetHandlers(mux, www_opts.WithGeometryHandlerAssets())

			if err != nil {
				return fmt.Errorf("Failed to append asset handler for geometry handler, %w", err)
			}
		}

		// START OF wasm stuff

		wasm_exec_opts := wasm_exec.DefaultWASMOptions()
		wasm_exec_opts.AppendJavaScriptAtEOF = settings.JavaScriptAtEOF
		wasm_exec_opts.RollupAssets = capabilities.RollupAssets
		wasm_exec_opts.Prefix = uris_table.URIPrefix
		wasm_exec_opts.Logger = logger

		err = wasm_exec.AppendAssetHandlers(mux, wasm_exec_opts)

		if err != nil {
			return fmt.Errorf("Failed to append wasm asset handlers, %w", err)
		}

		err = wasm_placetypes.AppendAssetHandlersWithPrefix(mux, uris_table.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append wasm placetypes asset handlers, %w", err)
		}

		err = wasm_validate.AppendAssetHandlersWithPrefix(mux, uris_table.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append wasm validate asset handlers, %w", err)
		}

		// START OF I don't like having to do this but since the default 'whosonfirst.validate.feature.js'
		// package (in go-whosonfirst-validate-wasm) has a relative path and, importantly, no well-defined
		// way to specify the wasm path (yet) this is what we're going to do in conjunction with writing
		// our own 'whosonfirst.browser.validate' package that fetches the custom URI from 'whosonfirst.browser.uris'.
		// I suppose it would be easy enough to add a 'setWasmURI' or a 'setWasmPrefix' method to 'whosonfirst.validate.feature.js'
		// but today that hasn't happened.

		// This file is served by the http/www/static.go handlers
		wasm_validate_uri := "/wasm/validate_feature.wasm"
		wasm_placetypes_uri := "/wasm/whosonfirst_placetypes.wasm"

		if uris_table.URIPrefix != "" {

			validate_uri, err := url.JoinPath(uris_table.URIPrefix, wasm_validate_uri)

			if err != nil {
				return fmt.Errorf("Failed to assign URI prefix to validate wasm path, %w", err)
			}

			placetypes_uri, err := url.JoinPath(uris_table.URIPrefix, wasm_placetypes_uri)

			if err != nil {
				return fmt.Errorf("Failed to assign URI prefix to placetypes wasm path, %w", err)
			}

			wasm_validate_uri = validate_uri
			wasm_placetypes_uri = placetypes_uri
		}

		uris_table.AddCustomURI("validate_wasm", wasm_validate_uri)
		uris_table.AddCustomURI("placetypes_wasm", wasm_placetypes_uri)

		// END OF I don't like having	to do this

		handler_routes[uris_table.CreateFeature] = wwwCreateGeometryHandlerFunc

		// TBD
		if capabilities.RollupAssets {
			err = www.AppendAssetHandlers(mux, www_opts.WithCreateHandlerAssets())

			if err != nil {
				return fmt.Errorf("Failed to append asset handler for create handler, %w", err)
			}
		}

		// Custom client-side WASM/JS validation (optional)

		if settings.CustomEditValidationWasm != nil {

			aa_log.Debug(logger, "Enable custom edit validation WASM")

			validation_handler_opts := &www.CustomValidationWasmHandlerOptions{
				CustomValidationWasm: settings.CustomEditValidationWasm,
			}

			validation_handler, err := www.CustomValidationWasmHandler(validation_handler_opts)

			if err != nil {
				return fmt.Errorf("Failed to create custom edit validation WASM handler, %w", err)
			}

			fname := filepath.Base(settings.CustomEditValidationWasm.Path)
			custom_validate_uri := filepath.Join("/wasm", fname)

			if uris_table.URIPrefix != "" {

				uri, err := url.JoinPath(uris_table.URIPrefix, custom_validate_uri)

				if err != nil {
					return fmt.Errorf("Failed to assign URI prefix to custom validate wasm path, %w", err)
				}

				custom_validate_uri = uri
			}

			uris_table.AddCustomURI("custom_validate_wasm", custom_validate_uri)

			aa_log.Debug(logger, "Handle custom validation WASM endpoint at %s\n", custom_validate_uri)
			mux.Handle(custom_validate_uri, validation_handler)
		}
	}

	if capabilities.EditAPI {

		handler_routes[uris_table.DeprecateFeatureAPI] = apiDeprecateHandlerFunc
		handler_routes[uris_table.CessateFeatureAPI] = apiCessateHandlerFunc
		handler_routes[uris_table.EditGeometryAPI] = apiEditGeometryHandlerFunc
		handler_routes[uris_table.CreateFeatureAPI] = apiCreateFeatureHandlerFunc
	}

	// START OF TBD...

	/*
		for uri, h := range settings.CustomWWWHandlers {

			parts := strings.Split(uri, "#")

			path := parts[0]
			label := path

			if len(parts) == 2 {
				label = parts[1]
			}

			h = www.AppendResourcesHandler(h, www_opts)
			h = maps.AppendResourcesHandlerWithProvider(h, settings.MapProvider, maps_opts)
			h = bootstrap.AppendResourcesHandler(h, bootstrap_opts)
			h = settings.CustomChrome.WrapHandler(h, label)
			h = settings.Authenticator.WrapHandler(h)

			aa_log.Debug(logger, "Handle custom www endpoint at %s\n", path)
			mux.Handle(path, h)
		}

		for path, h := range settings.CustomAPIHandlers {

			h = settings.Authenticator.WrapHandler(h)

			aa_log.Debug(logger, "Handle custom API endpoint at %s\n", path)
			mux.Handle(path, h)
		}

		if len(settings.CustomAssetHandlerFunctions) > 0 {

			aa_log.Debug(logger, "Register custom asset handlers")

			for idx, handler_func := range settings.CustomAssetHandlerFunctions {

				err := handler_func(mux, uris_table.URIPrefix)

				if err != nil {
					return fmt.Errorf("Failed to register custom asset handler function at offset %d, %w", idx, err)
				}
			}
		}

	*/

	// END OF TBD...

	// START OF uris.js
	// Publish uris_table as a JS file so that other JS knows where to find things

	if err != nil {
		return fmt.Errorf("Failed to create URIs handler, %w", err)
	}

	uris_path := "/javascript/whosonfirst.browser.uris.js"

	if uris_table.URIPrefix != "" {

		path, err := url.JoinPath(uris_table.URIPrefix, uris_path)

		if err != nil {
			return fmt.Errorf("Failed to assign URI prefix to %s, %w", uris_path, err)
		}

		uris_path = path
	}

	route_handlers[uris_path] = urisHandlerFunc

	// END OF uris.js

	aa_log.Info(logger, "Time to set up: %v\n", time.Since(t1))

	route_handler, err := handler.RouteHandler(route_handlers)

	if err != nil {
		return fmt.Errorf("Failed to create route handler, %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", route_handler)

	// Finally, start the server

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new search for '%s', %w", server_uri, err)
	}

	aa_log.Info(logger, "Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}

/*
func appendCustomMiddlewareHandlers(settings *Settings, path string, handler http.Handler) http.Handler {

	custom_handlers, exists := settings.CustomMiddlewareHandlers[path]

	if !exists {
		return handler
	}

	for _, middleware_func := range custom_handlers {
		handler = middleware_func(handler)
	}

	return handler
}
*/
