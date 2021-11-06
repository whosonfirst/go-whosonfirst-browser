package main

/*

> bin/whosonfirst-browser -enable-all -nextzen-tilepack-database /usr/local/data/nextzen-world-2019-1-10.db -reader-uri 'findingaid://awsdynamodb/findingaid?region=local&endpoint=http://localhost:8000&credentials=static:local:local:local&partition_key=id&template=https://raw.githubusercontent.com/sfomuseum-data/{repo}/main/data/' -reader-uri 'findingaid://awsdynamodb/findingaid?region=local&endpoint=http://localhost:8000&credentials=static:local:local:local&partition_key=id&template=https://raw.githubusercontent.com/sfomuseum-data/{repo}/master/data/'
2021/11/06 12:12:42 Listening on http://localhost:8080

*/

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-browser/v4/application/browser"
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
