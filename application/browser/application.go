package browser

import (
	_ "github.com/whosonfirst/go-reader-cachereader"
	_ "github.com/whosonfirst/go-reader-findingaid"
)

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-ping/v2"
	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/sfomuseum/go-flags/flagset"
	tzhttp "github.com/sfomuseum/go-http-tilezen/http"
	tiles_http "github.com/tilezen/go-tilepacks/http"
	"github.com/tilezen/go-tilepacks/tilepack"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/v4/application"
	"github.com/whosonfirst/go-whosonfirst-browser/v4/templates/html"
	"github.com/whosonfirst/go-whosonfirst-browser/v4/www"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
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

// type BrowserApplication implements the application.Application and provides a web-based application for browsing Who's On First records in different formats.
type BrowserApplication struct {
	application.Application
}

// NewBrowserApplication will return a new application.Application instance implementing the BrowserApplication application.
func NewBrowserApplication(ctx context.Context) (application.Application, error) {
	app := &BrowserApplication{}
	return app, nil
}

// DefaultFlagSet returns a `flag.FlagSet` instance with flags and defaults values assigned for use with `app`.
func (app *BrowserApplication) DefaultFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("browser")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.StringVar(&static_prefix, "static-prefix", "", "Prepend this prefix to URLs for static assets.")

	fs.Var(&reader_uris, "reader-uri", "One or more valid go-reader Reader URI strings.")
	fs.StringVar(&cache_uri, "cache-uri", "gocache://", "A valid go-cache Cache URI string.")

	fs.StringVar(&nextzen_api_key, "nextzen-api-key", "", "A valid Nextzen API key (https://developers.nextzen.org/).")
	fs.StringVar(&nextzen_style_url, "nextzen-style-url", "/tangram/refill-style.zip", "A valid Tangram scene file URL.")
	fs.StringVar(&nextzen_tile_url, "nextzen-tile-url", tangramjs.NEXTZEN_MVT_ENDPOINT, "A valid Nextzen MVT tile URL.")

	fs.StringVar(&tilepack_db, "nextzen-tilepack-database", "", "The path to a valid MBTiles database (tilepack) containing Nextzen MVT tiles.")
	fs.StringVar(&tilepack_uri, "nextzen-tilepack-uri", "/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt", "The relative URI to serve Nextzen MVT tiles from a MBTiles database (tilepack).")

	fs.BoolVar(&proxy_tiles, "proxy-tiles", false, "Proxy (and cache) Nextzen tiles.")
	fs.StringVar(&proxy_tiles_url, "proxy-tiles-url", "/tiles/", "The URL (a relative path) for proxied tiles.")
	fs.StringVar(&proxy_tiles_cache, "proxy-tiles-cache", "gocache://", "A valid tile proxy DSN string.")
	fs.IntVar(&proxy_tiles_timeout, "proxy-tiles-timeout", 30, "The maximum number of seconds to allow for fetching a tile from the proxy.")

	fs.BoolVar(&enable_all, "enable-all", false, "Enable all the available output handlers EXCEPT the search handlers which need to be explicitly enable using the -enable-search* flags.")
	fs.BoolVar(&enable_graphics, "enable-graphics", false, "Enable the 'png' and 'svg' output handlers.")
	fs.BoolVar(&enable_data, "enable-data", false, "Enable the 'geojson' and 'spr' and 'select' output handlers.")

	fs.BoolVar(&enable_png, "enable-png", false, "Enable the 'png' output handler.")
	fs.BoolVar(&enable_svg, "enable-svg", false, "Enable the 'svg' output handler.")

	fs.BoolVar(&enable_geojson, "enable-geojson", true, "Enable the 'geojson' output handler.")
	fs.BoolVar(&enable_geojsonld, "enable-geojson-ld", true, "Enable the 'geojson-ld' output handler.")
	fs.BoolVar(&enable_navplace, "enable-navplace", true, "Enable the IIIF 'navPlace' output handler.")
	fs.BoolVar(&enable_spr, "enable-spr", true, "Enable the 'spr' (or \"standard places response\") output handler.")
	fs.BoolVar(&enable_select, "enable-select", false, "Enable the 'select' output handler.")
	fs.StringVar(&select_pattern, "select-pattern", "properties(?:.[a-zA-Z0-9-_]+){1,}", "A valid regular expression for sanitizing select parameters.")

	fs.BoolVar(&enable_html, "enable-html", true, "Enable the 'html' (or human-friendly) output handlers.")

	fs.BoolVar(&enable_search_api, "enable-search-api", false, "Enable the (API) search handlers.")
	fs.BoolVar(&enable_search_api_geojson, "enable-search-api-geojson", false, "Enable the (API) search handlers to return results as GeoJSON.")
	fs.BoolVar(&enable_search_html, "enable-search-html", false, "Enable the (human-friendly) search handlers.")
	fs.BoolVar(&enable_search, "enable-search", false, "Enable both the API and human-friendly search handlers.")
	fs.StringVar(&search_database_uri, "search-database-uri", "", "A valid whosonfirst/go-whosonfist-search/fulltext URI.")

	fs.StringVar(&path_png, "path-png", "/png/", "The path that PNG requests should be served from.")
	fs.StringVar(&path_svg, "path-svg", "/svg/", "The path that SVG requests should be served from.")
	fs.StringVar(&path_geojson, "path-geojson", "/geojson/", "The path that GeoJSON requests should be served from.")
	fs.StringVar(&path_geojsonld, "path-geojson-ld", "/geojson-ld/", "The path that GeoJSON-LD requests should be served from.")
	fs.StringVar(&path_navplace, "path-navplace", "/navplace/", "The path that IIIF navPlace requests should be served from.")
	fs.StringVar(&path_spr, "path-spr", "/spr/", "The path that SPR requests should be served from.")
	fs.StringVar(&path_select, "path-select", "/select/", "The path that 'select' requests should be served from.")

	fs.StringVar(&path_search_api, "path-search-api", "/search/spr/", "The path that API 'search' requests should be served from.")
	fs.StringVar(&path_search_html, "path-search-html", "/search/", "The path that API 'search' requests should be served from.")

	fs.StringVar(&path_id, "path-id", "/id/", "The that Who's On First documents should be served from.")

	fs.IntVar(&navplace_max_features, "navplace-max-features", 3, "The maximum number of features to allow in a /navplace/{ID} URI string.")

	return fs, nil
}

