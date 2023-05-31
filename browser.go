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
	"time"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-server/handler"
	aa_log "github.com/aaronland/go-log/v2"
	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
	// wasm_placetypes "github.com/whosonfirst/go-whosonfirst-placetypes-wasm/http"
	// wasm_validate "github.com/whosonfirst/go-whosonfirst-validate-wasm/http"
)

func Run(ctx context.Context, run_logger *log.Logger) error {

	fs, err := DefaultFlagSet(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create default flagset, %w", err)
	}

	return RunWithFlagSet(ctx, fs, run_logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, run_logger *log.Logger) error {

	run_cfg, err := ConfigFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive config from flagset, %w", err)
	}

	return RunWithConfig(ctx, run_cfg, run_logger)
}

func RunWithConfig(ctx context.Context, run_cfg *Config, run_logger *log.Logger) error {

	cfg = run_cfg
	logger = run_logger

	if cfg.Verbose {
		aa_log.SetMinLevelWithPrefix(aa_log.DEBUG_PREFIX)
	} else {
		aa_log.SetMinLevelWithPrefix(aa_log.INFO_PREFIX)
	}

	// START OF auto-set/update config flags

	if cfg.EnableAll {
		cfg.EnableGraphics = true
		cfg.EnableData = true
		cfg.EnableHTML = true
	}

	if cfg.EnableGraphics {
		cfg.EnablePNG = true
		cfg.EnableSVG = true
	}

	if cfg.EnableData {
		cfg.EnableGeoJSON = true
		cfg.EnableGeoJSONLD = true
		cfg.EnableNavPlace = true
		cfg.EnableSPR = true
		cfg.EnableSelect = true
		cfg.EnableWebFinger = true
	}

	if cfg.EnableSearch {
		cfg.EnableSearchAPI = true
		cfg.EnableHTML = true
	}

	if cfg.EnablePointInPolygon {
		cfg.EnablePointInPolygonAPI = true
		cfg.EnableHTML = true
	}

	if cfg.EnableHTML {
		cfg.EnableGeoJSON = true
		cfg.EnablePNG = true
		cfg.EnableIndex = true
		cfg.EnableId = true
	}

	if cfg.EnableEdit {
		cfg.EnableEditAPI = true
		cfg.EnableEditUI = true
	}

	if cfg.EnableEditUI {
		cfg.EnableEditAPI = true
	}

	if cfg.DisableGeoJSON {
		cfg.EnableGeoJSON = false
	}

	if cfg.DisableGeoJSONLD {
		cfg.EnableGeoJSONLD = false
	}

	if cfg.DisableId {
		cfg.EnableId = false
	}

	if cfg.DisableIndex {
		cfg.EnableIndex = false
	}

	if cfg.DisableNavPlace {
		cfg.EnableNavPlace = false
	}

	if cfg.DisablePNG {
		cfg.EnablePNG = false
	}

	if cfg.DisableSearch {
		cfg.EnableSearch = false
	}

	if cfg.DisableSelect {
		cfg.EnableSelect = false
	}

	if cfg.DisableSPR {
		cfg.EnableSPR = false
	}

	if cfg.DisableSVG {
		cfg.EnableSVG = false
	}

	if cfg.DisableWebFinger {
		cfg.EnableWebFinger = false
	}

	// END OF auto-set/update config flags

	// START OF set up capabilities and uris

	capabilities = &browser_capabilities.Capabilities{}
	capabilities.RollupAssets = cfg.RollupAssets

	uris_table = &browser_uris.URIs{
		URIPrefix: cfg.URIPrefix,
		Ping:      cfg.PathPing,
	}

	if cfg.EnableIndex {
		capabilities.Index = true
		uris_table.Index = cfg.PathIndex
	}

	if cfg.EnableId {
		capabilities.Id = true
		uris_table.Id = cfg.PathId
	}

	if cfg.EnableGeoJSON {
		capabilities.GeoJSON = true
		uris_table.GeoJSON = cfg.PathGeoJSON
		uris_table.GeoJSONAlt = cfg.PathGeoJSONAlt
	}

	if cfg.EnableGeoJSONLD {
		capabilities.GeoJSONLD = true
		uris_table.GeoJSONLD = cfg.PathGeoJSONLD
		uris_table.GeoJSONLDAlt = cfg.PathGeoJSONLDAlt
	}

	if cfg.EnableSVG {
		capabilities.SVG = true
		uris_table.SVG = cfg.PathSVG
		uris_table.SVGAlt = cfg.PathSVGAlt
	}

	if cfg.EnablePNG {
		capabilities.PNG = true
		uris_table.PNG = cfg.PathPNG
		uris_table.PNGAlt = cfg.PathPNGAlt
	}

	if cfg.EnableSelect {
		capabilities.Select = true
		uris_table.Select = cfg.PathSelect
	}

	if cfg.EnableNavPlace {
		capabilities.NavPlace = true
		uris_table.NavPlace = cfg.PathNavPlace
		uris_table.NavPlaceAlt = cfg.PathNavPlaceAlt
	}

	if cfg.EnableSPR {
		capabilities.SPR = true
		uris_table.SPR = cfg.PathSPR
		uris_table.SPRAlt = cfg.PathSPRAlt
	}

	if cfg.EnableWebFinger {
		capabilities.WebFinger = true
		uris_table.WebFinger = cfg.PathWebFinger
		uris_table.WebFingerAlt = cfg.PathWebFingerAlt
	}

	if cfg.EnableEditUI {
		capabilities.EditUI = true
		capabilities.CreateFeature = true
		capabilities.DeprecateFeature = true
		capabilities.CessateFeature = true
		capabilities.EditGeometry = true
		uris_table.CreateFeature = cfg.PathCreateFeature
		uris_table.EditGeometry = cfg.PathEditGeometry
	}

	if cfg.EnableEditAPI {
		capabilities.EditAPI = true
		capabilities.CreateFeatureAPI = true
		capabilities.DeprecateFeatureAPI = true
		capabilities.CessateFeatureAPI = true
		capabilities.EditGeometryAPI = true
		uris_table.CreateFeatureAPI = cfg.PathCreateFeatureAPI
		uris_table.DeprecateFeatureAPI = cfg.PathDeprecateFeatureAPI
		uris_table.CessateFeatureAPI = cfg.PathCessateFeatureAPI
		uris_table.EditGeometryAPI = cfg.PathEditGeometryAPI
	}

	if cfg.EnableSearchAPI {
		capabilities.SearchAPI = true
		uris_table.Search = cfg.PathSearchAPI
	}

	if cfg.EnableSearch {
		capabilities.Search = true
		uris_table.Search = cfg.PathSearch
	}

	if cfg.EnablePointInPolygonAPI {
		capabilities.PointInPolygonAPI = true
		uris_table.Search = cfg.PathPointInPolygonAPI
	}

	if cfg.EnablePointInPolygon {
		capabilities.PointInPolygon = true
		uris_table.PointInPolygon = cfg.PathPointInPolygon
	}

	if cfg.URIPrefix != "" {

		err := uris_table.ApplyPrefix(cfg.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to apply prefix to Uris_table table, %w", err)
		}
	}

	// END OF set up capabilities and uris_table

	t1 := time.Now()

	mux := http.NewServeMux()

	route_handlers := make(map[string]handler.RouteHandlerFunc)

	// Start setting up handlers

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

	// Common code for HTML handler (public and/or edit handlers)

	if capabilities.HasHTMLCapabilities() {

		setupStaticOnce.Do(setupStatic)

		if setupStaticError != nil {
			return fmt.Errorf("Failed to configure static setup, %w", setupStaticError)
		}

		setupMapsOnce.Do(setupStatic)

		if setupMapsError != nil {
			setupWWWError = fmt.Errorf("Failed to configure static setup, %w", setupMapsError)
		}

		// START OF TBD...

		err := bootstrap.AppendAssetHandlers(mux, bootstrap_opts)

		if err != nil {
			return fmt.Errorf("Failed to append Bootstrap asset handlers, %w", err)
		}

		err = www.AppendAssetHandlers(mux, www_opts)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %w", err)
		}

		/*
			err = settings.CustomChrome.AppendStaticAssetHandlersWithPrefix(mux, uris_table.URIPrefix)

			if err != nil {
				return fmt.Errorf("Failed to append custom asset handlers, %w", err)
			}
		*/

		// Final map stuff

		maps_opts = maps.DefaultMapsOptions()
		maps_opts.AppendJavaScriptAtEOF = cfg.JavaScriptAtEOF
		maps_opts.RollupAssets = capabilities.RollupAssets
		maps_opts.Prefix = uris_table.URIPrefix
		maps_opts.Logger = logger

		err = maps.AppendAssetHandlers(mux, maps_opts)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %v")
		}

		err = map_provider.AppendAssetHandlers(mux)

		if err != nil {
			return fmt.Errorf("Failed to append provider asset handlers, %v", err)
		}

		// Null handler for annoying things like favicons
		// null_handler := www.NewNullHandler()
		// favicon_path := filepath.Join(uris_table.Id, "favicon.ico")
		// mux.Handle(favicon_path, null_handler)
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

			err := www.AppendAssetHandlers(mux, www_opts.WithIdHandlerAssets())

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

		route_handlers[uris_table.EditGeometry] = wwwEditGeometryHandlerFunc

		// TBD...

		if capabilities.RollupAssets {

			err := www.AppendAssetHandlers(mux, www_opts.WithGeometryHandlerAssets())

			if err != nil {
				return fmt.Errorf("Failed to append asset handler for geometry handler, %w", err)
			}
		}

		// START OF wasm stuff

		/*
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

			route_handlers[uris_table.CreateFeature] = wwwCreateGeometryHandlerFunc

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

		*/
	}

	if capabilities.EditAPI {

		route_handlers[uris_table.DeprecateFeatureAPI] = apiDeprecateHandlerFunc
		route_handlers[uris_table.CessateFeatureAPI] = apiCessateHandlerFunc
		route_handlers[uris_table.EditGeometryAPI] = apiEditGeometryHandlerFunc
		route_handlers[uris_table.CreateFeatureAPI] = apiCreateFeatureHandlerFunc
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
