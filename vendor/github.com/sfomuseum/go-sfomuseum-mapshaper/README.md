# go-sfomuseum-mapshaper

Go package for interacting with the mapserver-cli tool.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-sfomuseum-mapshaper.svg)](https://pkg.go.dev/github.com/sfomuseum/go-sfomuseum-mapshaper)

Documentation is incomplete.

## Tools

### server

A simple HTTP server to expose the mapserver-cli tool. Currently, only the '-points inner' functionality is exposed.

```
$> ./bin/server -h
A simple HTTP server to expose the mapserver-cli tool. Currently, only the '-points inner' functionality is exposed.
Usage:
	 ./bin/server [options]

Valid options are:
  -allowed-origins string
    	A comma-separated list of hosts to allow CORS requests from.
  -enable-cors
    	Enable support for CORS headers
  -mapshaper-path string
    	The path to your mapshaper binary. (default "/usr/local/bin/mapshaper")
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
  -uploads-max-bytes int
    	The maximum allowed size (in bytes) for uploads. (default 1048576)
```

## Docker

```
$> docker build -t mapshaper-server .

$> docker run -it -p 8080:8080 -e MAPSHAPER_SERVER_URI=http://0.0.0.0:8080 mapshaper-server

$> curl -s http://localhost:8080/api/innerpoint \
	-d @/usr/local/data/sfomuseum-data-architecture/data/174/588/208/3/1745882083.geojson \

| jq '.features[].geometry'

{
  "type": "Point",
  "coordinates": [
    -122.38875600604932,
    37.61459515528007
  ]
}
```

## See also

* https://github.com/mbloch/mapshaper