# go-reader-cachereader

Go package implementing the `whosonfirst/go-reader` interface for use with a caching layer.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-reader-cachereader.svg)](https://pkg.go.dev/github.com/whosonfirst/go-reader-cacheread)

## Example

```
package main

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-cachereader"	
	"io/ioutil"
	"log"
)

func main(){

	ctx := context.Background()

	reader_uri := "fs:///usr/local/data"
	cache_uri := "gocache://"

	cr_uri := fmt.Sprintf("cachereader://?reader=%s&cache=%s", reader_uri, cache_uri)
	
	r, _ := reader.NewReader(ctx, cr_uri)

	path := "101736545.geojson"
	
	for i := 0; i < 3; i++ {

		fh, _ := r.Read(ctx, path)
		defer fh.Close()

		io.Copy(ioutil.Discard, fh)

		status, _ := cachereader.GetLastRead(r, path)

		switch i {
		case 0:
			if status != cachereader.CacheMiss {
				log.Printf("Expected cache miss on first read of %s", path)
			}
		default:
			if status != cachereader.CacheHit {
				log.Printf("Expected cache hit after first read of %s", path)
			}
		}
	}
}
```

_Error handling omitted for the sake of brevity._

## See also

* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-cache