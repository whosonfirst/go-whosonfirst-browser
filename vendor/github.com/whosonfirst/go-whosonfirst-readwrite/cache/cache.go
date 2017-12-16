package cache

import (
	"errors"
	"io"
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

	switch source {
	case "gocache":

		opts, opts_err := DefaultGoCacheOptions()

		if opts_err != nil {
			err = opts_err
		} else {
			c, err = NewGoCache(opts)
		}
	case "null":
		c, err = NewNullCache()
	default:
		err = errors.New("Unknown or invalid source")
	}

	return c, err
}
