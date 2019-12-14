package cachereader

import (
	"context"
	"io"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	
)

type CacheReader struct {
	reader.Reader
	reader reader.Reader
	cache cache.Cache
}

func NewCacheReader(r reader.Reader, c cache.Cache) (reader.Reader, error) {

	cr := &CacheReader{
		reader: r,
		cache: c,
	}

	return cr, nil
}

func (cr *CacheReader) Open(ctx context.Context, uri string) error {
	return nil
}

func (cr *CacheReader) Read(ctx context.Context, key string) (io.ReadCloser, error) {

	fh, err := cr.cache.Get(ctx, key)

	if err == nil {
		return fh, nil
	}

	if !cache.IsCacheMiss(err){
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
