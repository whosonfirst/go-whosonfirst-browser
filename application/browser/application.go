package browser

import (
	_ "github.com/whosonfirst/go-reader-cachereader"
	_ "github.com/whosonfirst/go-reader-findingaid"
	_ "github.com/aaronland/go-http-server-tsnet"
)

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-ping/v2"
	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-http-auth"
	tzhttp "github.com/sfomuseum/go-http-tilezen/http"
	tiles_http "github.com/tilezen/go-tilepacks/http"
	"github.com/tilezen/go-tilepacks/tilepack"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/v5/http/api"
	"github.com/whosonfirst/go-whosonfirst-browser/v5/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v5/templates/html"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-writer/v3"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "BROWSER", true)

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

		log.Println(tempdir)
		defer os.RemoveAll(tempdir)

		cache_uri = fmt.Sprintf("fs://%s", tempdir)
	}

	cr_q := url.Values{}

	// go-reader-cachereader is configured to accept multiple readers
	// and to manage them all using reader.NewMultiReader
	cr_q["reader"] = reader_uris

	cr_q.Set("cache", cache_uri)

	cr_uri := url.URL{}
	cr_uri.Scheme = "cachereader"
	cr_uri.RawQuery = cr_q.Encode()

	cr, err := reader.NewReader(ctx, cr_uri.String())

	if err != nil {
		return fmt.Errorf("Failed to create reader for '%s', %w", cr_uri.String(), err)
	}

	// start of sudo put me in a package

	t := template.New("whosonfirst-browser").Funcs(template.FuncMap{
		"Add": func(i int, offset int) int {
			return i + offset
		},
	})

	t, err = t.ParseFS(html.FS, "*.html")

	if err != nil {
		return fmt.Errorf("Failed to parse templates, %w", err)
	}

	// end of sudo put me in a package

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

		if proxy_tiles {

			tile_cache, err := cache.NewCache(ctx, proxy_tiles_cache)

			if err != nil {
				return fmt.Errorf("Failed to create proxy tiles cache for '%s', %w", proxy_tiles_cache, err)
			}

			timeout := time.Duration(proxy_tiles_timeout) * time.Second

			proxy_opts := &tzhttp.TilezenProxyHandlerOptions{
				Cache:   tile_cache,
				Timeout: timeout,
			}

			proxy_handler, err := tzhttp.TilezenProxyHandler(proxy_opts)

			if err != nil {
				return fmt.Errorf("Failed to create proxy tiles handler, %w", err)
			}

			// the order here is important - we don't have a general-purpose "add to
			// mux with prefix" function here, like we do in the tangram handler so
			// we set the nextzen tile url with proxy_tiles_url and then update it
			// (proxy_tiles_url) with a prefix if necessary (20190911/thisisaaronland)

			nextzen_tile_url = fmt.Sprintf("%s{z}/{x}/{y}.mvt", proxy_tiles_url)

			if static_prefix != "" {

				proxy_tiles_url = filepath.Join(static_prefix, proxy_tiles_url)

				if !strings.HasSuffix(proxy_tiles_url, "/") {
					proxy_tiles_url = fmt.Sprintf("%s/", proxy_tiles_url)
				}
			}

			mux.Handle(proxy_tiles_url, proxy_handler)
		}

		bootstrap_opts := bootstrap.DefaultBootstrapOptions()

		tangramjs_opts := tangramjs.DefaultTangramJSOptions()
		tangramjs_opts.NextzenOptions.APIKey = nextzen_api_key
		tangramjs_opts.NextzenOptions.StyleURL = nextzen_style_url
		tangramjs_opts.NextzenOptions.TileURL = nextzen_tile_url

		if tilepack_db != "" {
			tangramjs_opts.NextzenOptions.TileURL = tilepack_uri
		}

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

		if enable_index {
			index_opts := www.IndexHandlerOptions{
				Templates: t,
				Endpoints: endpoints,
			}

			index_handler, err := www.IndexHandler(index_opts)

			if err != nil {
				return fmt.Errorf("Failed to create index handler, %w", err)
			}

			index_handler = bootstrap.AppendResourcesHandlerWithPrefix(index_handler, bootstrap_opts, static_prefix)

			mux.Handle("/", index_handler)
		}

		id_opts := www.IDHandlerOptions{
			Templates: t,
			Endpoints: endpoints,
			Reader:    cr,
			Logger:    logger,
		}

		id_handler, err := www.IDHandler(id_opts)

		if err != nil {
			return fmt.Errorf("Failed to create ID handler, %w", err)
		}

		id_handler = bootstrap.AppendResourcesHandlerWithPrefix(id_handler, bootstrap_opts, static_prefix)
		id_handler = tangramjs.AppendResourcesHandlerWithPrefix(id_handler, tangramjs_opts, static_prefix)

		mux.Handle(path_id, id_handler)

		if enable_search_html {

			search_db, err := fulltext.NewFullTextDatabase(ctx, search_database_uri)

			if err != nil {
				return fmt.Errorf("Failed to create fulltext database for '%s', %w", search_database_uri, err)
			}

			search_opts := www.SearchHandlerOptions{
				Templates: t,
				Endpoints: endpoints,
				Database:  search_db,
			}

			search_handler, err := www.SearchHandler(search_opts)

			if err != nil {
				return fmt.Errorf("Failed to create search handler, %w", err)
			}

			search_handler = bootstrap.AppendResourcesHandlerWithPrefix(search_handler, bootstrap_opts, static_prefix)
			search_handler = tangramjs.AppendResourcesHandlerWithPrefix(search_handler, tangramjs_opts, static_prefix)

			mux.Handle(path_search_html, search_handler)
		}

		err = bootstrap.AppendAssetHandlersWithPrefix(mux, static_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append Bootstrap asset handlers, %w", err)
		}

		err = tangramjs.AppendAssetHandlersWithPrefix(mux, static_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append Tangram.js asset handlers, %w", err)
		}

		err = www.AppendStaticAssetHandlersWithPrefix(mux, static_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append static asset handlers, %w", err)
		}

		if tilepack_db != "" {

			tiles_reader, err := tilepack.NewMbtilesReader(tilepack_db)

			if err != nil {
				return fmt.Errorf("Failed to load tilepack, %v", err)
			}

			u := strings.TrimLeft(tilepack_uri, "/")
			p := strings.Split(u, "/")
			path_tiles := fmt.Sprintf("/%s/", p[0])

			tiles_handler := tiles_http.MbtilesHandler(tiles_reader)
			mux.Handle(path_tiles, tiles_handler)
		}

	}

	if enable_api {

		ex, err := export.NewExporter(ctx, exporter_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new exporter, %w", err)
		}

		writers := make([]writer.Writer, len(writer_uris))

		for idx, wr_uri := range writer_uris {

			wr, err := writer.NewWriter(ctx, wr_uri)

			if err != nil {
				return fmt.Errorf("Failed to create writer for '%s', %w", wr_uri, err)
			}

			writers[idx] = wr
		}

		multi_opts := &writer.MultiWriterOptions{
			Logger: logger,
			Writers: writers,
		}
			
		multi_wr, err := writer.NewMultiWriterWithOptions(ctx, multi_opts)

		if err != nil {
			return fmt.Errorf("Failed to create multi writer, %w", err)
		}
		
		deprecate_opts := &api.DeprecateFeatureHandlerOptions{
			Reader:        cr,
			Logger:        logger,
			Authenticator: authenticator,
			Exporter:      ex,
			Writer:        multi_wr,
		}

		deprecate_handler, err := api.DeprecateFeatureHandler(deprecate_opts)

		if err != nil {
			return fmt.Errorf("Failed to create deprecate feature handler, %w", err)
		}

		deprecate_handler = authenticator.WrapHandler(deprecate_handler)
		mux.Handle(path_api_deprecate, deprecate_handler)

		cessate_opts := &api.CessateFeatureHandlerOptions{
			Reader:        cr,
			Logger:        logger,
			Authenticator: authenticator,
			Exporter:      ex,
			Writer:        multi_wr,
		}

		cessate_handler, err := api.CessateFeatureHandler(cessate_opts)

		if err != nil {
			return fmt.Errorf("Failed to create cessate feature handler, %w", err)
		}

		cessate_handler = authenticator.WrapHandler(cessate_handler)
		mux.Handle(path_api_cessate, cessate_handler)

	}

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
