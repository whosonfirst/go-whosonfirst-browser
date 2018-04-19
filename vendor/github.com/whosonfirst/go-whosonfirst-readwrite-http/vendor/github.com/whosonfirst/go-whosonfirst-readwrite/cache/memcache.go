package cache

// https://godoc.org/github.com/bradfitz/gomemcache/memcache

import (
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/whosonfirst/go-whosonfirst-readwrite/bytes"
	"io"
	"io/ioutil"
	"strings"
	"sync/atomic"
)

type MemcacheCache struct {
	Cache
	Options   *MemcacheCacheOptions
	cache     *memcache.Client
	hits      int64
	misses    int64
	evictions int64
	keys      int64
}

type MemcacheCacheOptions struct {
	Hosts []string
}

func DefaultMemcacheCacheOptions() (*MemcacheCacheOptions, error) {

	hosts := make([]string, 0)

	opts := MemcacheCacheOptions{
		Hosts: hosts,
	}

	return &opts, nil
}

func MemcacheCacheOptionsFromArgs(args map[string]string) (*MemcacheCacheOptions, error) {

	opts, err := DefaultMemcacheCacheOptions()

	if err != nil {
		return nil, err
	}

	str_hosts, ok := args["Hosts"]

	if ok {

		hosts := strings.Split(str_hosts, " ")
		opts.Hosts = hosts
	}

	return opts, nil
}

func NewMemcacheCache(opts *MemcacheCacheOptions) (Cache, error) {

	mc := memcache.New(opts.Hosts...)

	lc := MemcacheCache{
		Options:   opts,
		cache:     mc,
		hits:      int64(0),
		misses:    int64(0),
		evictions: int64(0),
		keys:      0,
	}

	return &lc, nil
}

func (c *MemcacheCache) Get(key string) (io.ReadCloser, error) {

	it, err := c.cache.Get(key)

	if err != nil {
		atomic.AddInt64(&c.misses, 1)
		return nil, errors.New("CACHE MISS")
	}

	atomic.AddInt64(&c.hits, 1)

	body := it.Value
	return bytes.ReadCloserFromBytes(body)
}

func (c *MemcacheCache) Set(key string, fh io.ReadCloser) (io.ReadCloser, error) {

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

	it := memcache.Item{
		Key:   key,
		Value: body,
	}

	err = c.cache.Set(&it)

	if err != nil {
		return fh, err
	}

	atomic.AddInt64(&c.keys, 1)

	return fh, nil
}

func (c *MemcacheCache) Unset(key string) error {
	return c.cache.Delete(key)
}

func (c *MemcacheCache) Size() int64 {
	return atomic.LoadInt64(&c.keys)
}

func (c *MemcacheCache) Hits() int64 {
	return atomic.LoadInt64(&c.hits)
}

func (c *MemcacheCache) Misses() int64 {
	return atomic.LoadInt64(&c.misses)
}

func (c *MemcacheCache) Evictions() int64 {
	return atomic.LoadInt64(&c.evictions)
}
