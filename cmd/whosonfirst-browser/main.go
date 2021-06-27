package main

import (
	_ "github.com/whosonfirst/go-reader-http"
	_ "github.com/whosonfirst/go-reader-whosonfirst-data"
)

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-browser/v3/application/browser"
	"log"
)

func main() {

	ctx := context.Background()

	app, err := browser.NewBrowserApplication(ctx)

	if err != nil {
		log.Fatalf("Failed to create browser application, %v", err)
	}

	err = app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run browser application, %v", err)
	}

}
