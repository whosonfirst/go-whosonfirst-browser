package cache

import (
	"errors"
	"io"
	_ "log"
	"strings"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type Cache interface {
	// GetWithReader(string, reader.Reader) (io.ReadCloser, error)
	Get(string) (io.ReadCloser, error)
	Set(string, io.ReadCloser) (io.ReadCloser, error)
	Hits() int64
	Misses() int64
	Evictions() int64
	Size() int64
}

func NewCacheFromSource(source string, args ...interface{}) (Cache, error) {

	var c Cache
	var err error

	cache_args := make(map[string]string)

	if len(args) >= 1 {
		cache_args = args[0].(map[string]string)
	}

	switch strings.ToLower(source) {
	case "gocache":

		opts, opts_err := GoCacheOptionsFromArgs(cache_args)

		if opts_err != nil {
			err = opts_err
		} else {
			c, err = NewGoCache(opts)
		}

	case "lru":

		opts, opts_err := LRUCacheOptionsFromArgs(cache_args)

		if opts_err != nil {
			err = opts_err
		} else {
			c, err = NewLRUCache(opts)
		}

	case "null":
		c, err = NewNullCache()
	default:
		err = errors.New("Unknown or invalid source")
	}

	return c, err
}
