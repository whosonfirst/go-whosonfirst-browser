package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-whosonfirst-findingaid"
	"github.com/whosonfirst/go-whosonfirst-uri"
	_ "log"
	"net/url"
	"path/filepath"
	"strings"
)

// CacheResolver is a struct that implements the findingaid.Resolver interface for information about Who's On First repositories.
type CacheResolver struct {
	findingaid.Resolver
	cache cache.Cache
}

func init() {

	ctx := context.Background()

	schemes := []string{
		"repo", // deprecated
		"repo-cache",
	}

	for _, s := range schemes {
		err := findingaid.RegisterResolver(ctx, s, NewCacheResolver)

		if err != nil {
			panic(err)
		}
	}
}

// NewCacheResolver returns a findingaid.Resolver instance for exposing information about Who's On First repositories
func NewCacheResolver(ctx context.Context, uri string) (findingaid.Resolver, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	cache_uri := q.Get("cache")

	if cache_uri == "" {
		return nil, errors.New("Missing cache URI")
	}

	_, err = url.Parse(cache_uri)

	if err != nil {
		return nil, err
	}

	c, err := cache.NewCache(ctx, cache_uri)

	if err != nil {
		return nil, err
	}

	fa := &CacheResolver{
		cache: c,
	}

	return fa, nil
}

// ResolveURI will return 'repo.FindingAidResponse' for 'str_response' if it present in the finding aid.
func (fa *CacheResolver) ResolveURI(ctx context.Context, str_uri string) (interface{}, error) {

	key, err := cacheKeyFromURI(str_uri)

	if err != nil {
		return nil, err
	}

	fh, err := fa.cache.Get(ctx, key)

	if err != nil {
		return nil, err
	}

	var rsp *FindingAidResponse

	dec := json.NewDecoder(fh)
	err = dec.Decode(&rsp)

	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func cacheKeyFromURI(str_uri string) (string, error) {

	id, uri_args, err := uri.ParseURI(str_uri)

	if err != nil {
		return "", err
	}

	rel_path, err := uri.Id2RelPath(id, uri_args)

	if err != nil {
		return "", err
	}

	return cacheKeyFromRelPath(rel_path)
}

func cacheKeyFromRelPath(rel_path string) (string, error) {

	ext := filepath.Ext(rel_path)
	rel_path = strings.Replace(rel_path, ext, "", 1)

	key := fmt.Sprintf("%s.json", rel_path)
	return key, nil
}
