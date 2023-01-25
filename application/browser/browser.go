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
	"path/filepath"
	"regexp"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	map_www "github.com/aaronland/go-http-maps/http/www"
	"github.com/aaronland/go-http-ping/v2"
	"github.com/aaronland/go-http-server"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/api"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/templates/html"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
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

	// To do: Convert Config struct in to a Settings struct and
	// then call RunWithSettings (below)

	if cfg.EnableAll {
		cfg.EnableGraphics = true
		cfg.EnableData = true
		cfg.EnableHTML = true
	}

	if cfg.EnableSearch {
		enable_search_api = true
		enable_search_api_geojson = true
		enable_search_html = true
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

	if enable_search_html {
		cfg.EnableHTML = true
	}

	if cfg.EnableHTML {
		cfg.EnableGeoJSON = true
		cfg.EnablePNG = true
	}

	if cfg.EnableEdit {
		cfg.EnableEditAPI = true
		cfg.EnableEditUI = true
	}

	if cfg.EnableEditUI {
		cfg.EnableEditAPI = true
	}

	// To do: pre-fill defaults

	settings, err := SettingsFromConfig(ctx, cfg, logger)

	if err != nil {
		return fmt.Errorf("Failed to create settings from config, %w", err)
	}

	return RunWithSettings(ctx, settings, logger)
}

