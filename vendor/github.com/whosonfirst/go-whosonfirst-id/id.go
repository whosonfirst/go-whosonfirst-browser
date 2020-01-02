package id

import (
	"context"
	_ "github.com/aaronland/go-brooklynintegers-api"
	"github.com/aaronland/go-uid"
	"github.com/aaronland/go-uid-artisanal"
	"strconv"
)

type Provider interface {
	NewID() (int64, error)
}

type WOFProvider struct {
	Provider
	uid_provider uid.Provider
}

func NewProvider(ctx context.Context) (Provider, error) {

	opts := &artisanal.ArtisanalProviderURIOptions{
		Pool:    "memory://",
		Minimum: 0,
		Clients: []string{
			"brooklynintegers://",
		},
	}

	uri, err := artisanal.NewArtisanalProviderURI(opts)

	if err != nil {
		return nil, err
	}

	// str_uri ends up looking like this:
	// artisanal:?client=brooklynintegers%3A%2F%2F&minimum=5&pool=memory%3A%2F%2F

	str_uri := uri.String()

	return NewProviderWithURI(ctx, str_uri)
}

func NewProviderWithURI(ctx context.Context, uri string) (Provider, error) {

	uid_pr, err := uid.NewProvider(ctx, uri)

	if err != nil {
		return nil, err
	}

	wof_pr := &WOFProvider{
		uid_provider: uid_pr,
	}

	return wof_pr, nil
}

func (wof_pr *WOFProvider) NewID() (int64, error) {

	uid, err := wof_pr.uid_provider.UID()

	if err != nil {
		return -1, err
	}

	str_id := uid.String()
	return strconv.ParseInt(str_id, 10, 64)
}
