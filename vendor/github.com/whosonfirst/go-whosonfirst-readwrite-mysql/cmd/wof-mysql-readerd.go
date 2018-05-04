package main

import (
	"flag"
	"fmt"
	mysql_reader "github.com/whosonfirst/go-whosonfirst-readwrite-mysql/reader"
	"github.com/whosonfirst/go-whosonfirst-readwrite/http"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"log"
	gohttp "net/http"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")
	var table = flag.String("table", "geojson", "The name of the MySQL table (indexed by go-whosonfirst-mysql) to query")

	var dsn = flag.String("dsn", "", "")

	flag.Parse()

	var r reader.Reader

	if *table == "geojson" {

		gr, err := mysql_reader.NewMySQLGeoJSONReader(*dsn)

		if err != nil {
			log.Fatal(err)
		}

		r = gr
	} else {

		wr, err := mysql_reader.NewMySQLWhosonfirstReader(*dsn)

		if err != nil {
			log.Fatal(err)
		}

		r = wr

	}

	read_handler, err := http.ReadHandler(r)

	if err != nil {
		log.Fatal(err)
	}

	mux := gohttp.NewServeMux()
	mux.Handle("/", read_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening for requests on %s\n", endpoint)

	err = gohttp.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
