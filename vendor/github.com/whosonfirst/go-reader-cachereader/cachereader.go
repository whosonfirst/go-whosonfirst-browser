package cachereader

import (
	"context"
	"errors"
	"github.com/whosonfirst/go-cache"
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

	reader_uri := q.Get("reader")

	if reader_uri == "" {
		return nil, errors.New("Missing ?reader= parameter")
	}

	cache_uri := q.Get("cache")

	if cache_uri == "" {
		return nil, errors.New("Missing ?cache= parameter")
	}

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		return nil, err
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

func (cr *CacheReader) Read(ctx context.Context, key string) (io.ReadCloser, error) {

	fh, err := cr.cache.Get(ctx, key)

	if err == nil {
		return fh, nil
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

	return fh, nil
}

func (cr *CacheReader) URI(key string) string {
	return cr.reader.URI(key)
}
