package cache

import (
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
