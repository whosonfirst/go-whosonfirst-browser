package main

import (
	"context"
	_ "github.com/whosonfirst/go-reader-http"
	_ "github.com/whosonfirst/go-reader-whosonfirst-data"
	"github.com/whosonfirst/go-whosonfirst-browser/v3"
	"log"
)

func main() {

	ctx := context.Background()
	err := browser.Start(ctx)

	if err != nil {
		log.Fatal(err)
	}

}
