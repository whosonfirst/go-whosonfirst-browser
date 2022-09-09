// whosonfirst-browser is a command line tool that launches a web application for browsing Who's On First -style records.
// For example:
//
//	$> bin/whosonfirst-browser -enable-all -nextzen-tilepack-database /usr/local/data/nextzen-world-2019-1-10.db \
//		-reader-uri 'findingaid://awsdynamodb/findingaid?region=local&endpoint=http://localhost:8000&credentials=static:local:local:local&partition_key=id&template=https://raw.githubusercontent.com/sfomuseum-data/{repo}/main/data/' \
//		-reader-uri 'findingaid://awsdynamodb/findingaid?region=local&endpoint=http://localhost:8000&credentials=static:local:local:local&partition_key=id&template=https://raw.githubusercontent.com/sfomuseum-data/{repo}/master/data/'
//	2021/11/06 12:12:42 Listening on http://localhost:8080
//
// Or:
//
//	$> ./bin/whosonfirst-browser -reader-uri 'findingaid://awsdynamodb/findingaid?partition_key=id&region=us-west-2&credentials=session&template=https://raw.githubusercontent.com/sfomuseum-data/{repo}/main/data/'
//	2022/02/08 10:38:48 Listening on http://localhost:8080
package main

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-browser/v5/application/browser"
	"log"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := browser.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run browser application, %v", err)
	}

}
