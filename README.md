# go-whosonfirst-static

![](docs/images/wof-static-sf.png)

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### wof-staticd

```
./bin/wof-staticd -h
Usage of ./bin/wof-staticd:
  -cache string
    	The named source to use for caching requests. (default "null")
  -cache-arg value
    	(0) or more user-defined '{KEY}={VALUE}' arguments to pass to the caching layer
  -config string
    	Read some or all flags from an ini-style config file. Values in the config file take precedence over command line flags.
  -data-endpoint string
    	The endpoint your HTML handler should fetch data files from.
  -debug
    	Enable debugging.
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
  -httptest.serve string
    	if non-empty, httptest.NewServer serves on this address and blocks
  -nextzen-api-key string
    	A valid Nextzen API key (https://developers.nextzen.org/). (default "xxxxxxx")
  -path-geojson string
    	The path that GeoJSON requests should be served from (default "/geojson/")
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
  -section string
    	A valid ini-style config file section. (default "wof-static")
  -source string
    	Valid sources are: fs, http, mysql, s3, sqlite (default "fs")
  -source-dsn string
    	A valid DSN string specific to the source you've chosen.
  -test-reader string
    	Perform some basic sanity checking on the reader at startup
```

## Example

```
bin/wof-staticd -port 8080 -source fs -source-dsn /usr/local/data/whosonfirst-data/data -cache lru -cache-arg 'CacheSize=500' -debug -nextzen-api-key ${NEXTZEN_APIKEY}
```

Or with a fixed-sized cache:

```
bin/wof-staticd -port 8080 -source fs -fs-root /usr/local/data/whosonfirst-data/data -cache bigcache -cache-arg HardMaxCacheSize=100 -cache-arg MaxEntrySize=1024 -debug -nextzen-api-key ${NEXTZEN_APIKEY}
2017/12/29 18:37:10 listening on localhost:8080
2017/12/29 18:37:54 REQUEST /id/85633793/
2017/12/29 18:37:54 GET 856/337/93/85633793.geojson CACHE MISS
2017/12/29 18:37:54 MISS 856/337/93/85633793.geojson
2017/12/29 18:37:54 READ 856/337/93/85633793.geojson <nil>
2017/12/29 18:37:56 SET 856/337/93/85633793.geojson entry is bigger than max shard size
```

_`HardMaxCacheSize` is measured in MB and `MaxEntrySize` in bytes._

## Lambda

Yes, it is possible to run `wof-staticd` as an AWS Lambda function.

To create the Lambda function you're going to upload to AWS simply use the handy `lambda` target in the Makefile. This will create a file called `deployment.zip` which you will need to upload to AWS (those details are out of scope for this document).

Your `wof-staticd` function should be configured with (Lambda) environment variables. Environment variables map to the standard command line flags as follows:

* The command line flag is upper-cased
* All instances of `-` are replaced with `_`
* Each flag is prefixed with `WOF_STATICD`

For example the command line flag `-protocol` would be mapped to the `WOF_STATICD_PROTOCOL` environment variable. Which is a good example because it is the one environment variable you _must_ to specify for `wof-staticd` to work as a Lambda function. Specifically you need to define the protocol as... "lambda". For example

```
WOF_STATICD_PROTOCOL = lambda
```

In reality you'll need to specify other flags, like `WOF_STATICD_SOURCE` and `WOF_STATICD_SOURCE_DSN`. For example here's how you might configure your function to render all the data and graphics formats (but not static HTML webpages) for your data:

```
WOF_STATICD_SOURCE = s3
WOF_STATICD_SOURCE_DSN = bucket={BUCKET} prefix={PREFIX} region={REGION} credentials=iam:
WOF_STATICD_ENABLE_HTML = false
WOF_STATICD_ENABLE_GRAPHICS = true
WOF_STATIC_ENABLE_DATA = true
```

### Lambda, API Gateway and images

In order for requests to produce PNG output (rather than a base64 encoded string) you will need to do a few things. Even then it's not clear that it will work and I'm uncertain whether it's AWS itself, the way AWS is configure or this code. This is what you're _supposed_ to do and... sometimes it works?

1. Make sure your API Gateway settings list `image/png` as a known and valid binary type:

![](docs/images/20180625-agw-binary.png)

2. If you've put a CloudFront distribution in front of your API Gateway then you
will to ensure that you blanket enable all HTTP headers or whitelist the
`Accept:` header , via the `Cache Based on Selected Request Headers` option (for
the CloudFront behaviour that points to your gateway):

![](docs/images/20180625-cf-cache.png)

3. Make sure you pass an `Accept: image/png` header when you request the PNG rendering.

## Docker

[Yes](Dockerfile). For example:

First, do the usual Docker `build` dance:

```
docker build -t wof-staticd .
```

Then:

```
docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='http' -e SOURCE_DSN='https://data.whosonfirst.org/' -e NEXTZEN_APIKEY=****' wof-staticd
```

Or:

```
docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='s3' -e SOURCE_DSN='bucket=whosonfirst region=us-east-1 credentials=env:' -e AWS_ACCESS_KEY_ID='***' -e AWS_SECRET_ACCESS_KEY='***' -e NEXTZEN_APIKEY='***' wof-staticd
```

Or even still, with caching:

```
docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='http' -e SOURCE_DSN='https://data.whosonfirst.org/' -e CACHE='gocache' -e CACHE_ARGS='DefaultExpiration=300 CleanupInterval=600' -e DEBUG='debug' -e NEXTZEN_APIKEY='*****' wof-static
```

## Nextzen

You will need a [valid Nextzen API key](https://developers.nextzen.org/) in order for map tiles to work.

## See also

* https://github.com/whosonfirst/go-whosonfirst-readwrite
