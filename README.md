# go-whosonfirst-browser

![](docs/images/wof-static-sf.png)

_This package used to be called `go-whosonfirst-static`. Now it is called `go-whosonfirst-browser.`_

## Install

You will need to have both `Go` (specifically version [1.12](https://golang.org/dl/) or higher) and the `make` programs installed on your computer.

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### browser

```
go build -mod vendor cmd/browser/main.go bin/browser
```

```
$> bin/browser -h
Usage of ./bin/browser:
  -cache-source string
    	... (default "gocache://")
  -enable-all
    	Enable all the available output handlers.
  -enable-data
    	Enable the 'geojson' and 'spr' output handlers.
  -enable-geojson
    	Enable the 'geojson' output handler. (default true)
  -enable-graphics
    	Enable the 'png' and 'svg' output handlers.
  -enable-html
    	Enable the 'html' (or human-friendly) output handler. (default true)
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
    	... (default "/tangram/refill-style.zip")
  -nextzen-tile-url string
    	... (default "https://{s}.tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt")
  -path-geojson string
    	The path that GeoJSON requests should be served from (default "/geojson/")
  -path-id string
    	... (default "/id/")
  -path-png string
    	The path that PNG requests should be served from (default "/png/")
  -path-spr string
    	The path that SPR requests should be served from (default "/spr/")
  -path-svg string
    	The path that PNG requests should be served from (default "/svg/")
  -port int
    	The port number to listen for requests on (default 8080)
  -protocol string
    	The protocol for wof-staticd server to listen on. Valid protocols are: http, lambda. (default "http")
  -reader-source string
    	... (default "https://data.whosonfirst.org")
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

## go-reader.Reader(s) and go-cache.Cache(s)

This is what the code for default `browser` tool looks like, with error handling omitted for the sake of brevity:

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

But if you wanted to using the [Go Cloud](#)

```
package main

import (
	"context"
	_ "github.com/whosonfirst/go-reader-blob"
	_ "github.com/whosonfirst/go-reader-http"	
	"github.com/whosonfirst/go-whosonfirst-browser"
)

func main() {
	ctx := context.Background()
	browser.Start(ctx)
}
```

And then you would start the `browser` tool like this:

```
$> bin/browser -reader-source 's3://{BUCKET}?region={REGION}' -enable-all -nextzen-api-key ${NEXTZEN_APIKEY}
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

There used to be. There isn't now. There will probably be one again.

## Nextzen

You will need a [valid Nextzen API key](https://developers.nextzen.org/) in order for map tiles to work.

## See also

* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-cache
* https://github.com/aaronland/go-http-bootstrap
* https://github.com/aaronland/go-http-tangramjs
