package cachewriter

import (
	"bytes"
	"context"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-writer"
	"io"
	"io/ioutil"
)

type CacheWriter struct {
	writer.Writer
	writer writer.Writer
	cache  cache.Cache
}

func NewCacheWriter(r writer.Writer, c cache.Cache) (writer.Writer, error) {

	cw := &CacheWriter{
		writer: r,
		cache:  c,
	}

	return cw, nil
}

func (cw *CacheWriter) Open(ctx context.Context, uri string) error {
	return nil
}

func (cw *CacheWriter) Write(ctx context.Context, key string, fh io.ReadCloser) error {

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return err
	}

	br := bytes.NewReader(body)
	fh = ioutil.NopCloser(br)

	err = cw.writer.Write(ctx, key, fh)

	if err != nil {
		return err
	}

	br = bytes.NewReader(body)
	fh = ioutil.NopCloser(br)

	_, err = cw.cache.Set(ctx, key, fh)

	if err == nil {
		return err
	}

	return nil
}
