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
    	 ... (default "null")
  -cache-arg value
    	     (0) or more user-defined '{KEY}={VALUE}' arguments to pass to the caching layer
  -data-endpoint string
    		 
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
  -port int
    	The port number to listen for requests on (default 8080)
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

## Docker

[Yes](Dockerfile). For example:

First, do the usual Docker `build` dance:

```
docker build -t wof-static .
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
