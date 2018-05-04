# go-whosonfirst-readwrite-mysql

This package assumes a few things:

1. You are using MySQL 5.7 or higher
2. You have indexed a `whosonfirst` or `geojson` table (or both) using the [go-whosonfirst-mysql](https://github.com/whosonfirst/go-whosonfirst-mysql) package (or equivalent code)
3. The side-effect of `2` is that this package is still Who's On First (and not arbitrary GeoJSON) specific

## Tools

### wof-mysql-readerd

```
./bin/wof-mysql-readerd -h
Usage of ./bin/wof-mysql-readerd:
  -dsn string
       
  -host string
    	The hostname to listen for requests on (default "localhost")
  -port int
    	The port number to listen for requests on (default 8080)
  -table string
    	 The name of the MySQL table (indexed by go-whosonfirst-mysql) to query (default "geojson")
```

For example:

```
./bin/wof-mysql-readerd -dsn '{USER}:{PASSWORD}@/{DATABASE}' -port 7778
2018/05/03 16:53:57 listening for requests on localhost:7778

curl -s localhost:7778/102/547/905/102547905.geojson | jq '.properties["wof:name"]'
"Suvarnabhumi International Airport"
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-readwrite
* https://github.com/whosonfirst/go-whosonfirst-mysql