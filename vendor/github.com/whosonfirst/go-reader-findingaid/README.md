# go-reader-findingaid

Go package implementing the whosonfirst/go-reader interface for use with Who's On First "finding aids".

## Important

This package targets "version 2" of the `whosonfirst/go-whosonfirst-findingaid` package:

https://github.com/whosonfirst/go-whosonfirst-findingaid/

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-reader-findingaid.svg)](https://pkg.go.dev/github.com/whosonfirst/go-reader-findingaid)

## Example

_Error handling omitted for brevity._

```
package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-reader"
	_ "github.com/whosonfirst/go-reader-findingaid"
	"io"
	"os"
)

func main() {

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader-findingaid URI.")

	flag.Parse()

	ctx := context.Background()

	r, _ := reader.NewReader(ctx, *reader_uri)

	for _, path := range flag.Args() {

		fh, _ := r.Read(ctx, path)
		defer fh.Close()

		io.Copy(os.Stdout, fh)
	}
}
```

For a complete working example see [cmd/read](cmd/read/main.go).

## Finding aids

A Who's On First finding aid is meant to map a given ID to its corresponding `whosonfirst-data/whosonfirst-data` repository (although other sources are possible through the use of URI templates described below).

One use case for finding aids is something like the [go-whosonfirst-browser](https://github.com/whosonfirst/go-whosonfirst-browser) which doesn't have a database of IDs but instead uses one or more [go-reader](https://github.com/whosonfirst/go-reader) instances to retrieve records. That is: The `go-whosonfirst-browser` doesn't actually know anything about _where_ the data is coming from. It lets the "reader" handle all those details.

One goal with the finding aids has been to create a "finding aid reader" that when given an ID would look up its corresponding repository and fetch the data over the wire from GitHub. That way the `go-whosonfirst-browser` could run with a minimal footprint (read: No database with a bazillion WOF records).

The finding aid model has two "tables". One is to store the WOF ID lookup and looks like this:

```
whosonfirst_id, repo_id
```

And one to store the repo ID and it's corresponding name:

```
repo_name, repo_id
```

The idea being that storing string repo names for every record is a waste of space and processing time. Although it may probably be the case that any given finding aids will map to a single WOF repo it is possible for a finding aid to contain pointers to records from multiple repositories.

This package _does not produce Who's On First finding aids_, it only consumes them. For information and tools to produce finding aids please consult:

* https://github.com/whosonfirst/go-whosonfirst-findingaid/tree/v2

## Finding aid URIs

Finding aid URIs take the form of:

```
findingaid://{RESOLVER}?{QUERY_PARAMETERS}
```

Valid finding aid query parameters are:

| Name | Type | Notes | Required
| --- | --- | --- | --- |
| dsn | string | A valid `matt/go-sqlite3` DSN string | yes |
| template | string | A valid `jtacoma/uritemplates` for resolving final reader URIs. If empty the default URI template mapping to `whosonfirst-data/whosonfirst-data` repositories will be used. | no |

For example:

```
findingaid://sqlite?dsn=/usr/local/data/findingaids/wof.db
```

Note: Although it is possible to produce Who's On First finding aids that are not SQLite database this package _only_ works with SQLite-based finding aids.

## How does it work?

First, given a URI (for example `174/616/026/9/1746160269.geojson` or `1746160269.geojson` or `1746160269`) that URI is resolved to its Who's On First ID: `1746160269`

Second, the repository name associated with that ID (for example `sfomuseum-data-maps`) is retrieved from the finding aid database you specified in the reader's constructor.

Third, the Who's On First ID derived in step (1) is then expanded in to it nested relative URL (for example `174/616/026/9/1746160269.geojson`).

Fourth, the repository name derived in step (2) is used to expand a URI template used to create a new internal `go-reader.Reader` instance. The default URI template points the `https://github.com/whosonfirst-data` repositories but you can define your own. For example:

```
findingaid_uri := "findingaid://sqlite?dsn=fa.db&template=https://raw.githubusercontent.com/sfomuseum-data/{repo}/main/data/"

ctx := context.Background()
r, _ := reader.NewReader(ctx, findingaid_uri)
r.Read(ctx, 1746160269)
```

In the example the final document would be read from https://raw.githubusercontent.com/sfomuseum-data/sfomuseum-data-maps/main/data/174/616/026/9/1746160269.geojson.

This means that URI template can define a constructor for [any go-reader.Reader implementation](https://github.com/whosonfirst/go-reader#custom-readers) so long as you've imported that package in your code.

## Tools

```
$> make cli
go build -mod vendor -o bin/read cmd/read/main.go
```

### read

Resolve one or more URIs, using a Who's On First finding aid and read their corresponding Who's On First documents, outputting each to `STDOUT`.

```
> ./bin/read -h
Resolve one or more URIs, using a Who's On First finding aid and read their corresponding Who's On First documents, outputting each to STDOUT.
Usage:
	 ./bin/read uri(N) uri(N)
  -reader-uri string
    	A valid whosonfirst/go-reader-findingaid URI
```

For example:

```
$> ./bin/read \
	-reader-uri 'findingaid://sqlite?dsn=/usr/local/data/findingaids/wof.db' \
	102527513 \
	
| jq '.["properties"]["wof:name"]'

"San Francisco International Airport"
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-reader
* https://github.com/whosonfirst/go-whosonfirst-reader-http
* https://github.com/whosonfirst/go-whosonfirst-findingaid
