package main

import (
	"context"
	_ "github.com/whosonfirst/go-reader-http"
	_ "github.com/whosonfirst/go-reader-github"	
	"github.com/whosonfirst/go-whosonfirst-browser"
	_ "github.com/whosonfirst/go-writer"
	_ "github.com/whosonfirst/go-writer-github"		
	"log"
)

func main() {

	ctx := context.Background()
	err := browser.Start(ctx)

	if err != nil {
		log.Fatal(err)
	}

}
