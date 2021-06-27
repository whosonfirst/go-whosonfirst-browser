package findingaid

import (
	"context"
	"github.com/aaronland/go-roster"
	"io"
	"net/url"
)

type Indexer interface {
	IndexURIs(context.Context, ...string) error
	IndexReader(context.Context, io.Reader) error
}

type IndexerInitializationFunc func(ctx context.Context, uri string) (Indexer, error)

var indexers roster.Roster

func ensureIndexers() error {

	if indexers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		indexers = r
	}

	return nil
}

func RegisterIndexer(ctx context.Context, name string, c IndexerInitializationFunc) error {

	err := ensureIndexers()

	if err != nil {
		return err
	}

	return indexers.Register(ctx, name, c)
}

func NewIndexer(ctx context.Context, uri string) (Indexer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := indexers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init := i.(IndexerInitializationFunc)
	c, err := init(ctx, uri)

	if err != nil {
		return nil, err
	}

	return c, nil
}
