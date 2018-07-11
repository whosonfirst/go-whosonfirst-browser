package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/whosonfirst/go-http-mapzenjs"
	mapzenjs_utils "github.com/whosonfirst/go-http-mapzenjs/utils"
	fs_reader "github.com/whosonfirst/go-whosonfirst-readwrite-fs/reader"
	http_reader "github.com/whosonfirst/go-whosonfirst-readwrite-http/reader"
	mysql_reader "github.com/whosonfirst/go-whosonfirst-readwrite-mysql/reader"
	s3_config "github.com/whosonfirst/go-whosonfirst-readwrite-s3/config"
	s3_reader "github.com/whosonfirst/go-whosonfirst-readwrite-s3/reader"
	sqlite_reader "github.com/whosonfirst/go-whosonfirst-readwrite-sqlite/reader"
	"github.com/whosonfirst/go-whosonfirst-readwrite/cache"
	"github.com/whosonfirst/go-whosonfirst-readwrite/flags"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-readwrite/utils"
	"github.com/whosonfirst/go-whosonfirst-static/http"
	"log"
	gohttp "net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func main() {

	config := flag.String("config", "", "Read some or all flags from an ini-style config file. Values in the config file take precedence over command line flags.")
	section := flag.String("section", "wof-static", "A valid ini-style config file section.")

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var source = flag.String("source", "fs", "Valid sources are: fs, http, mysql, s3, sqlite")
	var source_dsn = flag.String("source-dsn", "", "A valid DSN string specific to the source you've chosen.")

	var cache_source = flag.String("cache", "null", "...")

	var cache_args flags.KeyValueArgs
	flag.Var(&cache_args, "cache-arg", "(0) or more user-defined '{KEY}={VALUE}' arguments to pass to the caching layer")

	var test_reader = flag.String("test-reader", "", "Perform some basic sanity checking on the reader at startup")

	var data_endpoint = flag.String("data-endpoint", "", "")

	var api_key = flag.String("nextzen-api-key", "xxxxxxx", "A valid Nextzen API key (https://developers.nextzen.org/).")

	var debug = flag.Bool("debug", false, "Enable debugging.")

	var enable_all = flag.Bool("enable-all", false, "Enable all the available output handlers.")
	var enable_graphics = flag.Bool("enable-graphics", false, "Enable the 'png' and 'svg' output handlers.")
	var enable_data = flag.Bool("enable-data", false, "Enable the 'geojson' and 'spr' output handlers.")

	var enable_png = flag.Bool("enable-png", false, "Enable the 'png' output handler.")
	var enable_svg = flag.Bool("enable-svg", false, "Enable the 'svg' output handler.")

	var enable_geojson = flag.Bool("enable-geojson", true, "Enable the 'geojson' output handler.")
	var enable_spr = flag.Bool("enable-spr", true, "Enable the 'spr' (or \"standard places response\") output handler.")

	var enable_html = flag.Bool("enable-html", true, "Enable the 'html' (or human-friendly) output handler.")

	flag.Parse()

	if *config != "" {

		cfg, err := ini.LoadSources(ini.LoadOptions{
			AllowBooleanKeys: true,
		}, *config)

		if err != nil {
			logger.Fatal("Unable to load config file because %s", err)
		}

		sect, err := cfg.GetSection(*section)

		if err != nil {
			logger.Fatal("Config file is missing '%s' section, %s", *section, err)
		}

		flag.VisitAll(func(fl *flag.Flag) {

			name := fl.Name

			if name == "section" {
				return
			}

			if sect.HasKey(name) {

				k := sect.Key(name)
				v := k.Value()

				flag.Set(name, v)

				logger.Status("Reset %s flag from config file", name)
			}
		})

	} else {

		flag.VisitAll(func(fl *flag.Flag) {

			name := fl.Name
			env_name := name

			env_name = strings.Replace(name, "-", "_", -1)
			env_name = strings.ToUpper(name)

			env_name = fmt.Sprintf("WOF_STATIC_%s", env_name)

			v, ok := os.LookupEnv(env_name)

			if ok {

				flag.Set(name, v)
				logger.Status("Reset %s flag from %s environment variable", name, env_name)
			}
		})
	}

	var r reader.Reader
	var e error

	switch *source {
	case "fs":
		r, e = fs_reader.NewFSReader(*source_dsn)
	case "http":
		r, e = http_reader.NewHTTPReader(*source_dsn)
	case "mysql":
		r, e = mysql_reader.NewMySQLGeoJSONReader(*source_dsn)
	case "s3":

		cfg, err := s3_config.NewS3ConfigFromString(*source_dsn)

		if e != nil {
			e = err
		} else {
			r, e = s3_reader.NewS3Reader(cfg)
		}
	case "sqlite":
		r, e = sqlite_reader.NewSQLiteReader(*source_dsn)
	default:
		e = errors.New("Unknown or unsupported source")
	}

	if e != nil {
		log.Fatal(e)
	}

	// all this cache stuff _will_ change as the different cache providers
	// get migrated in to discrete packages (20180419/thisisaaronland)

	c, err := cache.NewCacheFromSource(*cache_source, cache_args.ToMap())

	if err != nil {
		log.Fatal(err)
	}

	opts, err := utils.NewDefaultCacheReaderOptions()

	if err != nil {
		log.Fatal(err)
	}

	opts.Debug = *debug

	cr, err := utils.NewCacheReader(r, c, opts)

	if err != nil {
		log.Fatal(err)
	}

	if *test_reader != "" {

		_, err := utils.TestReader(cr, *test_reader)

		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {

		for {

			select {
			case <-time.After(1 * time.Minute):
				log.Printf("CACHE KEYS: %d HITS: %d MISSES: %d\n", c.Size(), c.Hits(), c.Misses())
			}
		}
	}()

	html_opts := http.NewDefaultHTMLOptions()

	if *data_endpoint != "" {

		_, err = url.Parse(*data_endpoint)

		if err != nil {
			log.Fatal(err)
		}

		html_opts.DataEndpoint = *data_endpoint
	}

	html_handler, err := http.HTMLHandler(cr, html_opts)

	if err != nil {
		log.Fatal(err)
	}

	mapzenjs_opts := mapzenjs.DefaultMapzenJSOptions()
	mapzenjs_opts.APIKey = *api_key

	mapzenjs_handler, err := mapzenjs.MapzenJSHandler(html_handler, mapzenjs_opts)

	if err != nil {
		log.Fatal(err)
	}

	var png_handler gohttp.Handler
	var svg_handler gohttp.Handler

	var geojson_handler gohttp.Handler
	var spr_handler gohttp.Handler

	mux := gohttp.NewServeMux()

	ping_handler, err := http.PingHandler()

	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/ping", ping_handler)

	if *enable_all || *enable_graphics || *enable_png {

		png_opts, err := http.NewDefaultRasterOptions()

		if err != nil {
			log.Fatal(err)
		}

		h, err := http.RasterHandler(cr, png_opts)

		if err != nil {
			log.Fatal(err)
		}

		png_handler = h
		mux.Handle("/png/", png_handler)
	}

	if *enable_all || *enable_graphics || *enable_svg {

		svg_opts, err := http.NewDefaultSVGOptions()

		if err != nil {
			log.Fatal(err)
		}

		h, err := http.SVGHandler(cr, svg_opts)

		if err != nil {
			log.Fatal(err)
		}

		svg_handler = h
		mux.Handle("/svg/", svg_handler)
	}

	if *enable_all || *enable_data || *enable_spr {

		h, err := http.SPRHandler(cr)

		if err != nil {
			log.Fatal(err)
		}

		spr_handler = h
		mux.Handle("/spr/", spr_handler)
	}

	if *enable_all || *enable_data || *enable_geojson {

		h, err := http.GeoJSONHandler(cr)

		if err != nil {
			log.Fatal(err)
		}

		geojson_handler = h
		mux.Handle("/data/", geojson_handler)
	}

	if *enable_all || *enable_html {

		static_handler, err := http.StaticHandler()

		if err != nil {
			log.Fatal(err)
		}

		id_func := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

			path := req.URL.Path
			ext := filepath.Ext(path)

			if ext == ".geojson" && (*enable_data || *enable_geojson) {
				geojson_handler.ServeHTTP(rsp, req)
			} else if ext == ".spr" && (*enable_data || *enable_spr) {
				spr_handler.ServeHTTP(rsp, req)
			} else if ext == ".png" && (*enable_png || *enable_graphics) {
				png_handler.ServeHTTP(rsp, req)
			} else if ext == ".svg" && (*enable_svg || *enable_graphics) {
				svg_handler.ServeHTTP(rsp, req)
			} else {
				mapzenjs_handler.ServeHTTP(rsp, req)
			}

			return
		}

		id_handler := gohttp.HandlerFunc(id_func)

		mux.Handle("/id/", id_handler)
		mux.Handle("/fonts/", static_handler)
		mux.Handle("/javascript/", static_handler)
		mux.Handle("/css/", static_handler)

		err = mapzenjs_utils.AppendMapzenJSAssets(mux)

		if err != nil {
			log.Fatal(err)
		}
	}

	address := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening on %s\n", address)

	err = gohttp.ListenAndServe(address, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
