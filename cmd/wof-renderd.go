package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-http-mapzenjs"
	"github.com/whosonfirst/go-whosonfirst-render/http"
	"github.com/whosonfirst/go-whosonfirst-render/reader"
	"log"
	gohttp "net/http"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var source = flag.String("source", "fs", "...")
	var root = flag.String("fs-root", "", "...")

	var s3_bucket = flag.String("s3-bucket", "whosonfirst.mapzen.com", "...")
	var s3_prefix = flag.String("s3-prefix", "", "...")
	var s3_region = flag.String("s3-region", "us-east-1", "...")
	var s3_creds = flag.String("s3-credentials", "", "...")

	var api_key = flag.String("mapzen-apikey", "mapzen-xxxxxxx", "")

	flag.Parse()

	var r reader.Reader
	var err error

	switch *source {
	case "fs":
		r, err = reader.NewFSReader(*root)
	case "s3":

		cfg := reader.S3Config{
			Bucket:      *s3_bucket,
			Prefix:      *s3_prefix,
			Region:      *s3_region,
			Credentials: *s3_creds,
		}

		r, err = reader.NewS3Reader(cfg)
	default:
		err = errors.New("Unknown or invalid source")
	}

	if err != nil {
		log.Fatal(err)
	}

	handlers := make(map[string]gohttp.Handler)

	html_opts := http.NewDefaultHTMLOptions()
	html_opts.MapzenAPIKey = *api_key

	html_handler, err := http.HTMLHandler(r, html_opts)

	if err != nil {
		log.Fatal(err)
	}

	handlers["/"] = html_handler

	ping_handler, err := http.PingHandler()

	if err != nil {
		log.Fatal(err)
	}

	handlers["/ping"] = ping_handler

	static_handler, err := http.StaticHandler()

	if err != nil {
		log.Fatal(err)
	}

	mapzenjs_handler, err := mapzenjs.MapzenJSHandler()

	if err != nil {
		log.Fatal(err)
	}

	handlers["/javascript/mapzen.min.js"] = mapzenjs_handler
	handlers["/javascript/tangram.min.js"] = mapzenjs_handler
	handlers["/javascript/mapzen.js"] = mapzenjs_handler
	handlers["/javascript/tangram.js"] = mapzenjs_handler
	handlers["/css/mapzen.js.css"] = mapzenjs_handler
	handlers["/tangram/refill-style.zip"] = mapzenjs_handler

	handlers["/javascript/slippymap.crosshairs.js"] = static_handler
	handlers["/javascript/whosonfirst.spr.js"] = static_handler
	handlers["/css/whosonfirst.spr.css"] = static_handler

	mux := gohttp.NewServeMux()

	for uri, handler := range handlers {
		mux.Handle(uri, handler)
	}

	address := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening on %s\n", address)

	err = gohttp.ListenAndServe(address, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
