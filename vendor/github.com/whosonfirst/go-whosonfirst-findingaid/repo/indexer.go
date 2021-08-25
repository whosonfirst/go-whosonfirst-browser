package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-findingaid"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"io"
	_ "log"
	"net/url"
)

// Indexer is a struct that implements the findingaid.Indexer interface for information about Who's On First repositories.
type Indexer struct {
	findingaid.Indexer
	cache        cache.Cache
	iterator_uri string
}

func init() {

	ctx := context.Background()
	err := findingaid.RegisterIndexer(ctx, "repo", NewIndexer)

	if err != nil {
		panic(err)
	}
}

// NewIndexer returns a findingaid.Indexer instance for exposing information about Who's On First repositories
func NewIndexer(ctx context.Context, uri string) (findingaid.Indexer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	cache_uri := q.Get("cache")
	iterator_uri := q.Get("iterator")

	if cache_uri == "" {
		return nil, errors.New("Missing ?cache= parameter.")
	}

	c, err := cache.NewCache(ctx, cache_uri)

	if err != nil {
		return nil, err
	}

	if iterator_uri == "" {
		return nil, errors.New("Missing ?iterator= parameter.")
	}

	// We defer creating the iterator until the 'IndexURIs' method is
	// invoked because the iterator callback has a reference to this
	// (findingaid indexer) instance which hasn't been created at this
	// point.

	_, err = url.Parse(iterator_uri)

	if err != nil {
		return nil, fmt.Errorf("Invalid ?iterator= parameter, %w", err)
	}

	fa := &Indexer{
		cache:        c,
		iterator_uri: iterator_uri,
	}

	return fa, nil
}

// Index will index records defined by 'sources...' in the finding aid, using the whosonfirst/go-whosonfirst-iterate package.
func (fa *Indexer) IndexURIs(ctx context.Context, sources ...string) error {

	cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		return fa.IndexReader(ctx, fh)
	}

	iter, err := iterator.NewIterator(ctx, fa.iterator_uri, cb)

	if err != nil {
		return err
	}

	return iter.IterateURIs(ctx, sources...)
}

// IndexReader will index an individual Who's On First record in the finding aid.
func (fa *Indexer) IndexReader(ctx context.Context, fh io.Reader) error {

	rsp, err := FindingAidResponseFromReader(ctx, fh)

	if err != nil {
		return err
	}

	enc, err := json.Marshal(rsp)

	if err != nil {
		return err
	}

	br := bytes.NewReader(enc)
	rsc, err := ioutil.NewReadSeekCloser(br)

	if err != nil {
		return err
	}

	key, err := cacheKeyFromRelPath(rsp.URI)

	if err != nil {
		return err
	}

	_, err = fa.cache.Set(ctx, key, rsc)

	if err != nil {
		return err
	}

	return nil
}
