package findingaid

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
)

type Resolver interface {
	ResolveURI(context.Context, string) (interface{}, error)
}

type ResolverInitializationFunc func(ctx context.Context, uri string) (Resolver, error)

var resolvers roster.Roster

func ensureResolvers() error {

	if resolvers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		resolvers = r
	}

	return nil
}

func RegisterResolver(ctx context.Context, name string, c ResolverInitializationFunc) error {

	err := ensureResolvers()

	if err != nil {
		return err
	}

	return resolvers.Register(ctx, name, c)
}

func NewResolver(ctx context.Context, uri string) (Resolver, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := resolvers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init := i.(ResolverInitializationFunc)
	c, err := init(ctx, uri)

	if err != nil {
		return nil, err
	}

	return c, nil
}
