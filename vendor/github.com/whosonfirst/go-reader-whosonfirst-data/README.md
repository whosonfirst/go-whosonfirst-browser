# go-reader-whosonfirst-data

Who's On First (GitHub) data support for the go-reader.Reader interface.

## Example

```
package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-reader"
	_ "github.com/whosonfirst/go-reader-whosonfirst-data"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"os"
)

func main() {

	reader_uri := flag.String("reader-uri", "whosonfirst-data://", "A valid go-reader-whosonfirst-data URI.")
	
	flag.Parse()

	ctx := context.Background()
	r, _ := reader.NewReader(ctx, *reader_uri)

	uris := flag.Args()

	for _, raw := range uris {

		id, _, _ := uri.ParseURI(raw)
		rel_path, _ := uri.Id2RelPath(id)

		fh, _ := r.Read(ctx, rel_path)
		defer fh.Close()

		io.Copy(os.Stdout, fh)
	}
}
```

_Error handling omitted for the sake of brevity._

### Notes

The default behaviour for resolving a Who's On First ID to its corresponding GitHub repository is to use the [data.whosonfirst.org/findingaid](https://data.whosonfirst.org/findingaid/) lookup service. If you know the repository that all of the pathN arguments are part of you can skip the lookup service by appending the repo name as a query parameter to the `-reader-uri` argument.

For example:

```
-reader-uri whosonfirst-data://?repo=whosonfirst-data-admin-ca
```

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/read cmd/read/main.go
```

### read

```
$> ./bin/read -h
Read one or more Who's On First records from the whosonfirst-data GitHub organization.

Usage:
  ./bin/read [options] [path1 path2 ... pathN]

Options:
  -reader-uri string
    	A valid go-reader-whosonfirst-data URI. (default "whosonfirst-data://")

Notes:

pathN may be any valid Who's On First ID or URI that can be parsed by the
go-whosonfirst-uri package.

The default behaviour for resolving a Who's On First ID to its corresponding
GitHub repository is to use the 'data.whosonfirst.org/findingaid' lookup
service. If you know the repository that all of the pathN arguments are part of
you can skip the lookup service by appending the repo name as a query parameter
to the '-reader-uri' argument. For example: -reader-uri
'whosonfirst-data://?repo=whosonfirst-data-admin-ca'.
```

Read one or more Who's On First records from the whosonfirst-data GitHub organization.

For example:

```
$> ./bin/read 85633041 | jq '.properties["wof:name"]'
"Canada"
```

## See also

* https://github.com/whosonfirst-data
* https://github.com/go-whosonfirst-reader
* https://github.com/go-whosonfirst-reader-github
* https://github.com/go-whosonfirst-findingaid
* https://github.com/go-whosonfirst-whosonfirst-reader