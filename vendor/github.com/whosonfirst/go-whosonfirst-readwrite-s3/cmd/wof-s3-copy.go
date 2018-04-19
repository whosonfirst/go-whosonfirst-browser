package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-readwrite-s3/config"
	"github.com/whosonfirst/go-whosonfirst-readwrite-s3/reader"
	"github.com/whosonfirst/go-whosonfirst-readwrite-s3/writer"
	"log"
)

func main() {

	var source = flag.String("source", "", "...")
	var target = flag.String("target", "", "...")

	flag.Parse()

	r_cfg, err := config.NewS3ConfigFromString(*source)

	if err != nil {
		log.Fatal(err)
	}

	r, err := reader.NewS3Reader(r_cfg)

	if err != nil {
		log.Fatal(err)
	}

	w_cfg, err := config.NewS3ConfigFromString(*target)

	if err != nil {
		log.Fatal(err)
	}

	w, err := writer.NewS3Writer(w_cfg)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		fh, err := r.Read(path)

		if err != nil {
			log.Fatal(err)
		}

		err = w.Write(path, fh)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("copied %s\n", path)
	}

}
