package cachereader

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-reader"
	"io"
	"net/url"
	"sync"
)

// type CacheReaderOptions is a struct for use with the `NewCacheReaderWithOptions` method.
type CacheReaderOptions struct {
	Reader reader.Reader
	Cache  cache.Cache
}

// type CacheReader implements the `whosonfirst/go-reader` interface for use with a caching layer for reading documents.
type CacheReader struct {
	reader.Reader
	// A `whosonfirst/go-reader.Reader` instance used to read source documents.
	reader reader.Reader
	// A `whosonfirst/go-cache.Cache` instance used to cache source documents.
	cache cache.Cache
}

// type CacheStatus is a constant indicating the cache status for a read action
type CacheStatus uint8

const (
	// CacheNotFound signals that there are no recorded events for a cache lookup
	CacheNotFound CacheStatus = iota
	// CacheHit signals that the last read event was a cache HIT
	CacheHit
	// CacheHit signals that the last read event was a cache MISS
	CacheMiss
)

var lastReadMap *sync.Map

func init() {

	ctx := context.Background()
	err := reader.RegisterReader(ctx, "cachereader", NewCacheReader)

	if err != nil {
		panic(err)
	}

	lastReadMap = new(sync.Map)
}

func lastReadKey(cr reader.Reader, key string) string {
	return fmt.Sprintf("%p#%s", cr, key)
}

func setLastRead(cr reader.Reader, key string, s CacheStatus) {
	fq_key := lastReadKey(cr, key)
	lastReadMap.Store(fq_key, s)
}

// GetLastRead returns the `CacheStatus` value for the last event using 'cr' to read 'path'.
func GetLastRead(cr reader.Reader, key string) (CacheStatus, bool) {

	fq_key := lastReadKey(cr, key)

	v, _ := lastReadMap.Load(fq_key)

	if v == nil {
		return CacheNotFound, false
	}

	s := v.(CacheStatus)
	return s, true
}

// NewCacheReader will return a new `CacheReader` instance configured by 'uri' which
// is expected to take the form of:
//
//	cachereader://?reader={READER_URI}&cache={CACHE_URI}
//
// Where {READER_URI} is expected to be a valid `whosonfirst/go-reader.Reader` URI and
// {CACHE_URI} is expected to be a valid `whosonfirst/go-cache.Cache` URI. Note that multiple
// "?reader=" parameter are supported; internally the `CacheReader` implementation uses a
// `whosonfirst/go-reader.MultiReader` instance for reading documents.
func NewCacheReader(ctx context.Context, uri string) (reader.Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	cache_uri := q.Get("cache")

	if cache_uri == "" {
		return nil, fmt.Errorf("Missing ?cache= parameter")
	}

	reader_uris := q["reader"]

	r, err := reader.NewMultiReaderFromURIs(ctx, reader_uris...)

	if err != nil {
		return nil, fmt.Errorf("Failed to create reader, %w", err)
	}

	c, err := cache.NewCache(ctx, cache_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new cache, %w", err)
	}

	opts := &CacheReaderOptions{
		Reader: r,
		Cache:  c,
	}

	return NewCacheReaderWithOptions(ctx, opts)
}

// NewCacheReader will return a new `CacheReader` instance configured by 'opts'.
func NewCacheReaderWithOptions(ctx context.Context, opts *CacheReaderOptions) (reader.Reader, error) {

	cr := &CacheReader{
		reader: opts.Reader,
		cache:  opts.Cache,
	}

	return cr, nil
}

// Read returns an `io.ReadSeekCloser` instance for the document resolved by `uri`. The document
// will also be added to the internal cache maintained by 'cr' if it not already present.
func (cr *CacheReader) Read(ctx context.Context, key string) (io.ReadSeekCloser, error) {

	fh, err := cr.cache.Get(ctx, key)

	if err == nil {
		setLastRead(cr, key, CacheHit)
		return ioutil.NewReadSeekCloser(fh)
	}

	if !cache.IsCacheMiss(err) {
		return nil, fmt.Errorf("Failed to read from cache for %s with %T, %w", key, cr.cache, err)
	}

	setLastRead(cr, key, CacheMiss)
	fh, err = cr.reader.Read(ctx, key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read %s from %T, %w", key, cr.reader, err)
	}

	fh, err = cr.cache.Set(ctx, key, fh)

	if err != nil {
		return nil, fmt.Errorf("Failed to set cache for %s with %T, %w", key, cr.cache, err)
	}

	// https://github.com/whosonfirst/go-cache/issues/1
	return ioutil.NewReadSeekCloser(fh)
}

// ReaderURI returns final URI resolved by `uri` for this reader.
func (cr *CacheReader) ReaderURI(ctx context.Context, key string) string {
	return cr.reader.ReaderURI(ctx, key)
}
