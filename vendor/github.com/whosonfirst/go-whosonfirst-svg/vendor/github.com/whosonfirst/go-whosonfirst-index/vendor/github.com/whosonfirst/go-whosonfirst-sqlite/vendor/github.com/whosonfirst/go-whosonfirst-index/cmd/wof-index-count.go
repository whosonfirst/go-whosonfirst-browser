package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-index"
	"io"
	"log"
	"runtime"
	"sync/atomic"
)

func main() {

	var mode = flag.String("mode", "repo", "")
	var procs = flag.Int("processes", runtime.NumCPU()*2, "")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	var count int64
	count = 0

	f := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		i := atomic.AddInt64(&count, 1)

		log.Println("PATH", path, i)
		return nil
	}

	i, err := index.NewIndexer(*mode, f)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		err := i.IndexPath(path)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println(count, i.Indexed)
}
