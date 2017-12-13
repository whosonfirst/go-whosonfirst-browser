package main

import (
       "errors"
	"flag"
	"fmt"
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

	html_handler, err := http.HTMLHandler(r)

	if err != nil {
		log.Fatal(err)
	}

	ping_handler, err := http.PingHandler()

	if err != nil {
		log.Fatal(err)
	}

	address := fmt.Sprintf("%s:%d", *host, *port)

	mux := gohttp.NewServeMux()
	mux.Handle("/", html_handler)
	mux.Handle("/ping", ping_handler)

	err = gohttp.ListenAndServe(address, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
