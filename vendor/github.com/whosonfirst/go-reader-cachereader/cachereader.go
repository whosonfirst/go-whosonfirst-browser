package cachereader

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-reader"
	"io"
	"net/url"
)

type CacheReader struct {
	reader.Reader
	reader reader.Reader
	cache  cache.Cache
}

func init() {

	ctx := context.Background()
	err := reader.RegisterReader(ctx, "cachereader", NewCacheReader)

	if err != nil {
		panic(err)
	}
}

func NewCacheReader(ctx context.Context, uri string) (reader.Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
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
		return nil, err
	}

	cr := &CacheReader{
		reader: r,
		cache:  c,
	}

	return cr, nil
}

func (cr *CacheReader) Read(ctx context.Context, key string) (io.ReadSeekCloser, error) {

	fh, err := cr.cache.Get(ctx, key)

	if err == nil {
		// https://github.com/whosonfirst/go-cache/issues/1
		return ioutil.NewReadSeekCloser(fh)
	}

	if !cache.IsCacheMiss(err) {
		return nil, err
	}

	fh, err = cr.reader.Read(ctx, key)

	if err != nil {
		return nil, err
	}

	fh, err = cr.cache.Set(ctx, key, fh)

	if err != nil {
		return nil, err
	}

	// https://github.com/whosonfirst/go-cache/issues/1
	return ioutil.NewReadSeekCloser(fh)
}

func (cr *CacheReader) ReaderURI(ctx context.Context, key string) string {
	return cr.reader.ReaderURI(ctx, key)
}
