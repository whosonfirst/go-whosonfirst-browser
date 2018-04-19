package cache

// https://godoc.org/github.com/allegro/bigcache

import (
	"errors"
	"github.com/allegro/bigcache"
	"github.com/whosonfirst/go-whosonfirst-readwrite/bytes"
	"io"
	"io/ioutil"
	_ "log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type BigCacheCache struct {
	Cache
	Options   *BigCacheCacheOptions
	cache     *bigcache.BigCache
	hits      int64
	misses    int64
	evictions int64
	keys      int64
}

type BigCacheCacheOptions struct {
	HardMaxCacheSize int
	MaxEntrySize     int
	TimeToExpire     time.Duration
}

func DefaultBigCacheCacheOptions() (*BigCacheCacheOptions, error) {

	opts := BigCacheCacheOptions{
		HardMaxCacheSize: 0,
		MaxEntrySize:     0,
		TimeToExpire:     300 * time.Second,
	}

	return &opts, nil
}

func BigCacheCacheOptionsFromArgs(args map[string]string) (*BigCacheCacheOptions, error) {

	opts, err := DefaultBigCacheCacheOptions()

	if err != nil {
		return nil, err
	}

	str_cache_sz, ok := args["HardMaxCacheSize"]

	if ok {

		sz, err := strconv.Atoi(str_cache_sz)

		if err != nil {
			return nil, err
		}

		opts.HardMaxCacheSize = sz
	}

	str_max_sz, ok := args["MaxEntrySize"]

	if ok {

		sz, err := strconv.Atoi(str_max_sz)

		if err != nil {
			return nil, err
		}

		opts.MaxEntrySize = sz
	}

	str_tte, ok := args["TimeToExpire"]

	if ok {

		sz, err := strconv.Atoi(str_tte)

		if err != nil {
			return nil, err
		}

		opts.TimeToExpire = time.Duration(sz) * time.Second
	}

	return opts, nil
}

func NewBigCacheCache(opts *BigCacheCacheOptions) (Cache, error) {

	evict := opts.TimeToExpire * time.Second
	config := bigcache.DefaultConfig(evict)

	if opts.MaxEntrySize > 0 {
		config.MaxEntrySize = opts.MaxEntrySize
	}

	if opts.HardMaxCacheSize > 0 {
		config.HardMaxCacheSize = opts.HardMaxCacheSize
	}

	c, err := bigcache.NewBigCache(config)

	if err != nil {
		return nil, err
	}

	lc := BigCacheCache{
		Options:   opts,
		cache:     c,
		hits:      int64(0),
		misses:    int64(0),
		evictions: int64(0),
		keys:      0,
	}

	return &lc, nil
}

func (c *BigCacheCache) Get(key string) (io.ReadCloser, error) {

	body, err := c.cache.Get(key)

	if err != nil {
		atomic.AddInt64(&c.misses, 1)

		if strings.HasSuffix(err.Error(), "not found") {
			return nil, errors.New("CACHE MISS")
		}

		return nil, err
	}

	atomic.AddInt64(&c.hits, 1)

	return bytes.ReadCloserFromBytes(body)
}

func (c *BigCacheCache) Set(key string, fh io.ReadCloser) (io.ReadCloser, error) {

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

	fh, err = bytes.ReadCloserFromBytes(body)

	if err != nil {
		return nil, err
	}

	err = c.cache.Set(key, body)

	if err != nil {
		return fh, err
	}

	atomic.AddInt64(&c.keys, 1)

	return fh, nil
}

func (c *BigCacheCache) Unset(key string) error {
	return c.cache.Delete(key)
}

func (c *BigCacheCache) Size() int64 {
	// return atomic.LoadInt64(&c.keys)
	i := c.cache.Len()
	return int64(i)
}

func (c *BigCacheCache) Hits() int64 {
	return atomic.LoadInt64(&c.hits)
}

func (c *BigCacheCache) Misses() int64 {
	return atomic.LoadInt64(&c.misses)
}

func (c *BigCacheCache) Evictions() int64 {
	return atomic.LoadInt64(&c.evictions)
}