func RunWithSettings(ctx context.Context, settings *Settings, logger *log.Logger) error {

	// Set up templates
	// To do: Once we have config/settings stuff working this needs to be able to
	// specify a custom fs.FS for reading templates from

	t, err := html.LoadTemplates(ctx)

	if err != nil {
		return fmt.Errorf("Failed to load templates, %w", err)
	}

	// Start setting up handlers

	mux := http.NewServeMux()

	ping_handler, err := ping.PingPongHandler()

	if err != nil {
		return fmt.Errorf("Failed to create ping handler, %w", err)
	}

	mux.Handle(path_ping, ping_handler)

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

		mux.Handle(path_png, png_handler)

		if settings.Verbose {
			logger.Printf("Handle png endpoint at %s\n", path_png)
		}

		for _, alt_path := range path_png_alt {
			mux.Handle(alt_path, png_handler)
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

		mux.Handle(path_svg, svg_handler)

		if settings.Verbose {
			log.Printf("handle svg endpoint at %s\n", path_svg)
		}

		for _, alt_path := range path_svg_alt {
			mux.Handle(alt_path, svg_handler)
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

		mux.Handle(path_spr, spr_handler)

		if settings.Verbose {
			log.Printf("handle spr endpoint at %s\n", path_spr)
		}

		for _, alt_path := range path_spr_alt {
			mux.Handle(alt_path, spr_handler)
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

		mux.Handle(path_geojson, geojson_handler)

		if settings.Verbose {
			logger.Printf("Handle geojson endpoint at %s\n", path_geojson)
		}

		for _, alt_path := range path_geojson_alt {
			mux.Handle(alt_path, geojson_handler)
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

		mux.Handle(path_geojsonld, geojsonld_handler)

		if settings.Verbose {
			logger.Printf("Handle geojsonld endpoint at %s\n", path_geojsonld)
		}

		for _, alt_path := range path_geojsonld_alt {
			mux.Handle(alt_path, geojsonld_handler)
		}
	}

	if settings.Capabilities.NavPlace {

		if settings.Verbose {
			log.Println("navplace support enabled")
		}

		navplace_opts := &www.NavPlaceHandlerOptions{
			Reader:      settings.Reader,
			MaxFeatures: navplace_max_features, // UPDATE ME
			Logger:      logger,
		}

		navplace_handler, err := www.NavPlaceHandler(navplace_opts)

		if err != nil {
			return fmt.Errorf("Failed to create IIIF navPlace handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			navplace_handler = settings.CORSWrapper.Handler(navplace_handler)
		}

		mux.Handle(path_navplace, navplace_handler)

		if settings.Verbose {
			logger.Printf("Handle navplace endpoint at %s\n", path_navplace)
		}

		for _, alt_path := range path_navplace_alt {
			mux.Handle(alt_path, navplace_handler)
		}
	}

	if settings.Capabilities.Select {

		if settings.Verbose {
			log.Println("select support enabled")
		}

		// UPDATE ME

		if select_pattern == "" {
			return fmt.Errorf("Missing -select-pattern parameter.")
		}

		pat, err := regexp.Compile(select_pattern)

		if err != nil {
			return fmt.Errorf("Failed to compile select pattern (%s), %w", select_pattern, err)
		}

		select_opts := &www.SelectHandlerOptions{
			Pattern: pat,
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

		mux.Handle(path_select, select_handler)

		if settings.Verbose {
			log.Printf("handle select endpoint at %s\n", path_select)
		}

		for _, alt_path := range path_select_alt {
			mux.Handle(alt_path, select_handler)
		}
	}

	if settings.Capabilities.WebFinger {

		if settings.Verbose {
			log.Println("webfinger support enabled")
		}

		webfinger_opts := &www.WebfingerHandlerOptions{
			Reader:       settings.Reader,
			Logger:       logger,
			Paths:        settings.Paths,
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

		mux.Handle(path_webfinger, webfinger_handler)

		if settings.Verbose {
			log.Printf("handle webfinger endpoint at %s\n", path_webfinger)
		}

		for _, alt_path := range path_webfinger_alt {
			mux.Handle(alt_path, webfinger_handler)
		}
	}

	// START OF probably due a rethink shortly

	if enable_search_api {

		if search_database_uri == "" {
			return fmt.Errorf("-enable-search-api flag is true but -search-database-uri flag is empty.")
		}

		search_db, err := fulltext.NewFullTextDatabase(ctx, search_database_uri)

		if err != nil {
			return fmt.Errorf("Failed to create fulltext database for '%s', %w", search_database_uri, err)
		}

		search_opts := www.SearchAPIHandlerOptions{
			Database: search_db,
		}

		if enable_search_api_geojson {
			search_opts.EnableGeoJSON = true
			search_opts.GeoJSONReader = settings.Reader
		}

		search_handler, err := www.SearchAPIHandler(search_opts)

		if err != nil {
			return fmt.Errorf("Failed to create search handler, %w", err)
		}

		if settings.CORSWrapper != nil {
			search_handler = settings.CORSWrapper.Handler(search_handler)
		}

		mux.Handle(path_search_api, search_handler)
	}

	// END OF probably due a rethink shortly

	// Common code for HTML handler (public and/or edit handlers)

	var bootstrap_opts *bootstrap.BootstrapOptions
	var maps_opts *maps.MapsOptions

	if settings.Capabilities.HTML || settings.Capabilities.EditUI {

		bootstrap_opts = bootstrap.DefaultBootstrapOptions()

		err = bootstrap.AppendAssetHandlersWithPrefix(mux, settings.Paths.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append Bootstrap asset handlers, %w", err)
		}

		err = www.AppendStaticAssetHandlersWithPrefix(mux, settings.Paths.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %w", err)
		}

		err = settings.CustomChrome.AppendStaticAssetHandlersWithPrefix(mux, settings.Paths.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append custom asset handlers, %w", err)
		}

		// Final map stuff

		maps_opts = maps.DefaultMapsOptions()

		err = map_www.AppendStaticAssetHandlersWithPrefix(mux, settings.Paths.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %v")
		}

		err = settings.MapProvider.AppendAssetHandlersWithPrefix(mux, settings.Paths.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append provider asset handlers, %v", err)
		}

		// Null handler for annoying things like favicons

		null_handler := www.NewNullHandler()

		favicon_path := filepath.Join(path_id, "favicon.ico")
		mux.Handle(favicon_path, null_handler)
	}

	// Public HTML handlers

	// To do: Consider hooks to require auth?

	if settings.Capabilities.HTML {

		// Note that we append all the handler to mux at the end of this block so that they
		// can be updated with map-related middleware where necessary

		var index_handler http.Handler
		var id_handler http.Handler
		var search_handler http.Handler

		if enable_index {

			index_opts := www.IndexHandlerOptions{
				Templates:    t,
				Paths:        settings.Paths,
				Capabilities: settings.Capabilities,
			}

			index_h, err := www.IndexHandler(index_opts)

			if err != nil {
				return fmt.Errorf("Failed to create index handler, %w", err)
			}

			index_handler = index_h
			index_handler = bootstrap.AppendResourcesHandlerWithPrefix(index_handler, bootstrap_opts, settings.Paths.URIPrefix)
		}

		id_opts := www.IDHandlerOptions{
			Templates:    t,
			Paths:        settings.Paths,
			Capabilities: settings.Capabilities,
			Reader:       settings.Reader,
			Logger:       logger,
			MapProvider:  settings.MapProvider.Scheme(),
		}

		id_h, err := www.IDHandler(id_opts)

		if err != nil {
			return fmt.Errorf("Failed to create ID handler, %w", err)
		}

		id_handler = id_h
		id_handler = bootstrap.AppendResourcesHandlerWithPrefix(id_handler, bootstrap_opts, settings.Paths.URIPrefix)

		if enable_search_html {

			search_db, err := fulltext.NewFullTextDatabase(ctx, search_database_uri)

			if err != nil {
				return fmt.Errorf("Failed to create fulltext database for '%s', %w", search_database_uri, err)
			}

			search_opts := www.SearchHandlerOptions{
				Templates:    t,
				Paths:        settings.Paths,
				Capabilities: settings.Capabilities,
				Database:     search_db,
				MapProvider:  settings.MapProvider.Scheme(),
			}

			search_h, err := www.SearchHandler(search_opts)

			if err != nil {
				return fmt.Errorf("Failed to create search handler, %w", err)
			}

			search_handler = search_h
			search_handler = bootstrap.AppendResourcesHandlerWithPrefix(search_handler, bootstrap_opts, settings.Paths.URIPrefix)
		}

		id_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(id_handler, settings.MapProvider, maps_opts, settings.Paths.URIPrefix)
		id_handler = settings.CustomChrome.WrapHandler(id_handler)
		id_handler = settings.Authenticator.WrapHandler(id_handler)

		mux.Handle(path_id, id_handler)

		if settings.Verbose {
			log.Printf("handle ID endpoint at %s\n", path_id)
		}

		if enable_search_html {
			search_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(search_handler, settings.MapProvider, maps_opts, settings.Paths.URIPrefix)
			search_handler = settings.Authenticator.WrapHandler(search_handler)
			mux.Handle(path_search_html, search_handler)
		}

		index_handler = settings.Authenticator.WrapHandler(index_handler)
		mux.Handle("/", index_handler)
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
			Paths:         settings.Paths,
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
			Paths:         settings.Paths,
			Capabilities:  settings.Capabilities,
			Template:      create_t,
			Logger:        logger,
			Reader:        settings.Reader,
		}

		create_handler, err := www.CreateFeatureHandler(create_opts)

		if err != nil {
			return fmt.Errorf("Failed to create create feature handler, %w", err)
		}

		geom_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(geom_handler, settings.MapProvider, maps_opts, settings.Paths.URIPrefix)
		geom_handler = bootstrap.AppendResourcesHandlerWithPrefix(geom_handler, bootstrap_opts, settings.Paths.URIPrefix)
		geom_handler = settings.CustomChrome.WrapHandler(geom_handler)
		geom_handler = settings.Authenticator.WrapHandler(geom_handler)

		mux.Handle(path_edit_geometry, geom_handler)

		if settings.Verbose {
			log.Printf("handle edit geometry endpoint at %s\n", path_edit_geometry)
		}

		create_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(create_handler, settings.MapProvider, maps_opts, settings.Paths.URIPrefix)
		create_handler = bootstrap.AppendResourcesHandlerWithPrefix(create_handler, bootstrap_opts, settings.Paths.URIPrefix)
		create_handler = settings.CustomChrome.WrapHandler(create_handler)
		create_handler = settings.Authenticator.WrapHandler(create_handler)

		mux.Handle(path_create_feature, create_handler)

		if settings.Verbose {
			log.Printf("handle create feature endpoint at %s\n", path_create_feature)
		}
	}

	if settings.Capabilities.EditAPI {

		if settings.Verbose {
			log.Println("edit api support enabled")
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
		mux.Handle(path_api_deprecate_feature, deprecate_handler)

		if settings.Verbose {
			log.Printf("handle deprecate feature endpoint at %s\n", path_api_deprecate_feature)
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
		mux.Handle(path_api_cessate_feature, cessate_handler)

		if settings.Verbose {
			log.Printf("handle cessate feature endpoint at %s\n", path_api_cessate_feature)
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
		mux.Handle(path_api_edit_geometry, geom_handler)

		if settings.Verbose {
			log.Printf("handle edit geometry endpoint at %s\n", path_api_edit_geometry)
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
		mux.Handle(path_api_create_feature, create_handler)

		if settings.Verbose {
			log.Printf("handle create geometry endpoint at %s\n", path_api_edit_geometry)
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
