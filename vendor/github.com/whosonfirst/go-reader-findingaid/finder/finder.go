// package finder providers common methods and interfaces for retrieving repository data from a variety of storage systems.
package finder

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

// type Finder defines a storage-independent interface for retrieving a repository name given an ID.
type Finder interface {
	// GetRepo returns the repository name matching an ID.
	GetRepo(context.Context, int64) (string, error)
}

// type FinderInitializeFunc defines an initialization function for a storage-specific implementation of the Finder interface.
type FinderInitializeFunc func(ctx context.Context, uri string) (Finder, error)

var finders roster.Roster

func ensureSpatialRoster() error {

	if finders == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		finders = r
	}

	return nil
}

func RegisterFinder(ctx context.Context, scheme string, f FinderInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return finders.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range finders.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewFinder(ctx context.Context, uri string) (Finder, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := finders.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(FinderInitializeFunc)
	return f(ctx, uri)
}