// Run will run the `app` (BrowserApplication) using default flags and values.
func (app *BrowserApplication) Run(ctx context.Context) error {

	fs, err := app.DefaultFlagSet(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create default flagset, %w", err)
	}

	return app.RunWithFlagSet(ctx, fs)
}

// Run will run the `app` (BrowserApplication) using a custom `flag.FlagSet` instance.
func (app *BrowserApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

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

	log.Println("DEBUG", cr_uri.String())
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

	mux := http.NewServeMux()

	ping_handler, err := ping.PingPongHandler()

	if err != nil {
		return fmt.Errorf("Failed to create ping handler, %w", err)
	}

	mux.Handle("/ping", ping_handler)

	if enable_png {

		png_opts, err := www.NewDefaultRasterOptions()

		if err != nil {
			return fmt.Errorf("Failed to create raster/png options, %w", err)
		}

		png_handler, err := www.RasterHandler(cr, png_opts)

		if err != nil {
			return fmt.Errorf("Failed to create raster/png handler, %w", err)
		}

		mux.Handle(path_png, png_handler)
	}

	if enable_svg {

		svg_opts, err := www.NewDefaultSVGOptions()

		if err != nil {
			return fmt.Errorf("Failed to create SVG options, %w", err)
		}

		svg_handler, err := www.SVGHandler(cr, svg_opts)

		if err != nil {
			return fmt.Errorf("Failed to create SVG handler, %w", err)
		}

		mux.Handle(path_svg, svg_handler)
	}

	if enable_spr {

		spr_handler, err := www.SPRHandler(cr)

		if err != nil {
			return fmt.Errorf("Failed to create SPR handler, %w", err)
		}

		mux.Handle(path_spr, spr_handler)
	}

	if enable_geojson {

		geojson_handler, err := www.GeoJSONHandler(cr)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON handler, %w", err)
		}

		mux.Handle(path_geojson, geojson_handler)
	}

	if enable_geojsonld {

		geojsonld_handler, err := www.GeoJSONLDHandler(cr)

		if err != nil {
			return fmt.Errorf("Failed to create GeoJSON LD handler, %w", err)
		}

		mux.Handle(path_geojsonld, geojsonld_handler)
	}

	if enable_navplace {

		navplace_opts := &www.NavPlaceHandlerOptions{
			Reader:      cr,
			MaxFeatures: navplace_max_features,
		}

		navplace_handler, err := www.NavPlaceHandler(navplace_opts)

		if err != nil {
			return fmt.Errorf("Failed to create IIIF navPlace handler, %w", err)
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
		}

		select_handler, err := www.SelectHandler(cr, select_opts)

		if err != nil {
			return fmt.Errorf("Failed to create select handler, %w", err)
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

		id_opts := www.IDHandlerOptions{
			Templates: t,
			Endpoints: endpoints,
		}

		id_handler, err := www.IDHandler(cr, id_opts)

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
