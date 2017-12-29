package utils

import (
	"github.com/whosonfirst/go-whosonfirst-readwrite/cache"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"io"
	"log"
)

type CacheReader struct {
	reader.Reader
	reader reader.Reader
	cache  cache.Cache
	debug  bool
}

func NewCacheReader(r reader.Reader, c cache.Cache) (reader.Reader, error) {

	cr := CacheReader{
		reader: r,
		cache:  c,
		debug:  false,
	}

	return &cr, nil
}

func (r *CacheReader) Read(key string) (io.ReadCloser, error) {

	fh, err := r.cache.Get(key)

	if r.debug {
		log.Println("GET", key, fh, err)
	}

	if err == nil {

		if r.debug {
			log.Println("HIT", key)
		}

		return fh, nil
	}

	if err != nil && err.Error() != "CACHE MISS" {
		return nil, err
	}

	if r.debug {
		log.Println("MISS", key)
	}

	fh, err = r.reader.Read(key)

	if r.debug {
		log.Println("READ", key, fh, err)
	}

	if err != nil {
		return nil, err
	}

	fh, err = r.cache.Set(key, fh)

	if r.debug {
		log.Println("SET", key, fh, err)
	}

	if err != nil {
		return nil, err
	}

	return fh, nil
}
