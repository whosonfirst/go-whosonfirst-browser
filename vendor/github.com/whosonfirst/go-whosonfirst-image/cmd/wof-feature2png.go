package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-image"
	"github.com/whosonfirst/warning"
	"log"
	"os"
)

func main() {

	var width = flag.Int("width", 1024, "...")
	var height = flag.Int("height", 1024, "...")

	flag.Parse()

	opts := image.NewDefaultOptions()

	opts.Writer = os.Stdout // this is redundant but whatever
	opts.Width = *width
	opts.Height = *height

	for _, path := range flag.Args() {

		f, err := feature.LoadFeatureFromFile(path)

		if err != nil && !warning.IsWarning(err) {
			log.Fatal(err)
		}

		err = image.FeatureToPNG(f, opts)

		if err != nil {
			log.Fatal(err)
		}
	}
}
