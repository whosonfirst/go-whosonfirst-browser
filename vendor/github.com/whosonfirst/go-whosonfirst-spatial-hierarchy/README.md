# go-whosonfirst-spatial-hierarchy

Opionated Who's On First (WOF) hierarchy for `go-whosonfirst-spatial` packages.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-spatial-hierarchy.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-spatial-hierarchy)

Documentation is incomplete.

## Example

```
import (
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
)

import (
	"github.com/whosonfirst/go-whosonfirst-spatial-hierarchy"
	hierarchy_filter "github.com/whosonfirst/go-whosonfirst-spatial-hierarchy/filter"		
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	spatial_filter "github.com/whosonfirst/go-whosonfirst-spatial/filter"
)

body := []byte(`{"type":"Feature" ...}`)

spatial_db, _ := database.NewSpatialDatabase(ctx, "sqlite://?dsn=/usr/local/data/whosonfirst.db")

resolver_opts := &hierarchy.PointInPolygonHierarchyResolverOptions{
	Database: spatial_db,
}

resolver, _ := hierarchy.NewPointInPolygonHierarchyResolver(ctx, resolver_opts)

inputs := &spatial_filter.SPRInputs{}

results_cb := hierarchy_filter.FirstButForgivingSPRResultsFunc
update_cb := hierarchy.DefaultPointInPolygonHierarchyResolverUpdateCallback()
		
new_body, _ := resolver.PointInPolygonAndUpdate(ctx, inputs, results_cb, update_cb, body)
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-pip
* https://github.com/whosonfirst/go-whosonfirst-exporter
* https://github.com/sfomuseum/go-sfomuseum-mapshaper