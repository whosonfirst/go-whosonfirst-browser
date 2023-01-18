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
	"io/ioutil"
	"log"
	"net/http"
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
	"github.com/whosonfirst/go-whosonfirst-browser/v6/chrome"
	"github.com/whosonfirst/go-whosonfirst-browser/v6/http/api"
	"github.com/whosonfirst/go-whosonfirst-browser/v6/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v6/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-browser/v6/templates/html"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
)

func Run(ctx context.Context, logger *log.Logger) error {

	fs, err := DefaultFlagSet(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create default flagset, %w", err)
	}

	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "BROWSER")

	if err != nil {
		return fmt.Errorf("Failed to set flags from environment variables, %w", err)
	}

	if enable_all {
		enable_graphics = true
		enable_data = true
		enable_html = true
		// enable_search = true
	}

	if enable_search {
		enable_search_api = true
		enable_search_api_geojson = true
		enable_search_html = true
	}

	if enable_graphics {
		enable_png = true
		enable_svg = true
	}

	if enable_data {
		enable_geojson = true
		enable_geojsonld = true
		enable_navplace = true
		enable_spr = true
		enable_select = true
		enable_webfinger = true
	}

	if enable_search_html {
		enable_html = true
	}

	if enable_html {
		enable_geojson = true
		enable_png = true
	}

	var cors_wrapper *cors.Cors

	if enable_cors {

		if len(cors_origins) == 0 {
			cors_origins.Set("*")
		}

		cors_wrapper = cors.New(cors.Options{
			AllowedOrigins: cors_origins,
		})
	}

	if cache_uri == "tmp://" {

		now := time.Now()
		prefix := fmt.Sprintf("go-whosonfirst-browser-%d", now.Unix())

		tempdir, err := ioutil.TempDir("", prefix)

		if err != nil {
			return fmt.Errorf("Failed to derive tmp dir, %w", err)
		}

		defer os.RemoveAll(tempdir)

		cache_uri = fmt.Sprintf("fs://%s", tempdir)
	}

	if github_accesstoken_uri != "" {

		for idx, r_uri := range reader_uris {

			r_uri, err := github_reader.EnsureGitHubAccessToken(ctx, r_uri, github_accesstoken_uri)

			if err != nil {
				return fmt.Errorf("Failed to ensure GitHub access token for '%s', %w", r_uri, err)
			}

			reader_uris[idx] = r_uri
		}
	}

	browser_reader, err := reader.NewMultiReaderFromURIs(ctx, reader_uris...)

	if err != nil {
		return fmt.Errorf("Failed to create reader, %w", err)
	}

	browser_cache, err := cache.NewCache(ctx, cache_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new cache, %w", err)
	}

	cr_opts := &cachereader.CacheReaderOptions{
		Reader: browser_reader,
		Cache:  browser_cache,
	}

	cr, err := cachereader.NewCacheReaderWithOptions(ctx, cr_opts)

	if err != nil {
		return fmt.Errorf("Failed to create cache reader, %w", err)
	}

	t, err := html.LoadTemplates(ctx)

	if err != nil {
		return fmt.Errorf("Failed to load templates, %w", err)
	}

	if static_prefix != "" {

		static_prefix = strings.TrimRight(static_prefix, "/")

		if !strings.HasPrefix(static_prefix, "/") {
			return fmt.Errorf("Invalid -static-prefix value")
		}
	}

	authenticator, err := auth.NewAuthenticator(ctx, authenticator_uri)

	if err != nil {
		return fmt.Errorf("Failed to create authenticator, %w", err)
	}

	custom, err := chrome.NewChrome(ctx, custom_chrome_uri)

	if err != nil {
		return fmt.Errorf("Failed to create custom chrome, %w", err)
	}

	www_paths := &www.Paths{
		GeoJSON:   path_geojson,
		GeoJSONLD: path_geojsonld,
		SVG:       path_svg,
		PNG:       path_png,
		Select:    path_select,
		NavPlace:  path_navplace,
		SPR:       path_spr,
		HTML:      path_id,
	}

	www_capabilities := &www.Capabilities{
		HTML:      enable_html,
		GeoJSON:   enable_geojson,
		GeoJSONLD: enable_geojsonld,
		SVG:       enable_svg,
		PNG:       enable_png,
		Select:    enable_select,
		NavPlace:  enable_navplace,
		SPR:       enable_spr,
	}

	mux := http.NewServeMux()

	ping_handler, err := ping.PingPongHandler()

	if err != nil {
		return fmt.Errorf("Failed to create ping handler, %w", err)
	}

	mux.Handle("/ping", ping_handler)

	if enable_png {

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

		for _, alt_path := range path_png_alt {
			mux.Handle(alt_path, png_handler)
		}
	}

	if enable_svg {

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

		if enable_cors {
			svg_handler = cors_wrapper.Handler(svg_handler)
		}

		mux.Handle(path_svg, svg_handler)

		for _, alt_path := range path_svg_alt {
			mux.Handle(alt_path, svg_handler)
		}
	}

	if enable_spr {

		spr_opts := &www.SPRHandlerOptions{
			Reader: cr,
			Logger: logger,
		}

		spr_handler, err := www.SPRHandler(spr_opts)

		if err != nil {
			return fmt.Errorf("Failed to create SPR handler, %w", err)
		}

		if enable_cors {
			spr_handler = cors_wrapper.Handler(spr_handler)
		}

		mux.Handle(path_spr, spr_handler)

		for _, alt_path := range path_spr_alt {
			mux.Handle(alt_path, spr_handler)
		}
	}

	if enable_geojson {

		geojson_opts := &www.GeoJSONHandlerOptions{
			Reader: cr,
			Logger: logger,
		}

		geojson_handler, err := www.GeoJSONHandler(geojson_opts)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON handler, %w", err)
		}

		if enable_cors {
			geojson_handler = cors_wrapper.Handler(geojson_handler)
		}

		mux.Handle(path_geojson, geojson_handler)

		for _, alt_path := range path_geojson_alt {
			mux.Handle(alt_path, geojson_handler)
		}
	}

	if enable_geojsonld {

		geojsonld_opts := &www.GeoJSONLDHandlerOptions{
			Reader: cr,
			Logger: logger,
		}

		geojsonld_handler, err := www.GeoJSONLDHandler(geojsonld_opts)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON LD handler, %w", err)
		}

		if enable_cors {
			geojsonld_handler = cors_wrapper.Handler(geojsonld_handler)
		}

		mux.Handle(path_geojsonld, geojsonld_handler)

		for _, alt_path := range path_geojsonld_alt {
			mux.Handle(alt_path, geojsonld_handler)
		}
	}

	if enable_navplace {

		navplace_opts := &www.NavPlaceHandlerOptions{
			Reader:      cr,
			MaxFeatures: navplace_max_features,
			Logger:      logger,
		}

		navplace_handler, err := www.NavPlaceHandler(navplace_opts)

		if err != nil {
			return fmt.Errorf("Failed to create IIIF navPlace handler, %w", err)
		}

		if enable_cors {
			navplace_handler = cors_wrapper.Handler(navplace_handler)
		}

		mux.Handle(path_navplace, navplace_handler)

		for _, alt_path := range path_navplace_alt {
			mux.Handle(alt_path, navplace_handler)
		}
	}

	if enable_select {

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

		if enable_cors {
			select_handler = cors_wrapper.Handler(select_handler)
		}

		mux.Handle(path_select, select_handler)

		for _, alt_path := range path_select_alt {
			mux.Handle(alt_path, select_handler)
		}
	}

	if enable_webfinger {

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

		if enable_cors {
			webfinger_handler = cors_wrapper.Handler(webfinger_handler)
		}

		mux.Handle(path_webfinger, webfinger_handler)

		for _, alt_path := range path_webfinger_alt {
			mux.Handle(alt_path, webfinger_handler)
		}
	}

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

			/*
				if resolver_uri != "" {

				resolver_func, err := geojson.NewSPRPathResolverFunc(ctx, resolver_uri)

				if err != nil {
					return err
				}

				api_pip_opts.SPRPathResolver = resolver_func
			*/
		}

		search_handler, err := www.SearchAPIHandler(search_opts)

		if err != nil {
			return fmt.Errorf("Failed to create search handler, %w", err)
		}

		if enable_cors {
			search_handler = cors_wrapper.Handler(search_handler)
		}

		mux.Handle(path_search_api, search_handler)
	}

	if enable_html {

		endpoints := &www.Endpoints{
			Data:  path_geojson,
			Png:   path_png,
			Svg:   path_svg,
			Spr:   path_spr,
			Id:    path_id,
			Index: "/",
		}

		if enable_search_html {
			endpoints.Search = path_search_html
		}

		bootstrap_opts := bootstrap.DefaultBootstrapOptions()

		err = bootstrap.AppendAssetHandlersWithPrefix(mux, static_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append Bootstrap asset handlers, %w", err)
		}

		err = www.AppendStaticAssetHandlersWithPrefix(mux, static_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %w", err)
		}

		err = custom.AppendStaticAssetHandlersWithPrefix(mux, static_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append custom asset handlers, %w", err)
		}

		provider_uri, err := provider.ProviderURIFromFlagSet(fs)

		if err != nil {
			return fmt.Errorf("Failed to derive provider URI from flagset, %v", err)
		}

		map_provider, err := provider.NewProvider(ctx, provider_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new provider, %w", err)
		}

		err = map_provider.SetLogger(logger)

		if err != nil {
			return fmt.Errorf("Failed to set logger for provider, %w", err)
		}

		// Note that we append all the handler to mux at the end of this block so that they
		// can be updated with map-related middleware where necessary

		var index_handler http.Handler
		var id_handler http.Handler
		var search_handler http.Handler

		if enable_index {

			index_opts := www.IndexHandlerOptions{
				Templates: t,
				Endpoints: endpoints,
			}

			index_h, err := www.IndexHandler(index_opts)

			if err != nil {
				return fmt.Errorf("Failed to create index handler, %w", err)
			}

			index_handler = index_h
			index_handler = bootstrap.AppendResourcesHandlerWithPrefix(index_handler, bootstrap_opts, static_prefix)
		}

		id_opts := www.IDHandlerOptions{
			Templates:   t,
			Endpoints:   endpoints,
			Reader:      cr,
			Logger:      logger,
			MapProvider: map_provider.Scheme(),
		}

		id_h, err := www.IDHandler(id_opts)

		if err != nil {
			return fmt.Errorf("Failed to create ID handler, %w", err)
		}

		id_handler = id_h
		id_handler = bootstrap.AppendResourcesHandlerWithPrefix(id_handler, bootstrap_opts, static_prefix)

		if enable_search_html {

			search_db, err := fulltext.NewFullTextDatabase(ctx, search_database_uri)

			if err != nil {
				return fmt.Errorf("Failed to create fulltext database for '%s', %w", search_database_uri, err)
			}

			search_opts := www.SearchHandlerOptions{
				Templates:   t,
				Endpoints:   endpoints,
				Database:    search_db,
				MapProvider: map_provider.Scheme(),
			}

			search_h, err := www.SearchHandler(search_opts)

			if err != nil {
				return fmt.Errorf("Failed to create search handler, %w", err)
			}

			search_handler = search_h
			search_handler = bootstrap.AppendResourcesHandlerWithPrefix(search_handler, bootstrap_opts, static_prefix)
		}

		// Edit geometry

		geom_t := t.Lookup("geometry")

		if geom_t == nil {
			return fmt.Errorf("Failed to load 'geometry' template")
		}

		geom_opts := &www.EditGeometryHandlerOptions{
			Authenticator: authenticator,
			MapProvider:   map_provider.Scheme(),
			Endpoints:     endpoints,
			Template:      geom_t,
			Logger:        logger,
			Reader:        cr,
		}

		geom_handler, err := www.EditGeometryHandler(geom_opts)

		if err != nil {
			return fmt.Errorf("Failed to create edit geometry handler, %w", err)
		}

		maps_opts := maps.DefaultMapsOptions()

		// TBD... update maps_opts... but really only (n) handlers
		// maps_opts = custom.UpdateMapsOptions(maps_opts)
		
		err = map_www.AppendStaticAssetHandlersWithPrefix(mux, static_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %v")
		}

		err = map_provider.AppendAssetHandlersWithPrefix(mux, static_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append provider asset handlers, %v", err)
		}

		id_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(id_handler, map_provider, maps_opts, static_prefix)
		id_handler = authenticator.WrapHandler(id_handler)

		mux.Handle(path_id, id_handler)

		if enable_search_html {
			search_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(search_handler, map_provider, maps_opts, static_prefix)
			search_handler = authenticator.WrapHandler(search_handler)
			mux.Handle(path_search_html, search_handler)
		}

		geom_handler = maps.AppendResourcesHandlerWithPrefixAndProvider(geom_handler, map_provider, maps_opts, static_prefix)

		// TBD: custom middleware goes here?
		geom_handler = custom.WrapHandler(geom_handler)

		geom_handler = authenticator.WrapHandler(geom_handler)

		mux.Handle(path_edit_geometry, geom_handler)

		// log.Printf("Enable edit geometry handler at %s", path_edit_geometry)

		index_handler = authenticator.WrapHandler(index_handler)
		mux.Handle("/", index_handler)

		null_handler := www.NewNullHandler()

		favicon_path := filepath.Join(path_id, "favicon.ico")
		mux.Handle(favicon_path, null_handler)

	}

	if enable_api {

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
		mux.Handle(path_api_deprecate, deprecate_handler)

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
		mux.Handle(path_api_cessate, cessate_handler)

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
			return fmt.Errorf("Failed to create cessate feature handler, %w", err)
		}

		geom_handler = authenticator.WrapHandler(geom_handler)
		mux.Handle(path_api_edit_geometry, geom_handler)

		// logger.Printf("Listening at %s\n", path_api_edit_geometry)
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
