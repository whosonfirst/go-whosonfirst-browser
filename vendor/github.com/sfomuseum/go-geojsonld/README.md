# go-geojson-ld

Go package for converting GeoJSON `Feature` records in to GeoJSON-LD records.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-geojson-ld.svg)](https://pkg.go.dev/github.com/sfomuseum/go-geojson-ld)

## Example

```
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-geojsonld"
	"os"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	for _, path := range flag.Args() {
	
		fh, _ := os.Open(path)
		
		body, _ := geojsonld.AsGeoJSONLDWithReader(ctx, fh)
		fmt.Println(string(body))
	}
}
```

## See also

* http://geojson.org/geojson-ld/
* https://json-ld.org/