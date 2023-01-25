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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	map_www "github.com/aaronland/go-http-maps/http/www"
	"github.com/aaronland/go-http-maps/provider"
	"github.com/aaronland/go-http-ping/v2"
	"github.com/aaronland/go-http-server"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-cachereader"
	github_reader "github.com/whosonfirst/go-reader-github"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/chrome"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/api"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/templates/html"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	github_writer "github.com/whosonfirst/go-writer-github/v3"
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
		cfg.EnableGrapics = true
		cfg.EnableData = true
		cfg.EnableHTML = true
	}

	if cfg.EnableSearch {
		enable_search_api = true
		enable_search_api_geojson = true
		enable_search_html = true
	}

	if cfg.EnableGrapics {
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

	settings, err := SettingsFromConfig(ctx, cfg, logger)

	if err != nil {
		return fmt.Errorf("Failed to create settings from config, %w", err)
	}

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

	if cfg.EnablePNG {

		if cfg.Verbose {
			logger.Println("png support enabled")
		}

		sizes := www.DefaultRasterSizes()

		png_opts := &www.RasterHandlerOptions{
			Sizes:  sizes,
			Format: "png",
			Reader: cr,
			Logger: logger,
		}

		png_handler, err := www.RasterHandler(png_opts)

		if err != nil {
			return fmt.Errorf("Failed to create raster/png handler, %w", err)
		}

		mux.Handle(path_png, png_handler)

		if cfg.Verbose {
			logger.Printf("Handle png endpoint at %s\n", path_png)
		}

		for _, alt_path := range path_png_alt {
			mux.Handle(alt_path, png_handler)
		}
	}

	if cfg.EnableSVG {

		if cfg.Verbose {
			logger.Println("svg support enabled")
		}

		sizes := www.DefaultSVGSizes()

		svg_opts := &www.SVGHandlerOptions{
			Sizes:  sizes,
			Reader: cr,
			Logger: logger,
		}

		svg_handler, err := www.SVGHandler(svg_opts)

		if err != nil {
			return fmt.Errorf("Failed to create SVG handler, %w", err)
		}

		if cfg.EnableCORS {
			svg_handler = cors_wrapper.Handler(svg_handler)
		}

		mux.Handle(path_svg, svg_handler)

		if cfg.Verbose {
			log.Printf("handle svg endpoint at %s\n", path_svg)
		}

		for _, alt_path := range path_svg_alt {
			mux.Handle(alt_path, svg_handler)
		}
	}

	if cfg.EnableSPR {

		if cfg.Verbose {
			log.Println("spr support enabled")
		}

		spr_opts := &www.SPRHandlerOptions{
			Reader: cr,
			Logger: logger,
		}

		spr_handler, err := www.SPRHandler(spr_opts)

		if err != nil {
			return fmt.Errorf("Failed to create SPR handler, %w", err)
		}

		if cfg.EnableCORS {
			spr_handler = cors_wrapper.Handler(spr_handler)
		}

		mux.Handle(path_spr, spr_handler)

		if cfg.Verbose {
			log.Printf("handle spr endpoint at %s\n", path_spr)
		}

		for _, alt_path := range path_spr_alt {
			mux.Handle(alt_path, spr_handler)
		}
	}

	if cfg.EnableGeoJSON {

		if cfg.Verbose {
			log.Println("geojson support enabled")
		}

		geojson_opts := &www.GeoJSONHandlerOptions{
			Reader: cr,
			Logger: logger,
		}

		geojson_handler, err := www.GeoJSONHandler(geojson_opts)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON handler, %w", err)
		}

		if cfg.EnableCORS {
			geojson_handler = cors_wrapper.Handler(geojson_handler)
		}

		mux.Handle(path_geojson, geojson_handler)

		if cfg.Verbose {
			logger.Printf("Handle geojson endpoint at %s\n", path_geojson)
		}

		for _, alt_path := range path_geojson_alt {
			mux.Handle(alt_path, geojson_handler)
		}
	}

	if cfg.EnableGeoJSONLD {

		if cfg.Verbose {
			log.Println("geojsonld support enabled")
		}

		geojsonld_opts := &www.GeoJSONLDHandlerOptions{
			Reader: cr,
			Logger: logger,
		}

		geojsonld_handler, err := www.GeoJSONLDHandler(geojsonld_opts)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON LD handler, %w", err)
		}

		if cfg.EnableCORS {
			geojsonld_handler = cors_wrapper.Handler(geojsonld_handler)
		}

		mux.Handle(path_geojsonld, geojsonld_handler)

		if cfg.Verbose {
			logger.Printf("Handle geojsonld endpoint at %s\n", path_geojsonld)
		}

		for _, alt_path := range path_geojsonld_alt {
			mux.Handle(alt_path, geojsonld_handler)
		}
	}

	if cfg.EnableNavPlace {

		if cfg.Verbose {
			log.Println("navplace support enabled")
		}

		navplace_opts := &www.NavPlaceHandlerOptions{
			Reader:      cr,
			MaxFeatures: navplace_max_features,
			Logger:      logger,
		}

		navplace_handler, err := www.NavPlaceHandler(navplace_opts)

		if err != nil {
			return fmt.Errorf("Failed to create IIIF navPlace handler, %w", err)
		}

		if cfg.EnableCORS {
			navplace_handler = cors_wrapper.Handler(navplace_handler)
		}

		mux.Handle(path_navplace, navplace_handler)

		if cfg.Verbose {
			logger.Printf("Handle navplace endpoint at %s\n", path_navplace)
		}

		for _, alt_path := range path_navplace_alt {
			mux.Handle(alt_path, navplace_handler)
		}
	}

	if cfg.EnableSelect {

		if cfg.Verbose {
			log.Println("select support enabled")
		}

		if select_pattern == "" {
			return fmt.Errorf("Missing -select-pattern parameter.")
		}

		pat, err := regexp.Compile(select_pattern)

		if err != nil {
			return fmt.Errorf("Failed to compile select pattern (%s), %w", select_pattern, err)
		}

		select_opts := &www.SelectHandlerOptions{
			Pattern: pat,
			Reader:  cr,
			Logger:  logger,
		}

		select_handler, err := www.SelectHandler(select_opts)

		if err != nil {
			return fmt.Errorf("Failed to create select handler, %w", err)
		}

		if cfg.EnableCORS {
			select_handler = cors_wrapper.Handler(select_handler)
		}

		mux.Handle(path_select, select_handler)

		if cfg.Verbose {
			log.Printf("handle select endpoint at %s\n", path_select)
		}

		for _, alt_path := range path_select_alt {
			mux.Handle(alt_path, select_handler)
		}
	}

	if cfg.EnableWebFinger {

		if cfg.Verbose {
			log.Println("webfinger support enabled")
		}

		webfinger_opts := &www.WebfingerHandlerOptions{
			Reader:       cr,
			Logger:       logger,
			Paths:        www_paths,
			Capabilities: www_capabilities,
			Hostname:     webfinger_hostname,
		}

		webfinger_handler, err := www.WebfingerHandler(webfinger_opts)

		if err != nil {
			return fmt.Errorf("Failed to create webfinger handler, %w", err)
		}

		if cfg.EnableCORS {
			webfinger_handler = cors_wrapper.Handler(webfinger_handler)
		}

		mux.Handle(path_webfinger, webfinger_handler)

		if cfg.Verbose {
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
			search_opts.GeoJSONReader = cr
		}

		search_handler, err := www.SearchAPIHandler(search_opts)

		if err != nil {
			return fmt.Errorf("Failed to create search handler, %w", err)
		}

		if cfg.EnableCORS {
			search_handler = cors_wrapper.Handler(search_handler)
		}

		mux.Handle(path_search_api, search_handler)
	}

	// END OF probably due a rethink shortly

	// Common code for HTML handler (public and/or edit handlers)

	var bootstrap_opts *bootstrap.BootstrapOptions
	var map_provider provider.Provider
	var maps_opts *maps.MapsOptions

	if cfg.EnableHTML || cfg.EnableEditUI {

		bootstrap_opts = bootstrap.DefaultBootstrapOptions()

		err = bootstrap.AppendAssetHandlersWithPrefix(mux, cfg.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append Bootstrap asset handlers, %w", err)
		}

		err = www.AppendStaticAssetHandlersWithPrefix(mux, cfg.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %w", err)
		}

		err = custom.AppendStaticAssetHandlersWithPrefix(mux, cfg.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append custom asset handlers, %w", err)
		}

		provider_uri, err := provider.ProviderURIFromFlagSet(fs)

		if err != nil {
			return fmt.Errorf("Failed to derive provider URI from flagset, %v", err)
		}

		pr, err := provider.NewProvider(ctx, provider_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new provider, %w", err)
		}

		map_provider = pr
		err = map_provider.SetLogger(logger)

		if err != nil {
			return fmt.Errorf("Failed to set logger for provider, %w", err)
		}

		// Final map stuff

		maps_opts = maps.DefaultMapsOptions()

		err = map_www.AppendStaticAssetHandlersWithPrefix(mux, cfg.URIPrefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %v")
		}

		err = map_provider.AppendAssetHandlersWithPrefix(mux, cfg.URIPrefix)

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

	if cfg.EnableHTML {

		// Note that we append all the handler to mux at the end of this block so that they
		// can be updated with map-related middleware where necessary

		var index_handler http.Handler
		var id_handler http.Handler
		var search_handler http.Handler

		if enable_index {

			index_opts := www.IndexHandlerOptions{
				Templates:    t,
				Paths:        www_paths,
				Capabilities: www_capabilities,
			}

			index_h, err := www.IndexHandler(index_opts)

			if err != nil {
				return fmt.Errorf("Failed to create index handler, %w", err)
			}

			index_handler = index_h
			index_handler = bootstrap.AppendResourcesHandlerWithPrefix(index_handler, bootstrap_opts, cfg.URIPrefix)
		}

		id_opts := www.IDHandlerOptions{
			Templates:    t,
			Paths:        www_paths,
			Capabilities: www_capabilities,
			Reader:       cr,
			Logger:       logger,
			MapProvider:  map_provider.Scheme(),
		}

		id_h, err := www.IDHandler(id_opts)

		if err != nil {
			return fmt.Errorf("Failed to create ID handler, %w", err)
		}

		id_handler = id_h
		id_handler = bootstrap.AppendResourcesHandlerWithPrefix(id_handler, bootstrap_opts, cfg.URIPrefix)

		if enable_search_html {

			search_db, err := fulltext.NewFullTextDatabase(ctx, search_database_uri)

			if err != nil {
				return fmt.Errorf("Failed to create fulltext database for '%s', %w", search_database_uri, err)
			}

			search_opts := www.SearchHandlerOptions{
				Templates:    t,
				Paths:        www_paths,
				Capabilities: www_capabilities,
				Database:     search_db,
				MapProvider:  map_provider.Scheme(),
			}

			search_h, err := www.SearchHandler(search_opts)

			if err != nil {
				return fmt.Errorf("Failed to create search handler, %w", err)
			}

			search_handler = search_h
			search_handler = bootstrap.AppendResourcesHandlerWithPrefix(search_handler, bootstrap_opts, cfg.URIPrefix)
		}

		id_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(id_handler, map_provider, maps_opts, cfg.URIPrefix)
		id_handler = custom.WrapHandler(id_handler)
		id_handler = authenticator.WrapHandler(id_handler)

		mux.Handle(path_id, id_handler)

		if cfg.Verbose {
			log.Printf("handle ID endpoint at %s\n", path_id)
		}

		if enable_search_html {
			search_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(search_handler, map_provider, maps_opts, cfg.URIPrefix)
			search_handler = authenticator.WrapHandler(search_handler)
			mux.Handle(path_search_html, search_handler)
		}

		index_handler = authenticator.WrapHandler(index_handler)
		mux.Handle("/", index_handler)
	}

	// Edit/write HTML handlers

	if cfg.EnableEditUI {

		if cfg.Verbose {
			log.Println("edit user interface support enabled")
		}

		// Edit geometry

		geom_t := t.Lookup("geometry")

		if geom_t == nil {
			return fmt.Errorf("Failed to load 'geometry' template")
		}

		geom_opts := &www.EditGeometryHandlerOptions{
			Authenticator: authenticator,
			MapProvider:   map_provider.Scheme(),
			Paths:         www_paths,
			Capabilities:  www_capabilities,
			Template:      geom_t,
			Logger:        logger,
			Reader:        cr,
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
			Authenticator: authenticator,
			MapProvider:   map_provider.Scheme(),
			Paths:         www_paths,
			Capabilities:  www_capabilities,
			Template:      create_t,
			Logger:        logger,
			Reader:        cr,
		}

		create_handler, err := www.CreateFeatureHandler(create_opts)

		if err != nil {
			return fmt.Errorf("Failed to create create feature handler, %w", err)
		}

		geom_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(geom_handler, map_provider, maps_opts, cfg.URIPrefix)
		geom_handler = bootstrap.AppendResourcesHandlerWithPrefix(geom_handler, bootstrap_opts, cfg.URIPrefix)
		geom_handler = custom.WrapHandler(geom_handler)
		geom_handler = authenticator.WrapHandler(geom_handler)

		mux.Handle(path_edit_geometry, geom_handler)

		if cfg.Verbose {
			log.Printf("handle edit geometry endpoint at %s\n", path_edit_geometry)
		}

		create_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(create_handler, map_provider, maps_opts, cfg.URIPrefix)
		create_handler = bootstrap.AppendResourcesHandlerWithPrefix(create_handler, bootstrap_opts, cfg.URIPrefix)
		create_handler = custom.WrapHandler(create_handler)
		create_handler = authenticator.WrapHandler(create_handler)

		mux.Handle(path_create_feature, create_handler)

		if cfg.Verbose {
			log.Printf("handle create feature endpoint at %s\n", path_create_feature)
		}
	}

	if cfg.EnableEditAPI {

		if cfg.Verbose {
			log.Println("edit api support enabled")
		}

		// Exporter

		ex, err := export.NewExporter(ctx, exporter_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new exporter, %w", err)
		}

		// Writers are created at runtime using the http/api/publish.go#publishFeature
		// method which in turn calls writer/writer.go#NewWriter

		// START OF Point in polygon service setup

		// We're doing this the long way because we need/want to pass in 'cr' and I am
		// not sure what the interface/signature changes to do that should be...

		spatial_db, err := database.NewSpatialDatabase(ctx, spatial_database_uri)

		if err != nil {
			return fmt.Errorf("Failed to create spatial database, %w", err)
		}

		pip_service, err := pointinpolygon.NewPointInPolygonServiceWithDatabaseAndReader(ctx, spatial_db, cr)

		if err != nil {
			return fmt.Errorf("Failed to create point in polygon service, %w", err)
		}

		// END OF Point in polygon service setup

		// Deprecate a record

		deprecate_opts := &api.DeprecateFeatureHandlerOptions{
			Reader:        cr,
			Cache:         browser_cache,
			Logger:        logger,
			Authenticator: authenticator,
			Exporter:      ex,
			WriterURIs:    writer_uris,
		}

		deprecate_handler, err := api.DeprecateFeatureHandler(deprecate_opts)

		if err != nil {
			return fmt.Errorf("Failed to create deprecate feature handler, %w", err)
		}

		deprecate_handler = authenticator.WrapHandler(deprecate_handler)
		mux.Handle(path_api_deprecate_feature, deprecate_handler)

		if cfg.Verbose {
			log.Printf("handle deprecate feature endpoint at %s\n", path_api_deprecate_feature)
		}

		// Mark a record as ceased

		cessate_opts := &api.CessateFeatureHandlerOptions{
			Reader:        cr,
			Cache:         browser_cache,
			Logger:        logger,
			Authenticator: authenticator,
			Exporter:      ex,
			WriterURIs:    writer_uris,
		}

		cessate_handler, err := api.CessateFeatureHandler(cessate_opts)

		if err != nil {
			return fmt.Errorf("Failed to create cessate feature handler, %w", err)
		}

		cessate_handler = authenticator.WrapHandler(cessate_handler)
		mux.Handle(path_api_cessate_feature, cessate_handler)

		if cfg.Verbose {
			log.Printf("handle cessate feature endpoint at %s\n", path_api_cessate_feature)
		}

		// Edit geometry

		geom_opts := &api.UpdateGeometryHandlerOptions{
			Reader:                cr,
			Cache:                 browser_cache,
			Logger:                logger,
			Authenticator:         authenticator,
			Exporter:              ex,
			WriterURIs:            writer_uris,
			PointInPolygonService: pip_service,
		}

		geom_handler, err := api.UpdateGeometryHandler(geom_opts)

		if err != nil {
			return fmt.Errorf("Failed to create uupdate geometry handler, %w", err)
		}

		geom_handler = authenticator.WrapHandler(geom_handler)
		mux.Handle(path_api_edit_geometry, geom_handler)

		if cfg.Verbose {
			log.Printf("handle edit geometry endpoint at %s\n", path_api_edit_geometry)
		}

		// Create a new feature

		create_opts := &api.CreateFeatureHandlerOptions{
			Reader:                cr,
			Cache:                 browser_cache,
			Logger:                logger,
			Authenticator:         authenticator,
			Exporter:              ex,
			WriterURIs:            writer_uris,
			PointInPolygonService: pip_service,
		}

		create_handler, err := api.CreateFeatureHandler(create_opts)

		if err != nil {
			return fmt.Errorf("Failed to create create feature handler, %w", err)
		}

		create_handler = authenticator.WrapHandler(create_handler)
		mux.Handle(path_api_create_feature, create_handler)

		if cfg.Verbose {
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

func RunWithConfig(ctx context.Context, cfg *Config, logger *log.Logger) error {
	return fmt.Errorf("Not implemented")
}

func RunWithSettings(ctx context.Context, settings *Settings) error {
	return fmt.Errorf("Not implemented")
}
