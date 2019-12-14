package cache

import (
	"context"
	"github.com/aaronland/go-roster"
	"io"
	_ "log"
	"net/url"
)

type Cache interface {
	Open(context.Context, string) error
	Close(context.Context) error
	Name() string
	Get(context.Context, string) (io.ReadCloser, error)
	Set(context.Context, string, io.ReadCloser) (io.ReadCloser, error)
	Unset(context.Context, string) error
	Hits() int64
	Misses() int64
	Evictions() int64
	Size() int64
	SizeWithContext(context.Context) int64
}

var caches roster.Roster

func ensureRoster() error {

	if caches == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		caches = r
	}

	return nil
}

func RegisterCache(ctx context.Context, name string, c Cache) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return caches.Register(ctx, name, c)
}

func NewCache(ctx context.Context, uri string) (Cache, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := caches.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	c := i.(Cache)

	err = c.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return c, nil
}
