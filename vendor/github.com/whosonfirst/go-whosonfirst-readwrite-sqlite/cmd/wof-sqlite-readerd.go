package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-readwrite/http"	
	"github.com/whosonfirst/go-whosonfirst-readwrite-sqlite/reader"
	"log"
	"os"
	gohttp "net/http"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var dsn = flag.String("dsn", "", "")

	flag.Parse()

	r, err := reader.NewSQLiteReader(*dsn)

	if err != nil {
		log.Fatal(err)
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
