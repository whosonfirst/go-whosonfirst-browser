# go-whosonfirst-readwrite

## Important

It's probably too soon, for you. If nothing else the documentation is incomplete

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Example

```
./bin/wof-reader -cache memcache -cache-arg Hosts=127.0.0.1:11211 -fs-root /usr/local/data/whosonfirst-data/data/ -debug 0/0.geojson
2017/12/29 16:58:08 GET 0/0.geojson <nil>
2017/12/29 16:58:08 HIT 0/0.geojson
2017/12/29 16:58:08 0/0.geojson true
2017/12/29 16:58:08 GET 0/0.geojson <nil>
2017/12/29 16:58:08 HIT 0/0.geojson
```

## Caches

```
type Cache interface {
	Get(string) (io.ReadCloser, error)
	Set(string, io.ReadCloser) (io.ReadCloser, error)
	Unset(string) error
	Hits() int64
	Misses() int64
	Evictions() int64
	Size() int64
}
```

### bigcache

### gocache

### lru

### memcache

## Readers

```
type Reader interface {
	Read(string) (io.ReadCloser, error)
}
```

### fs

### github

### http

### null

### s3

## Writers

```
type Writer interface {
	Write(string, io.ReadCloser) error
}
```

### multi

### null

### s3

### stdout

