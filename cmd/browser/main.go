package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/http"
	"github.com/whosonfirst/go-whosonfirst-browser/server"
	"log"
	gohttp "net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func main() {

	var proto = flag.String("protocol", "http", "The protocol for wof-staticd server to listen on. Valid protocols are: http, lambda.")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var reader_source = flag.String("reader-source", "", "...")
	var cache_source = flag.String("cache-source", "", "...")

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

	var path_png = flag.String("path-png", "/png/", "The path that PNG requests should be served from")
	var path_svg = flag.String("path-svg", "/svg/", "The path that PNG requests should be served from")
	var path_geojson = flag.String("path-geojson", "/geojson/", "The path that GeoJSON requests should be served from")
	var path_spr = flag.String("path-spr", "/spr/", "The path that SPR requests should be served from")

	flag.Parse()

	r, err := reader.NewReader(*reader_source)

	if err != nil {
		log.Fatal(err)
	}

	c, err := cache.NewCache(*cache_source)

	if err != nil {
		log.Fatal(err)
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

	bootstrap_opts := bootstrap.DefaultBootstrapOptions()

	tangramjs_opts := tangramjs.DefaultTangramJSOptions()
	tangramjs_opts.Nextzen.APIKey = *nextzen_apikey
	tangramjs_opts.Nextzen.StyleURL = *nextzen_style_url
	tangramjs_opts.Nextzen.TileURL = *nextzen_tile_url

	err = bootstrap.AppendAssetHandlersWithPrefix(mux, *static_prefix)

	if err != nil {
		log.Fatal(err)
	}

	html_handler = bootstrap.AppendResourcesHandlerWithPrefix(html_handler, bootstrap_opts, *static_prefix)
	html_handler = tangramjs.AppendResourcesHandlerWithPrefix(html_handler, tangramjs_opts, *static_prefix)

	err = tangramjs.AppendAssetHandlersWithPrefix(mux, *static_prefix)

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
		mux.Handle(*path_png, png_handler)
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
		mux.Handle(*path_svg, svg_handler)
	}

	if *enable_all || *enable_data || *enable_spr {

		h, err := http.SPRHandler(cr)

		if err != nil {
			log.Fatal(err)
		}

		spr_handler = h
		mux.Handle(*path_spr, spr_handler)
	}

	if *enable_all || *enable_data || *enable_geojson {

		h, err := http.GeoJSONHandler(cr)

		if err != nil {
			log.Fatal(err)
		}

		geojson_handler = h
		mux.Handle(*path_geojson, geojson_handler)
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

	address := fmt.Sprintf("http://%s:%d", *host, *port)

	u, err := url.Parse(address)

	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewStaticServer(*proto, u)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
