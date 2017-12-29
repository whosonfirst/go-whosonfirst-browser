package cache

import (
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru"
	"github.com/whosonfirst/go-whosonfirst-readwrite/bytes"
	"io"
	"io/ioutil"
	"strconv"
	"sync/atomic"
)

type LRUCache struct {
	Cache
	Options   *LRUCacheOptions
	cache     *lru.Cache
	hits      int64
	misses    int64
	evictions int64
	keys      int64
}

type LRUCacheOptions struct {
	CacheSize int
}

func (o *LRUCacheOptions) String() string {
	return fmt.Sprintf("cache size %d", o.CacheSize)
}

func DefaultLRUCacheOptions() (*LRUCacheOptions, error) {

	opts := LRUCacheOptions{
		CacheSize: 100,
	}

	return &opts, nil
}

func LRUCacheOptionsFromArgs(args map[string]string) (*LRUCacheOptions, error) {

	opts, err := DefaultLRUCacheOptions()

	if err != nil {
		return nil, err
	}

	str_sz, ok := args["CacheSize"]

	if ok {

		sz, err := strconv.Atoi(str_sz)

		if err != nil {
			return nil, err
		}

		opts.CacheSize = sz
	}

	return opts, nil
}

func NewLRUCache(opts *LRUCacheOptions) (Cache, error) {

	c, err := lru.New(opts.CacheSize)

	if err != nil {
		return nil, err
	}

	lc := LRUCache{
		Options:   opts,
		cache:     c,
		hits:      int64(0),
		misses:    int64(0),
		evictions: int64(0),
		keys:      0,
	}

	return &lc, nil
}

func (c *LRUCache) Get(key string) (io.ReadCloser, error) {

	cache, ok := c.cache.Get(key)

	if !ok {
		atomic.AddInt64(&c.misses, 1)
		return nil, errors.New("CACHE MISS")
	}

	atomic.AddInt64(&c.hits, 1)

	body := cache.([]byte)
	return bytes.ReadCloserFromBytes(body)
}

func (c *LRUCache) Set(key string, fh io.ReadCloser) (io.ReadCloser, error) {

	/*

	   Assume an io.Reader is hooked up to a satellite dish receiving a message (maybe a 1TB message) from an
	   alien civilization who only transmits their message once every thousand years. There's no "rewinding"
	   that.

	   https://groups.google.com/forum/#!msg/golang-nuts/BzDAg0CFqyk/t3TvH9QV0xEJ

	*/

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	evicted := c.cache.Add(key, body)
	atomic.AddInt64(&c.keys, 1)

	if evicted {
		atomic.AddInt64(&c.evictions, 1)
		atomic.AddInt64(&c.keys, -1)
	}

	return bytes.ReadCloserFromBytes(body)
}

func (c *LRUCache) Size() int64 {
	return atomic.LoadInt64(&c.keys)
}

func (c *LRUCache) Hits() int64 {
	return atomic.LoadInt64(&c.hits)
}

func (c *LRUCache) Misses() int64 {
	return atomic.LoadInt64(&c.misses)
}

func (c *LRUCache) Evictions() int64 {
	return atomic.LoadInt64(&c.evictions)
}
