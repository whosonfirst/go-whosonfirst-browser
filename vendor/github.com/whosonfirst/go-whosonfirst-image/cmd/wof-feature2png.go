package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-image"
	"log"
	"os"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		f, err := feature.LoadFeatureFromFile(path)

		if err != nil {
			log.Fatal(err)
		}

		wr := os.Stdout

		err = image.FeatureToPNG(f, wr)

		if err != nil {
			log.Fatal(err)
		}
	}
}
