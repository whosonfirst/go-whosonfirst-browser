// package resolver providers common methods and interfaces for retrieving repository data from a variety of storage systems.
package resolver

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

// type Resolver defines a storage-independent interface for retrieving a repository name given an ID.
type Resolver interface {
	// GetRepo returns the repository name matching an ID.
	GetRepo(context.Context, int64) (string, error)
}

// type ResolverInitializeFunc defines an initialization function for a storage-specific implementation of the Resolver interface.
type ResolverInitializeFunc func(ctx context.Context, uri string) (Resolver, error)

var resolvers roster.Roster

func ensureSpatialRoster() error {

	if resolvers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		resolvers = r
	}

	return nil
}

func RegisterResolver(ctx context.Context, scheme string, f ResolverInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return resolvers.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range resolvers.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
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

	f := i.(ResolverInitializeFunc)
	return f(ctx, uri)
}
