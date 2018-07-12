# go-whosonfirst-mysql

Go package for working with Who's On First documents and MySQL databases.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## A few things before we get started

1. This package assumes you are running a version of [MySQL](https://dev.mysql.com/doc/refman/5.7/en/spatial-analysis-functions.html) (or [MariaDB](https://mariadb.com/kb/en/library/geographic-geometric-features/)) with spatial extensions, so version 5.7 or higher.
2. This package assumes Who's On First documents and is not yet able to index arbitrary GeoJSON documents.
3. This package shares the same basic model as the [go-whosonfirst-sqlite-*](https://github.com/whosonfirst?utf8=%E2%9C%93&q=go-whosonfirst-sqlite&type=&language=) packages. They should be reconciled. Today, they are not.
4. This is not an abstract package for working with databases and tables that aren't Who's On First specific, the way [go-whosonfirst-sqlite](https://github.com/whosonfirst/go-whosonfirst-sqlite) is. It probably _should_ be but that seems like something that will happen as a result of doing #3 (above). 

## Interfaces

### Database

```
type Database interface {
     Conn() (*sql.DB, error)
     DSN() string
     Close() error
}
```

### Table

```
type Table interface {
     Name() string
     Schema() string
     InitializeTable(Database) error
     IndexRecord(Database, interface{}) error
}
```

It is left up to people implementing the `Table` interface to figure out what to do with the second value passed to the `IndexRecord` method. For example:

```
func (t *WhosonfirstTable) IndexRecord(db mysql.Database, i interface{}) error {
	return t.IndexFeature(db, i.(geojson.Feature))
}

func (t *WhosonfirstTable) IndexFeature(db mysql.Database, f geojson.Feature) error {
	// code to index geojson.Feature here
}
```

## Tables

### geojson

```
CREATE TABLE IF NOT EXISTS geojson (
      id BIGINT UNSIGNED PRIMARY KEY,
      body LONGBLOB NOT NULL,
      lastmodified INT NOT NULL,
      KEY lastmodified (lastmodified)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

### whosonfirst

```
CREATE TABLE IF NOT EXISTS %s (
      id BIGINT UNSIGNED PRIMARY KEY,
      properties JSON NOT NULL,
      geometry GEOMETRY NOT NULL,
      centroid POINT NOT NULL COMMENT 'This is not necessary a math centroid',
      lastmodified INT NOT NULL,
      parent_id BIGINT       GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties,'$."wof:parent_id"'))) VIRTUAL,
      placetype VARCHAR(64)  GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties,'$."wof:placetype"'))) VIRTUAL,
      is_current TINYINT     GENERATED ALWAYS AS (JSON_CONTAINS_PATH(properties, 'one', '$."mz:is_current"') AND JSON_UNQUOTE(JSON_EXTRACT(properties,'$."mz:is_current"'))) VIRTUAL,
      is_nullisland TINYINT  GENERATED ALWAYS AS (JSON_CONTAINS_PATH(properties, 'one', '$."mz:is_nullisland"') AND JSON_LENGTH(JSON_EXTRACT(properties, '$."mz:is_nullisland"'))) VIRTUAL,
      is_approximate TINYINT GENERATED ALWAYS AS (JSON_CONTAINS_PATH(properties, 'one', '$."mz:is_approximate"') AND JSON_LENGTH(JSON_EXTRACT(properties, '$."mz:is_approximate"'))) VIRTUAL,
      is_ceased TINYINT      GENERATED ALWAYS AS (JSON_CONTAINS_PATH(properties, 'one', '$."edtf:cessation"') AND JSON_UNQUOTE(JSON_EXTRACT(properties,'$."edtf:cessation"')) != "" AND JSON_UNQUOTE(JSON_EXTRACT(properties,'$."edtf:cessation"')) != "open" AND json_unquote(json_extract(properties,'$."edtf:cessation"')) != "uuuu") VIRTUAL,
      is_deprecated TINYINT  GENERATED ALWAYS AS (JSON_CONTAINS_PATH(properties, 'one', '$."edtf:deprecated"') AND JSON_UNQUOTE(JSON_EXTRACT(properties,'$."edtf:deprecated"')) != "" AND json_unquote(json_extract(properties,'$."edtf:deprecated"')) != "uuuu") VIRTUAL,
      is_superseded TINYINT  GENERATED ALWAYS AS (JSON_LENGTH(JSON_EXTRACT(properties, '$."wof:superseded_by"')) > 0) VIRTUAL,
      is_superseding TINYINT GENERATED ALWAYS AS (JSON_LENGTH(JSON_EXTRACT(properties, '$."wof:supersedes"')) > 0) VIRTUAL,
      date_upper DATE	     GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties, '$."date:cessation_upper"'))) VIRTUAL,
      date_lower DATE	     GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(properties, '$."date:inception_lower"'))) VIRTUAL,
      KEY parent_id (parent_id),
      KEY placetype (placetype),
      KEY is_current (is_current),
      KEY is_nullisland (is_nullisland),
      KEY is_approximate (is_approximate),
      KEY is_deprecated (is_deprecated),
      KEY is_superseded (is_superseded),
      KEY is_superseding (is_superseding),
      KEY date_upper (date_upper),
      KEY date_lower (date_lower),
      SPATIAL KEY idx_geometry (geometry),
      SPATIAL KEY idx_centroid (centroid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
```

There are a few important things to note about the `whosonfirst` table:

1. It is technically possible to add VIRTUAL centroid along the lines of `centroid POINT GENERATED ALWAYS AS (ST_Centroid(geometry)) VIRTUAL` we don't because MySQL will return the math centroid and well we all know what that means for places like San Francisco (SF) - if you don't it means the [math centroid will be in the Pacific Ocean](https://spelunker.whosonfirst.org/id/85922583/) because technically the Farralon Islands are part of SF - so instead we we compute the centroid in the code (using the go-whosonfirst-geojson-v2 Centroid interface)
2. It's almost certainly going to be moved in to a different package (once this code base is reconciled with the `go-whosonfirst-sqlite` packages)
3. It is now a _third_ way to "spatially" store WOF records, along with the [go-whosonfirst-sqlite-features `geometries`](https://github.com/whosonfirst/go-whosonfirst-sqlite-features#geometries) and the [go-whosonfirst-spatialite-geojson geojson](https://github.com/whosonfirst/go-whosonfirst-spatialite-geojson#geojson) tables. It is entirely possible that this is "just how it is" and there is no value in a single unified table schema but, equally, it seems like it's something to have a think about.

## Custom tables

Sure. You just need to write a per-table package that implements the `Table` interface, described above.

## Tools

### wof-mysql-index 

```
./bin/wof-mysql-index -h
Usage of ./bin/wof-mysql-index:
  -all
	Index all the tables
  -config string
    	  Read some or all flags from an ini-style config file. Values in the config file take precedence over command line flags.
  -dsn string
       A valid go-sql-driver DSN string, for example '{USER}:{PASSWORD}@/{DATABASE}'
  -geojson
	Index the 'geojson' tables
  -mode string
    	The mode to use importing data. Valid modes are: directory,feature,feature-collection,files,geojson-ls,meta,path,repo,sqlite. (default "repo")
  -section string
    	   A valid ini-style config file section. (default "wof-mysql")
  -timings
	Display timings during and after indexing
  -whosonfirst
	Index the 'whosonfirst' tables
```

For example:

```
./bin/wof-mysql-index -dsn '{USER}:{PASSWORD}@/{DATABASE}' /usr/local/data/whosonfirst-data/
```

### Config files

You can read (or override) command line flags from a config file, by passing the `-config` flag with the path to a valid ini-style config file. For example, assuming a config file like this:

```
[wof-mysql]
dsn={USER}:{PASS}@/{DATABASE}
all
timings
```

You might invoke it like this:

```
./bin/wof-mysql-index -config ./test.cfg /usr/local/data/whosonfirst-data-*
13:47:57.021711 [wof-mysql-index] STATUS Reset all flag from config file
13:47:57.021840 [wof-mysql-index] STATUS Reset dsn flag from config file
13:47:57.021846 [wof-mysql-index] STATUS Reset timings flag from config file
13:48:57.037310 [wof-mysql-index] STATUS time to index geojson (3155) : 16.979713633s
13:48:57.037329 [wof-mysql-index] STATUS time to index whosonfirst (3155) : 29.342492075s
13:48:57.037334 [wof-mysql-index] STATUS time to index all (3155) : 1m0.013715096s
... and so on
```

### Environment variables

_Unless_ you are passing the `-config` flag you can set (or override) command line flags with environment variables. Environment variable are expected to:

* Be upper-cased
* Replace all instances of `-` with `_`
* Be prefixed with `WOF_MYSQL`

For example the `-dsn` flag would be overridden by the `WOF_MYSQL_DSN` environment variable.

## See also:

* https://dev.mysql.com/doc/refman/5.7/en/spatial-analysis-functions.html
* https://github.com/whosonfirst/go-whosonfirst-sqlite

* https://dev.mysql.com/doc/refman/8.0/en/json-functions.html
* https://www.percona.com/blog/2016/03/07/json-document-fast-lookup-with-mysql-5-7/
* https://archive.fosdem.org/2016/schedule/event/mysql57_json/attachments/slides/1291/export/events/attachments/mysql57_json/slides/1291/MySQL_57_JSON.pdf


