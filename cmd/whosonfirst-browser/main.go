package main

import (
	"context"
	_ "github.com/whosonfirst/go-reader-http"
	"github.com/whosonfirst/go-whosonfirst-browser"
	_ "github.com/whosonfirst/go-writer"
	"log"
)

func main() {

	ctx := context.Background()
	err := browser.Start(ctx)

	if err != nil {
		log.Fatal(err)
	}

}
