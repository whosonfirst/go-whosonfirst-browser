# go-whosonfirst-browser

![](docs/images/wof-browser-montreal.png)

Go package for rendering known Who's On First (WOF) IDs in a number of formats.

_This package used to be called `go-whosonfirst-static`. Now it is called `go-whosonfirst-browser.`_

## Install

You will need to have the `Go` programming language (specifically version [1.12](https://golang.org/dl/) or higher) installed. All of this package's dependencies are bundled with the code in the `vendor` directory.

## Things this package is not

### This is not a replacement for the [Who's On First Spelunker](https://spelunker.whosonfirst.org/).

At least not yet.

`go-whosonfirst-browser` was designed to be a simple display tool for known Who's On First (WOF) IDs and records. That constitutes a third to half of [what the Spelunker does](https://github.com/whosonfirst/whosonfirst-www-spelunker) (the remainder being list views and facets) so in principle it would be easy enough to add the same functionality here.

The principle advantage of migrating Spelunker functionality to this package is that it does not have any external dependencies and has been support for multiple data sources and caches and can be pre-compiled in to a standalone binary tool. The principle disadvantage would be that experimenting and developing code and functionality in Python (used by the existing Spelunker) has a lower barrier to entry than doing the same in Go (used by this package).

### This is not a search engine.

This is a tool that is primarily geared towards displaying _known_ Who's On First IDs. It does not maintain an index, or a list of known reocrds, before it displays them.

It would be easy enough to add flags to use an external instance of the [Pelias Placeholder API](https://millsfield.sfomuseum.org/blog/2019/11/04/placeholder/) for basic search functionality so we'll add that to the list of features for a "2.x" release.

### This is not a tool for editing Who's On First documents.

At least not yet.

Interestingly the code that renders Who's On First (WOF) property dictionaries in to pretty HTML tables is the same code used for the experimental Mapzen "[Yes No Fix](https://whosonfirst.org/blog/2016/04/08/yesnofix/) project". That functionality has not been enabled or tested with this tool yet.

## Tools

### whosonfirst-browser

To build the browser use the handy `go build` tool, like this:

```
go build -mod vendor cmd/whosonfirst-browser/main.go bin/whosonfirst-browser
```

```
$> bin/whosonfirst-browser -h
Usage of ./bin/whosonfirst-browser:
  -cache-source string
    	A valid go-cache Cache URI string. (default "gocache://")
  -enable-all
    	Enable all the available output handlers.
  -enable-data
    	Enable the 'geojson' and 'spr' output handlers.
  -enable-geojson
    	Enable the 'geojson' output handler. (default true)
  -enable-graphics
    	Enable the 'png' and 'svg' output handlers.
  -enable-html
    	Enable the 'html' (or human-friendly) output handlers. (default true)
  -enable-png
    	Enable the 'png' output handler.
  -enable-spr
    	Enable the 'spr' (or "standard places response") output handler. (default true)
  -enable-svg
    	Enable the 'svg' output handler.
  -host string
    	The hostname to listen for requests on (default "localhost")
  -nextzen-api-key string
    	A valid Nextzen API key (https://developers.nextzen.org/). (default "xxxxxxx")
  -nextzen-style-url string
    	A valid Tangram scene file URL. (default "/tangram/refill-style.zip")
  -nextzen-tile-url string
    	A valid Nextzen MVT tile URL. (default "https://{s}.tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt")
  -path-geojson string
    	The path that GeoJSON requests should be served from. (default "/geojson/")
  -path-id string
    	The that Who's On First documents should be served from. (default "/id/")
  -path-png string
    	The path that PNG requests should be served from. (default "/png/")
  -path-spr string
    	The path that SPR requests should be served from. (default "/spr/")
  -path-svg string
    	The path that SVG requests should be served from. (default "/svg/")
  -port int
    	The port number to listen for requests on (default 8080)
  -protocol string
    	The protocol for wof-staticd server to listen on. Valid protocols are: http, lambda. (default "http")
  -proxy-tiles
    	Proxy (and cache) Nextzen tiles.
  -proxy-tiles-cache string
    	A valid tile proxy DSN string. (default "gocache://")
  -proxy-tiles-timeout int
    	The maximum number of seconds to allow for fetching a tile from the proxy. (default 30)
  -proxy-tiles-url string
    	The URL (a relative path) for proxied tiles. (default "/tiles/")
  -data-source string
    	A valid go-reader Reader URI string. (default "https://data.whosonfirst.org")
  -static-prefix string
    	Prepend this prefix to URLs for static assets.
  -templates string
    	An optional string for local templates. This is anything that can be read by the 'templates.ParseGlob' method.
```

#### Example

```
$> bin/browser -enable-all -nextzen-api-key ${NEXTZEN_APIKEY}
2019/12/14 18:22:16 Listening on http://localhost:8080
```

## Output formats

The following output formats are available.

### GeoJSON

A raw Who's On First (WOF) GeoJSON document. For example:

![](docs/images/wof-browser-montreal-geojson.png)

`http://localhost:8080/geojson/101736545`

### HTML

A responsive HTML table and map for a given WOF ID. For example:

![](docs/images/wof-browser-montreal-html.png)

`http://localhost:8080/id/101736545`

### PNG

A PNG-encoded representation of the geometry for a given WOF ID. For example:

![](docs/images/wof-browser-montreal-png.png)

`http://localhost:8080/png/101736545`

### SPR (Standard Places Response)

A JSON-encoded "standard places response" for a given WOF ID. For example:

![](docs/images/wof-browser-montreal-spr.png)

`http://localhost:8080/spr/101736545`

### SVG

An XML-encoded SVG representation of the geometry for a given WOF ID.  For example:

![](docs/images/wof-browser-montreal-svg.png)

`http://localhost:8080/svg/101736545`

## Tiles

`go-whosonfirst-browser` uses [Nextzen](https://nextzen.org/) vector data tiles and the [Tangram.js](https://github.com/tangrams/tangram) rendering library for displaying maps. The Tangram code and styling assets are bundled with this tool and served directly but, by default, tile data is retrieved from the Nextzen servers.

It is possible to cache those tiles locally using the `-proxy-tiles` flag at start up. The default cache for proxying tiles is an ephemiral in-memory cache but you can also specify an alternative [go-cache](https://github.com/whosonfirst/go-cache) `cache.Cache` source using the `-proxy-tiles-cache` flag. Caches are discussed in detail below.

## Data sources and Caches

The `go-whosonfirst-browser` uses the [go-reader](#) `reader.Reader` and [go-cache](#) `cache.Cache` interfaces for reading and caching data respectively. This enables the "guts" of the code to be developed and operate independently of any individual data source or cache.

The default `whosonfirst-browser` tool allows data sources to be specified as a localfile system or a remote HTTP(S) endpoint and caching sources as a local filesystem or an ephemiral in-memory lookup.

This is what the code for default `whosonfirst-browser` tool looks like, with error handling omitted for the sake of brevity:

```
package main

import (
	"context"
	_ "github.com/whosonfirst/go-reader-http"
	"github.com/whosonfirst/go-whosonfirst-browser"
)

func main() {
	ctx := context.Background()
	browser.Start(ctx)
}
```

For example... :

```
$> bin/whosonfirst-browser -enable-all \
	-data-source 'file:///usr/local/data/whosonfirst-data-admin-us/data' \
	-nextzen-api-key ${NEXTZEN_APIKEY}	
```

Or... :

```
$> bin/whosonfirst-browser -enable-all \
	-data-source 'https://data.whosonfirst.org' \
	-cache-source 'file:///usr/local/cache/whosonfirst' \
	-nextzen-api-key ${NEXTZEN_APIKEY}	
```


But if you wanted to ... using the [Go Cloud Blob package](https://gocloud.dev/howto/blob/)

```
package main

import (
	"context"
	_ "github.com/whosonfirst/go-reader-blob"
	"github.com/whosonfirst/go-whosonfirst-browser"
)

func main() {
	ctx := context.Background()
	browser.Start(ctx)
}
```

And then you would start the `browser` tool like this:

```
$> bin/browser -enable-all \
	-data-source 's3://{BUCKET}?region={REGION}' \
	-nextzen-api-key ${NEXTZEN_APIKEY}
```

Or if you wanted to ... and cache tiles locally to a local directory:

```
package main

import (
	"context"
	_ "github.com/whosonfirst/go-cache-blob"	
	_ "github.com/whosonfirst/go-reader-blob"
	"github.com/whosonfirst/go-whosonfirst-browser"
	_ "gocloud.dev/blob/fileblob"
)

func main() {
	ctx := context.Background()
	browser.Start(ctx)
}
```

And then you would start the `browser` tool like this:

```
$> bin/browser -enable-all \
	-proxy-tiles \
	-proxy-tiles-cache 'file:///usr/local/data/tilezen'
	-data-source 's3://{BUCKET}?region={REGION}' \
	-nextzen-api-key ${NEXTZEN_APIKEY}
```

## Lambda

Yes, it is possible to run `browser` as an AWS Lambda function.

To create the Lambda function you're going to upload to AWS simply use the handy `lambda` target in the Makefile. This will create a file called `deployment.zip` which you will need to upload to AWS (those details are out of scope for this document).

Your `wof-staticd` function should be configured with (Lambda) environment variables. Environment variables map to the standard command line flags as follows:

* The command line flag is upper-cased
* All instances of `-` are replaced with `_`
* Each flag is prefixed with `BROWSER`

For example the command line flag `-protocol` would be mapped to the `BROWSER_PROTOCOL` environment variable. Which is a good example because it is the one environment variable you _must_ to specify for `wof-staticd` to work as a Lambda function. Specifically you need to define the protocol as... "lambda". For example

```
BROWSER_PROTOCOL = lambda
```

In reality you'll need to specify other flags, like `BROWSER_SOURCE` and `BROWSER_SOURCE_DSN`. For example here's how you might configure your function to render all the data and graphics formats (but not static HTML webpages) for your data:

```
BROWSER_CACHE_SOURCE = gocache://
BROWSER_READER_SOURCE = s3://{BUCKET}?prefix={PREFIX}&region={REGION}
BROWSER_ENABLE_HTML = false
BROWSER_ENABLE_GRAPHICS = true
BROWSER_ENABLE_DATA = true
```

### Lambda, API Gateway and images

In order for requests to produce PNG output (rather than a base64 encoded string) you will need to do a few things.

1. Make sure your API Gateway settings list `image/png` as a known and valid binary type:

![](docs/images/20180625-agw-binary.png)

2. If you've put a CloudFront distribution in front of your API Gateway then you
will to ensure that you blanket enable all HTTP headers or whitelist the
`Accept:` header , via the `Cache Based on Selected Request Headers` option (for
the CloudFront behaviour that points to your gateway):

2a. **Or:** Don't use a custom whitelist (in your behaviour settings) but make sure you pass a custom header in your origin settings (see `3a` for details).

![](docs/images/20180625-cf-cache.png)

3. Make sure you pass an `Accept: image/png` header when you request the PNG rendering.

3a. **Or:** make sure you specify a `Origin Custom Headers` header in your CloudFront origin settings (specifically `Accept: image/png`)

4. If you add another image (or binary) handler to this package you'll need to
repeat steps 1-3 _and_ update the `BinaryContentTypes` dictionary in
[server/lambda.go](server/lambda.go) code accordingly. Good times...

## Docker

[Yes](Dockerfile).

## Nextzen

You will need a [valid Nextzen API key](https://developers.nextzen.org/) in order for map tiles to work.

## See also

* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-cache
* https://github.com/aaronland/go-http-bootstrap
* https://github.com/aaronland/go-http-tangramjs
* https://github.com/sfomuseum/go-http-tilezen