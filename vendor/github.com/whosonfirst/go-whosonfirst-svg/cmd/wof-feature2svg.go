package main

import (
	"flag"
	_ "github.com/facebookgo/atomicfile"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-svg"
	"log"
)

func main() {

	var width = flag.Float64("width", 1024., "...")
	var height = flag.Float64("height", 1024., "...")

	flag.Parse()

	o := svg.NewDefaultOptions()

	o.Width = *width
	o.Height = *height

	for _, path := range flag.Args() {

		f, err := feature.LoadFeatureFromFile(path)

		if err != nil {
			log.Fatal(err)
		}

		err = svg.FeatureToSVG(f, o)

		if err != nil {
			log.Fatal(err)
		}
	}
}
