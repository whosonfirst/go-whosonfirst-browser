package browser

import (
	// _ "github.com/aaronland/go-http-server-tsnet"
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

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	map_www "github.com/aaronland/go-http-maps/http/www"
	"github.com/aaronland/go-http-ping/v2"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-template/html"
	"github.com/sfomuseum/go-template/text"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/api"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/templates/javascript"
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

	if settings.Verbose {
		logger.Printf("Handle ping endpoint at %s\n", settings.URIs.Ping)
	}

	// START OF uris.js

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

	mux.Handle(uris_path, uris_handler)

	if settings.Verbose {
		logger.Printf("Handle whosonfirst.browser.uris.js endpoint at %s\n", uris_path)
	}

	// END OF uris.js

	if settings.Capabilities.PNG {

		if settings.Verbose {
			logger.Println("png support enabled")
		}

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

		mux.Handle(settings.URIs.PNG, png_handler)

		if settings.Verbose {
			logger.Printf("Handle png endpoint at %s\n", settings.URIs.PNG)
		}

		for _, alt_path := range settings.URIs.PNGAlt {

			mux.Handle(alt_path, png_handler)

			if settings.Verbose {
				logger.Printf("Handle png endpoint at %s\n", alt_path)
			}
		}
	}

	if settings.Capabilities.SVG {

		if settings.Verbose {
			logger.Println("svg support enabled")
		}

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

		mux.Handle(settings.URIs.SVG, svg_handler)

		if settings.Verbose {
			log.Printf("handle svg endpoint at %s\n", settings.URIs.SVG)
		}

		for _, alt_path := range settings.URIs.SVGAlt {

			mux.Handle(alt_path, svg_handler)

			if settings.Verbose {
				log.Printf("handle svg endpoint at %s\n", alt_path)
			}
		}
	}

	if settings.Capabilities.SPR {

		if settings.Verbose {
			log.Println("spr support enabled")
		}

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

		mux.Handle(settings.URIs.SPR, spr_handler)

		if settings.Verbose {
			log.Printf("handle spr endpoint at %s\n", settings.URIs.SPR)
		}

		for _, alt_path := range settings.URIs.SPRAlt {

			mux.Handle(alt_path, spr_handler)

			if settings.Verbose {
				log.Printf("handle spr endpoint at %s\n", alt_path)
			}
		}
	}

	if settings.Capabilities.GeoJSON {

		if settings.Verbose {
			log.Println("geojson support enabled")
		}

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

		mux.Handle(settings.URIs.GeoJSON, geojson_handler)

		if settings.Verbose {
			logger.Printf("Handle geojson endpoint at %s\n", settings.URIs.GeoJSON)
		}

		for _, alt_path := range settings.URIs.GeoJSONAlt {

			mux.Handle(alt_path, geojson_handler)

			if settings.Verbose {
				logger.Printf("Handle geojson endpoint at %s\n", alt_path)
			}
		}
	}

	if settings.Capabilities.GeoJSONLD {

		if settings.Verbose {
			log.Println("geojsonld support enabled")
		}

		geojsonld_opts := &www.GeoJSONLDHandlerOptions{
			Reader: settings.Reader,
			Logger: logger,
		}

		geojsonld_handler, err := www.GeoJSONLDHandler(geojsonld_opts)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON LD handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			geojsonld_handler = settings.CORSWrapper.Handler(geojsonld_handler)
		}

		mux.Handle(settings.URIs.GeoJSONLD, geojsonld_handler)

		if settings.Verbose {
			logger.Printf("Handle geojsonld endpoint at %s\n", settings.URIs.GeoJSONLD)
		}

		for _, alt_path := range settings.URIs.GeoJSONLDAlt {

			mux.Handle(alt_path, geojsonld_handler)

			if settings.Verbose {
				logger.Printf("Handle geojsonld endpoint at %s\n", alt_path)
			}
		}
	}

	if settings.Capabilities.NavPlace {

		if settings.Verbose {
			log.Println("navplace support enabled")
		}

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

		mux.Handle(settings.URIs.NavPlace, navplace_handler)

		if settings.Verbose {
			logger.Printf("Handle navplace endpoint at %s\n", settings.URIs.NavPlace)
		}

		for _, alt_path := range settings.URIs.NavPlaceAlt {
			mux.Handle(alt_path, navplace_handler)

			if settings.Verbose {
				logger.Printf("Handle navplace endpoint at %s\n", alt_path)
			}
		}
	}

	if settings.Capabilities.Select {

		if settings.Verbose {
			log.Println("select support enabled")
		}

		select_opts := &www.SelectHandlerOptions{
			Pattern: settings.SelectPattern,
			Reader:  settings.Reader,
			Logger:  logger,
		}

		select_handler, err := www.SelectHandler(select_opts)

		if err != nil {
			return fmt.Errorf("Failed to create select handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			select_handler = settings.CORSWrapper.Handler(select_handler)
		}

		mux.Handle(settings.URIs.Select, select_handler)

		if settings.Verbose {
			log.Printf("handle select endpoint at %s\n", settings.URIs.Select)
		}

		for _, alt_path := range settings.URIs.SelectAlt {

			mux.Handle(alt_path, select_handler)

			if settings.Verbose {
				log.Printf("handle select endpoint at %s\n", alt_path)
			}
		}
	}

	if settings.Capabilities.WebFinger {

		if settings.Verbose {
			log.Println("webfinger support enabled")
		}

		webfinger_opts := &www.WebfingerHandlerOptions{
			Reader:       settings.Reader,
			Logger:       logger,
			URIs:         settings.URIs,
			Capabilities: settings.Capabilities,
			Hostname:     webfinger_hostname, // UPDATE M
		}

		webfinger_handler, err := www.WebfingerHandler(webfinger_opts)

		if err != nil {
			return fmt.Errorf("Failed to create webfinger handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			webfinger_handler = settings.CORSWrapper.Handler(webfinger_handler)
		}

		mux.Handle(settings.URIs.WebFinger, webfinger_handler)

		if settings.Verbose {
			log.Printf("handle webfinger endpoint at %s\n", settings.URIs.WebFinger)
		}

		for _, alt_path := range settings.URIs.WebFingerAlt {
			mux.Handle(alt_path, webfinger_handler)

			if settings.Verbose {
				log.Printf("handle webfinger endpoint at %s\n", alt_path)
			}
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
			return fmt.Errorf("Failed to create index handler, %w", err)
		}

		index_handler = bootstrap.AppendResourcesHandlerWithPrefix(index_handler, bootstrap_opts, settings.URIs.URIPrefix)
		index_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(index_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		index_handler = settings.CustomChrome.WrapHandler(index_handler)
		index_handler = settings.Authenticator.WrapHandler(index_handler)

		mux.Handle(settings.URIs.Index, index_handler)

		if settings.Verbose {
			log.Printf("handle index endpoint at %s\n", settings.URIs.Index)
		}
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
			return fmt.Errorf("Failed to create ID handler, %w", err)
		}

		id_handler = bootstrap.AppendResourcesHandlerWithPrefix(id_handler, bootstrap_opts, settings.URIs.URIPrefix)
		id_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(id_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		id_handler = settings.CustomChrome.WrapHandler(id_handler)
		id_handler = settings.Authenticator.WrapHandler(id_handler)

		mux.Handle(settings.URIs.Id, id_handler)

		if settings.Verbose {
			log.Printf("handle ID endpoint at %s\n", settings.URIs.Id)
		}

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
			return fmt.Errorf("Failed to create search handler, %w", err)
		}

		search_handler = bootstrap.AppendResourcesHandlerWithPrefix(search_handler, bootstrap_opts, settings.URIs.URIPrefix)
		search_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(search_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		search_handler = settings.CustomChrome.WrapHandler(search_handler)
		search_handler = settings.Authenticator.WrapHandler(search_handler)

		mux.Handle(settings.URIs.Search, search_handler)

		if settings.Verbose {
			log.Printf("handle search endpoint at %s\n", settings.URIs.Search)
		}
	}

	// Edit/write HTML handlers

	if settings.Capabilities.EditUI {

		if settings.Verbose {
			log.Println("edit user interface support enabled")
		}

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
			Authenticator: settings.Authenticator,
			MapProvider:   settings.MapProvider.Scheme(),
			URIs:          settings.URIs,
			Capabilities:  settings.Capabilities,
			Template:      create_t,
			Logger:        logger,
			Reader:        settings.Reader,
		}

		create_handler, err := www.CreateFeatureHandler(create_opts)

		if err != nil {
			return fmt.Errorf("Failed to create create feature handler, %w", err)
		}

		geom_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(geom_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		geom_handler = bootstrap.AppendResourcesHandlerWithPrefix(geom_handler, bootstrap_opts, settings.URIs.URIPrefix)
		geom_handler = settings.CustomChrome.WrapHandler(geom_handler)
		geom_handler = settings.Authenticator.WrapHandler(geom_handler)

		mux.Handle(settings.URIs.EditGeometry, geom_handler)

		if settings.Verbose {
			log.Printf("handle edit geometry endpoint at %s\n", path_edit_geometry)
		}

		err = wasm_validate.AppendAssetHandlersWithPrefix(mux, settings.URIs.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append wasm validate asset handlers, %w", err)
		}

		wasm_validate_opts := wasm_validate.DefaultWASMOptions()
		wasm_validate_opts.EnableWASMExec()
		
		create_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(create_handler, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		create_handler = wasm_validate.AppendResourcesHandlerWithPrefix(create_handler, wasm_validate_opts, settings.URIs.URIPrefix)
		create_handler = bootstrap.AppendResourcesHandlerWithPrefix(create_handler, bootstrap_opts, settings.URIs.URIPrefix)
		create_handler = settings.CustomChrome.WrapHandler(create_handler)
		create_handler = settings.Authenticator.WrapHandler(create_handler)

		mux.Handle(settings.URIs.CreateFeature, create_handler)

		if settings.Verbose {
			log.Printf("handle create feature endpoint at %s\n", path_create_feature)
		}
	}

	if settings.Capabilities.EditAPI {

		if settings.Verbose {
			log.Println("Edit api support enabled")
		}
		
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
		mux.Handle(settings.URIs.DeprecateFeatureAPI, deprecate_handler)

		if settings.Verbose {
			log.Printf("handle deprecate feature endpoint at %s\n", settings.URIs.DeprecateFeatureAPI)
		}

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
		mux.Handle(settings.URIs.CessateFeatureAPI, cessate_handler)

		if settings.Verbose {
			log.Printf("handle cessate feature endpoint at %s\n", settings.URIs.CessateFeatureAPI)
		}

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
			return fmt.Errorf("Failed to create uupdate geometry handler, %w", err)
		}

		geom_handler = settings.Authenticator.WrapHandler(geom_handler)
		mux.Handle(settings.URIs.EditGeometryAPI, geom_handler)

		if settings.Verbose {
			log.Printf("handle edit geometry endpoint at %s\n", settings.URIs.EditGeometryAPI)
		}

		// Create a new feature

		create_opts := &api.CreateFeatureHandlerOptions{
			Reader:                settings.Reader,
			Cache:                 settings.Cache,
			Logger:                logger,
			Authenticator:         settings.Authenticator,
			Exporter:              settings.Exporter,
			WriterURIs:            settings.WriterURIs,
			PointInPolygonService: settings.PointInPolygonService,
		}

		create_handler, err := api.CreateFeatureHandler(create_opts)

		if err != nil {
			return fmt.Errorf("Failed to create create feature handler, %w", err)
		}

		create_handler = settings.Authenticator.WrapHandler(create_handler)
		mux.Handle(settings.URIs.CreateFeatureAPI, create_handler)

		if settings.Verbose {
			log.Printf("handle create feature endpoint at %s\n", settings.URIs.CreateFeatureAPI)
		}

	}

	// Custom handlers

	for path, h := range settings.CustomWWWHandlers {

		h = maps.AppendResourcesHandlerWithPrefixAndProvider(h, settings.MapProvider, maps_opts, settings.URIs.URIPrefix)
		h = bootstrap.AppendResourcesHandlerWithPrefix(h, bootstrap_opts, settings.URIs.URIPrefix)
		h = settings.CustomChrome.WrapHandler(h)

		h = settings.Authenticator.WrapHandler(h)

		mux.Handle(path, h)

		if settings.Verbose {
			log.Printf("handle custom www endpoint at %s\n", path)
		}
	}

	for path, h := range settings.CustomAPIHandlers {

		h = settings.Authenticator.WrapHandler(h)

		mux.Handle(path, h)

		if settings.Verbose {
			log.Printf("handle custom api endpoint at %s\n", path)
		}
	}

	// Finally, start the server

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new search for '%s', %w", server_uri, err)
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}
