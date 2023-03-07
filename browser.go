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
	"time"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	map_www "github.com/aaronland/go-http-maps/http/www"
	"github.com/aaronland/go-http-ping/v2"
	"github.com/aaronland/go-http-server"
	aa_log "github.com/aaronland/go-log"
	wasm_exec "github.com/sfomuseum/go-http-wasm"
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

	return RunWithConfig(ctx, cfg, logger)
}

func RunWithConfig(ctx context.Context, cfg *Config, logger *log.Logger) error {

	settings, err := SettingsFromConfig(ctx, cfg)

	if err != nil {
		return fmt.Errorf("Failed to create settings from config, %w", err)
	}

	return RunWithSettings(ctx, settings, logger)
}

func RunWithSettings(ctx context.Context, settings *Settings, logger *log.Logger) error {

	if settings.Verbose {
		aa_log.SetMinLevelWithPrefix(aa_log.DEBUG_PREFIX)
	} else {
		aa_log.SetMinLevelWithPrefix(aa_log.INFO_PREFIX)
	}

	t1 := time.Now()

	t, err := html.LoadTemplates(ctx, settings.Templates...)

	if err != nil {
		return fmt.Errorf("Failed to load templates, %w", err)
	}

	js_t, err := text.LoadTemplatesMatching(ctx, "*.js", javascript.FS)

	if err != nil {
		return fmt.Errorf("Failed to load JS templates, %w", err)
	}

	// Start setting up handlers

	mux := http.NewServeMux()

	ping_handler, err := ping.PingPongHandler()

	if err != nil {
		return fmt.Errorf("Failed to create ping handler, %w", err)
	}

	mux.Handle(settings.URIs.Ping, ping_handler)

	aa_log.Debug(logger, "Handle ping endpoint at %s\n", settings.URIs.Ping)

	if settings.Capabilities.PNG {

		aa_log.Debug(logger, "PNG support enabled")

		sizes := www.DefaultRasterSizes()

		png_opts := &www.RasterHandlerOptions{
			Sizes:  sizes,
			Format: "png",
			Reader: settings.Reader,
			Logger: logger,
		}

		png_handler, err := www.RasterHandler(png_opts)

		if err != nil {
			return fmt.Errorf("Failed to create raster/png handler, %w", err)
		}

		aa_log.Debug(logger, "Handle PNG endpoint at %s\n", settings.URIs.PNG)
		mux.Handle(settings.URIs.PNG, png_handler)

		for _, alt_path := range settings.URIs.PNGAlt {
			aa_log.Debug(logger, "Handle png endpoint at %s\n", alt_path)
			mux.Handle(alt_path, png_handler)
		}
	}

	if settings.Capabilities.SVG {

		aa_log.Debug(logger, "SVG support enabled")

		sizes := www.DefaultSVGSizes()

		svg_opts := &www.SVGHandlerOptions{
			Sizes:  sizes,
			Reader: settings.Reader,
			Logger: logger,
		}

		svg_handler, err := www.SVGHandler(svg_opts)

		if err != nil {
			return fmt.Errorf("Failed to create SVG handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			svg_handler = settings.CORSWrapper.Handler(svg_handler)
		}

		aa_log.Debug(logger, "Handle svg endpoint at %s\n", settings.URIs.SVG)
		mux.Handle(settings.URIs.SVG, svg_handler)

		for _, alt_path := range settings.URIs.SVGAlt {

			aa_log.Debug(logger, "Handle SVG endpoint at %s\n", alt_path)
			mux.Handle(alt_path, svg_handler)
		}
	}

	if settings.Capabilities.SPR {

		aa_log.Debug(logger, "SPR support enabled")

		spr_opts := &www.SPRHandlerOptions{
			Reader: settings.Reader,
			Logger: logger,
		}

		spr_handler, err := www.SPRHandler(spr_opts)

		if err != nil {
			return fmt.Errorf("Failed to create SPR handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			spr_handler = settings.CORSWrapper.Handler(spr_handler)
		}

		aa_log.Debug(logger, "Handle spr endpoint at %s\n", settings.URIs.SPR)
		mux.Handle(settings.URIs.SPR, spr_handler)

		for _, alt_path := range settings.URIs.SPRAlt {

			aa_log.Debug(logger, "Handle SPR endpoint at %s\n", alt_path)
			mux.Handle(alt_path, spr_handler)
		}
	}

	if settings.Capabilities.GeoJSON {

		aa_log.Debug(logger, "GeoJSON support enabled")

		geojson_opts := &www.GeoJSONHandlerOptions{
			Reader: settings.Reader,
			Logger: logger,
		}

		geojson_handler, err := www.GeoJSONHandler(geojson_opts)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			geojson_handler = settings.CORSWrapper.Handler(geojson_handler)
		}

		aa_log.Debug(logger, "Handle GeoJSON endpoint at %s\n", settings.URIs.GeoJSON)
		mux.Handle(settings.URIs.GeoJSON, geojson_handler)

		for _, alt_path := range settings.URIs.GeoJSONAlt {

			aa_log.Debug(logger, "Handle GeoJSON endpoint at %s\n", alt_path)
			mux.Handle(alt_path, geojson_handler)
		}
	}

	if settings.Capabilities.GeoJSONLD {

		aa_log.Debug(logger, "GeoJSON-LD support enabled")

		geojsonld_opts := &www.GeoJSONLDHandlerOptions{
			Reader: settings.Reader,
			Logger: logger,
		}

		geojsonld_handler, err := www.GeoJSONLDHandler(geojsonld_opts)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON-LD handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			geojsonld_handler = settings.CORSWrapper.Handler(geojsonld_handler)
		}

		aa_log.Debug(logger, "Handle GeoJSON-LD endpoint at %s\n", settings.URIs.GeoJSONLD)
		mux.Handle(settings.URIs.GeoJSONLD, geojsonld_handler)

		for _, alt_path := range settings.URIs.GeoJSONLDAlt {

			aa_log.Debug(logger, "Handle GeoJSON-LD endpoint at %s\n", alt_path)
			mux.Handle(alt_path, geojsonld_handler)
		}
	}

	if settings.Capabilities.NavPlace {

		aa_log.Debug(logger, "navPlace support enabled")

		navplace_opts := &www.NavPlaceHandlerOptions{
			Reader:      settings.Reader,
			MaxFeatures: settings.NavPlaceMaxFeatures,
			Logger:      logger,
		}

		navplace_handler, err := www.NavPlaceHandler(navplace_opts)

		if err != nil {
			return fmt.Errorf("Failed to create IIIF navPlace handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			navplace_handler = settings.CORSWrapper.Handler(navplace_handler)
		}

		aa_log.Debug(logger, "Handle navPlace endpoint at %s\n", settings.URIs.NavPlace)
		mux.Handle(settings.URIs.NavPlace, navplace_handler)

		for _, alt_path := range settings.URIs.NavPlaceAlt {

			aa_log.Debug(logger, "Handle navPlace endpoint at %s\n", alt_path)
			mux.Handle(alt_path, navplace_handler)
		}
	}

	if settings.Capabilities.Select {

		aa_log.Debug(logger, "Select support enabled")

		select_opts := &www.SelectHandlerOptions{
			Pattern: settings.SelectPattern,
			Reader:  settings.Reader,
			Logger:  logger,
		}

		select_handler, err := www.SelectHandler(select_opts)

		if err != nil {
			return fmt.Errorf("Failed to create Select handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			select_handler = settings.CORSWrapper.Handler(select_handler)
		}

		aa_log.Debug(logger, "Handle Select endpoint at %s\n", settings.URIs.Select)
		mux.Handle(settings.URIs.Select, select_handler)

		for _, alt_path := range settings.URIs.SelectAlt {

			aa_log.Debug(logger, "Handle select endpoint at %s\n", alt_path)
			mux.Handle(alt_path, select_handler)
		}
	}

	if settings.Capabilities.WebFinger {

		aa_log.Debug(logger, "WebFinger support enabled")

		webfinger_opts := &www.WebfingerHandlerOptions{
			Reader:       settings.Reader,
			Logger:       logger,
			URIs:         settings.URIs,
			Capabilities: settings.Capabilities,
			Hostname:     settings.WebFingerHostname,
		}

		webfinger_handler, err := www.WebfingerHandler(webfinger_opts)

		if err != nil {
			return fmt.Errorf("Failed to create WebFinger handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			webfinger_handler = settings.CORSWrapper.Handler(webfinger_handler)
		}

		aa_log.Debug(logger, "Handle WebFinger endpoint at %s\n", settings.URIs.WebFinger)
		mux.Handle(settings.URIs.WebFinger, webfinger_handler)

		for _, alt_path := range settings.URIs.WebFingerAlt {

			aa_log.Debug(logger, "Handle WebFinger endpoint at %s\n", alt_path)
			mux.Handle(alt_path, webfinger_handler)
		}
	}

	if settings.Capabilities.SearchAPI {

		search_opts := api.SearchAPIHandlerOptions{
			Database:      settings.SearchDatabase,
			EnableGeoJSON: true,
			GeoJSONReader: settings.Reader,
		}

		search_api_handler, err := api.SearchAPIHandler(search_opts)

		if err != nil {
			return fmt.Errorf("Failed to create search handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			search_api_handler = settings.CORSWrapper.Handler(search_api_handler)
		}

		mux.Handle(settings.URIs.SearchAPI, search_api_handler)
	}

	if settings.Capabilities.PointInPolygonAPI {
		// To do: Need to sort out what's necessary to create *spatial_app.SpatialApplication
		// https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/app/app.go#L20
	}

	// Common code for HTML handler (public and/or edit handlers)

	var bootstrap_opts *bootstrap.BootstrapOptions
	var maps_opts *maps.MapsOptions

	if settings.HasHTMLCapabilities() {

		bootstrap_opts = bootstrap.DefaultBootstrapOptions()

		err = bootstrap.AppendAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append Bootstrap asset handlers, %w", err)
		}

		err = www.AppendStaticAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %w", err)
		}

		err = settings.CustomChrome.AppendStaticAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append custom asset handlers, %w", err)
		}

		// Final map stuff

		maps_opts = maps.DefaultMapsOptions()

		err = map_www.AppendStaticAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %v")
		}

		err = settings.MapProvider.AppendAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append provider asset handlers, %v", err)
		}

		// Null handler for annoying things like favicons

		null_handler := www.NewNullHandler()

		favicon_path := filepath.Join(settings.URIs.Id, "favicon.ico")
		mux.Handle(favicon_path, null_handler)
	}

	// Public HTML handlers

	// To do: Consider hooks to require auth?

	if settings.Capabilities.Index {

		index_opts := www.IndexHandlerOptions{
			Templates:    t,
			URIs:         settings.URIs,
			Capabilities: settings.Capabilities,
		}

		index_handler, err := www.IndexHandler(index_opts)

		if err != nil {
			return fmt.Errorf("Failed to create Index handler, %w", err)
		}

		index_handler = bootstrap.AppendResourcesHandlerWithPrefix(index_handler, bootstrap_opts, settings.URIs.URIPrefix)
		index_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(index_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		index_handler = settings.CustomChrome.WrapHandler(index_handler)
		index_handler = settings.Authenticator.WrapHandler(index_handler)

		aa_log.Debug(logger, "Handle Index endpoint at %s\n", settings.URIs.Index)
		mux.Handle(settings.URIs.Index, index_handler)
	}

	if settings.Capabilities.Id {

		id_opts := www.IDHandlerOptions{
			Templates:    t,
			URIs:         settings.URIs,
			Capabilities: settings.Capabilities,
			Reader:       settings.Reader,
			Logger:       logger,
			MapProvider:  settings.MapProvider.Scheme(),
		}

		id_handler, err := www.IDHandler(id_opts)

		if err != nil {
			return fmt.Errorf("Failed to create Id handler, %w", err)
		}

		id_handler = bootstrap.AppendResourcesHandlerWithPrefix(id_handler, bootstrap_opts, settings.URIs.URIPrefix)
		id_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(id_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		id_handler = settings.CustomChrome.WrapHandler(id_handler)
		id_handler = settings.Authenticator.WrapHandler(id_handler)

		aa_log.Debug(logger, "Handle Id endpoint at %s\n", settings.URIs.Id)
		mux.Handle(settings.URIs.Id, id_handler)
	}

	if settings.Capabilities.Search {

		search_opts := www.SearchHandlerOptions{
			Templates:    t,
			URIs:         settings.URIs,
			Capabilities: settings.Capabilities,
			Database:     settings.SearchDatabase,
			MapProvider:  settings.MapProvider.Scheme(),
		}

		search_handler, err := www.SearchHandler(search_opts)

		if err != nil {
			return fmt.Errorf("Failed to create Search handler, %w", err)
		}

		search_handler = bootstrap.AppendResourcesHandlerWithPrefix(search_handler, bootstrap_opts, settings.URIs.URIPrefix)
		search_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(search_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		search_handler = settings.CustomChrome.WrapHandler(search_handler)
		search_handler = settings.Authenticator.WrapHandler(search_handler)

		aa_log.Debug(logger, "Handle Search endpoint at %s\n", settings.URIs.Search)
		mux.Handle(settings.URIs.Search, search_handler)
	}

	// Edit/write HTML handlers

	if settings.Capabilities.EditUI {

		aa_log.Debug(logger, "Edit user interface support enabled")

		// Edit geometry

		geom_t := t.Lookup("geometry")

		if geom_t == nil {
			return fmt.Errorf("Failed to load 'geometry' template")
		}

		geom_opts := &www.EditGeometryHandlerOptions{
			Authenticator: settings.Authenticator,
			MapProvider:   settings.MapProvider.Scheme(),
			URIs:          settings.URIs,
			Capabilities:  settings.Capabilities,
			Template:      geom_t,
			Logger:        logger,
			Reader:        settings.Reader,
		}

		geom_handler, err := www.EditGeometryHandler(geom_opts)

		if err != nil {
			return fmt.Errorf("Failed to create edit geometry handler, %w", err)
		}

		// To do: Edit properties

		// Create feature

		create_t := t.Lookup("create")

		if create_t == nil {
			return fmt.Errorf("Failed to load 'create' template")
		}

		create_opts := &www.CreateFeatureHandlerOptions{
			Authenticator:    settings.Authenticator,
			MapProvider:      settings.MapProvider.Scheme(),
			URIs:             settings.URIs,
			Capabilities:     settings.Capabilities,
			Template:         create_t,
			Logger:           logger,
			Reader:           settings.Reader,
			CustomProperties: settings.CustomEditProperties,
		}

		create_handler, err := www.CreateFeatureHandler(create_opts)

		if err != nil {
			return fmt.Errorf("Failed to create create feature handler, %w", err)
		}

		geom_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(geom_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		geom_handler = bootstrap.AppendResourcesHandlerWithPrefix(geom_handler, bootstrap_opts, settings.URIs.URIPrefix)
		geom_handler = settings.CustomChrome.WrapHandler(geom_handler)
		geom_handler = settings.Authenticator.WrapHandler(geom_handler)

		aa_log.Debug(logger, "Handle edit geometry endpoint at %s\n", path_edit_geometry)
		mux.Handle(settings.URIs.EditGeometry, geom_handler)

		// START OF wasm stuff

		err = wasm_exec.AppendAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append wasm asset handlers, %w", err)
		}

		err = wasm_placetypes.AppendAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append wasm placetypes asset handlers, %w", err)
		}

		err = wasm_validate.AppendAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append wasm validate asset handlers, %w", err)
		}

		wasm_exec_opts := wasm_exec.DefaultWASMOptions()

		// START OF I don't like having to do this but since the default 'whosonfirst.validate.feature.js'
		// package (in go-whosonfirst-validate-wasm) has a relative path and, importantly, no well-defined
		// way to specify the wasm path (yet) this is what we're going to do in conjunction with writing
		// our own 'whosonfirst.browser.validate' package that fetches the custom URI from 'whosonfirst.browser.uris'.
		// I suppose it would be easy enough to add a 'setWasmURI' or a 'setWasmPrefix' method to 'whosonfirst.validate.feature.js'
		// but today that hasn't happened.

		// This file is served by the http/www/static.go handlers
		wasm_validate_uri := "/wasm/validate_feature.wasm"
		wasm_placetypes_uri := "/wasm/whosonfirst_placetypes.wasm"

		if settings.URIs.URIPrefix != "" {

			validate_uri, err := url.JoinPath(settings.URIs.URIPrefix, wasm_validate_uri)

			if err != nil {
				return fmt.Errorf("Failed to assign URI prefix to validate wasm path, %w", err)
			}

			placetypes_uri, err := url.JoinPath(settings.URIs.URIPrefix, wasm_placetypes_uri)

			if err != nil {
				return fmt.Errorf("Failed to assign URI prefix to placetypes wasm path, %w", err)
			}

			wasm_validate_uri = validate_uri
			wasm_placetypes_uri = placetypes_uri
		}

		settings.URIs.AddCustomURI("validate_wasm", wasm_validate_uri)
		settings.URIs.AddCustomURI("placetypes_wasm", wasm_placetypes_uri)

		// END OF I don't like having	to do this

		create_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(create_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		create_handler = wasm_exec.AppendResourcesHandlerWithPrefix(create_handler, wasm_exec_opts, settings.URIs.URIPrefix)

		create_handler = appendCustomMiddlewareHandlers(settings, settings.URIs.CreateFeature, create_handler)

		create_handler = bootstrap.AppendResourcesHandlerWithPrefix(create_handler, bootstrap_opts, settings.URIs.URIPrefix)
		create_handler = settings.CustomChrome.WrapHandler(create_handler)
		create_handler = settings.Authenticator.WrapHandler(create_handler)

		aa_log.Debug(logger, "Handle create feature endpoint at %s\n", path_create_feature)
		mux.Handle(settings.URIs.CreateFeature, create_handler)

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

			if settings.URIs.URIPrefix != "" {

				uri, err := url.JoinPath(settings.URIs.URIPrefix, custom_validate_uri)

				if err != nil {
					return fmt.Errorf("Failed to assign URI prefix to custom validate wasm path, %w", err)
				}

				custom_validate_uri = uri
			}

			settings.URIs.AddCustomURI("custom_validate_wasm", custom_validate_uri)

			aa_log.Debug(logger, "Handle custom validation WASM endpoint at %s\n", custom_validate_uri)
			mux.Handle(custom_validate_uri, validation_handler)
		}
	}

	if settings.Capabilities.EditAPI {

		aa_log.Debug(logger, "Edit API support enabled")

		// Writers are created at runtime using the http/api/publish.go#publishFeature
		// method which in turn calls writer/writer.go#NewWriter

		// Deprecate a record

		deprecate_opts := &api.DeprecateFeatureHandlerOptions{
			Reader:        settings.Reader,
			Cache:         settings.Cache,
			Logger:        logger,
			Authenticator: settings.Authenticator,
			Exporter:      settings.Exporter,
			WriterURIs:    settings.WriterURIs,
		}

		deprecate_handler, err := api.DeprecateFeatureHandler(deprecate_opts)

		if err != nil {
			return fmt.Errorf("Failed to create deprecate feature handler, %w", err)
		}

		deprecate_handler = settings.Authenticator.WrapHandler(deprecate_handler)

		aa_log.Debug(logger, "Handle deprecate feature API endpoint at %s\n", settings.URIs.DeprecateFeatureAPI)
		mux.Handle(settings.URIs.DeprecateFeatureAPI, deprecate_handler)

		// Mark a record as ceased

		cessate_opts := &api.CessateFeatureHandlerOptions{
			Reader:        settings.Reader,
			Cache:         settings.Cache,
			Logger:        logger,
			Authenticator: settings.Authenticator,
			Exporter:      settings.Exporter,
			WriterURIs:    settings.WriterURIs,
		}

		cessate_handler, err := api.CessateFeatureHandler(cessate_opts)

		if err != nil {
			return fmt.Errorf("Failed to create cessate feature handler, %w", err)
		}

		cessate_handler = settings.Authenticator.WrapHandler(cessate_handler)

		aa_log.Debug(logger, "Handle cessate feature API endpoint at %s\n", settings.URIs.CessateFeatureAPI)
		mux.Handle(settings.URIs.CessateFeatureAPI, cessate_handler)

		// Edit geometry

		geom_opts := &api.UpdateGeometryHandlerOptions{
			Reader:                settings.Reader,
			Cache:                 settings.Cache,
			Logger:                logger,
			Authenticator:         settings.Authenticator,
			Exporter:              settings.Exporter,
			WriterURIs:            settings.WriterURIs,
			PointInPolygonService: settings.PointInPolygonService,
		}

		geom_handler, err := api.UpdateGeometryHandler(geom_opts)

		if err != nil {
			return fmt.Errorf("Failed to create update geometry handler, %w", err)
		}

		geom_handler = settings.Authenticator.WrapHandler(geom_handler)

		aa_log.Debug(logger, "Handle edit geometry API endpoint at %s\n", settings.URIs.EditGeometryAPI)
		mux.Handle(settings.URIs.EditGeometryAPI, geom_handler)

		// Create a new feature

		create_opts := &api.CreateFeatureHandlerOptions{
			Reader:                settings.Reader,
			Cache:                 settings.Cache,
			Logger:                logger,
			Authenticator:         settings.Authenticator,
			Exporter:              settings.Exporter,
			WriterURIs:            settings.WriterURIs,
			PointInPolygonService: settings.PointInPolygonService,
			CustomProperties:      settings.CustomEditProperties,
			CustomValidationFunc:  settings.CustomEditValidationFunc,
		}

		create_handler, err := api.CreateFeatureHandler(create_opts)

		if err != nil {
			return fmt.Errorf("Failed to create create feature handler, %w", err)
		}

		create_handler = settings.Authenticator.WrapHandler(create_handler)

		aa_log.Debug(logger, "Handle create feature API endpoint at %s\n", settings.URIs.CreateFeatureAPI)
		mux.Handle(settings.URIs.CreateFeatureAPI, create_handler)
	}

	// START OF customizable handlers
	// Note that middleware handlers are applied inline above

	// Custom handlers

	for path, h := range settings.CustomWWWHandlers {

		h = maps.AppendResourcesHandlerWithPrefixAndProvider(h, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		h = bootstrap.AppendResourcesHandlerWithPrefix(h, bootstrap_opts, settings.URIs.URIPrefix)
		h = settings.CustomChrome.WrapHandler(h)

		h = settings.Authenticator.WrapHandler(h)

		aa_log.Debug(logger, "Handle custom www endpoint at %s\n", path)
		mux.Handle(path, h)
	}

	for path, h := range settings.CustomAPIHandlers {

		h = settings.Authenticator.WrapHandler(h)

		aa_log.Debug(logger, "Handle custom API endpoint at %s\n", path)
		mux.Handle(path, h)
	}

	// Custom asset handlers

	if len(settings.CustomAssetHandlerFunctions) > 0 {

		aa_log.Debug(logger, "Register custom asset handlers")

		for idx, handler_func := range settings.CustomAssetHandlerFunctions {

			err := handler_func(mux, settings.URIs.URIPrefix)

			if err != nil {
				return fmt.Errorf("Failed to register custom asset handler function at offset %d, %w", idx, err)
			}
		}
	}

	// END OF customizable handlers

	// START OF uris.js
	// Publish settings.URIs as a JS file so that other JS knows where to find things

	uris_t := js_t.Lookup("uris")

	if uris_t == nil {
		return fmt.Errorf("Failed to load 'uris' javascript template")
	}

	uris_opts := &www.URIsHandlerOptions{
		URIs:     settings.URIs,
		Template: uris_t,
	}

	uris_handler, err := www.URIsHandler(uris_opts)

	if err != nil {
		return fmt.Errorf("Failed to create URIs handler, %w", err)
	}

	uris_path := "/javascript/whosonfirst.browser.uris.js"

	if settings.URIs.URIPrefix != "" {

		path, err := url.JoinPath(settings.URIs.URIPrefix, uris_path)

		if err != nil {
			return fmt.Errorf("Failed to assign URI prefix to %s, %w", uris_path, err)
		}

		uris_path = path
	}

	aa_log.Debug(logger, "Handle whosonfirst.browser.uris.js endpoint at %s\n", uris_path)
	mux.Handle(uris_path, uris_handler)

	// END OF uris.js

	aa_log.Info(logger, "Time to set up: %v\n", time.Since(t1))

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
