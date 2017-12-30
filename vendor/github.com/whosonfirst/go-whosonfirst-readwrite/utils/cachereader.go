package utils

import (
	"github.com/whosonfirst/go-whosonfirst-readwrite/cache"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"io"
	"log"
)

type CacheReader struct {
	reader.Reader
	reader  reader.Reader
	cache   cache.Cache
	options *CacheReaderOptions
}

type CacheReaderOptions struct {
	Debug bool
	Strict	bool	
}

func NewDefaultCacheReaderOptions() (*CacheReaderOptions, error) {

	opts := CacheReaderOptions{
		Debug: false,
		Strict: false,
	}

	return &opts, nil
}

func NewCacheReader(r reader.Reader, c cache.Cache, opts *CacheReaderOptions) (reader.Reader, error) {

	cr := CacheReader{
		reader:  r,
		cache:   c,
		options: opts,
	}

	return &cr, nil
}

func (r *CacheReader) Read(key string) (io.ReadCloser, error) {

	fh, err := r.cache.Get(key)

	if r.options.Debug {
		log.Println("GET", key, err)
	}

	if err == nil {

		if r.options.Debug {
			log.Println("HIT", key)
		}

		return fh, nil
	}

	if err != nil && err.Error() != "CACHE MISS" {
		return nil, err
	}

	if r.options.Debug {
		log.Println("MISS", key)
	}

	fh, err = r.reader.Read(key)

	if r.options.Debug {
		log.Println("READ", key, err)
	}

	if err != nil {
		return nil, err
	}

	fh, err = r.cache.Set(key, fh)

	if r.options.Debug {
		log.Println("SET", key, err)
	}

	if err != nil && r.options.Strict {
		return nil, err
	}	
	
	return fh, nil
}
