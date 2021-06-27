package fulltext

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-search/filter"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"net/url"
	"sort"
	"strings"
)

type FullTextDatabase interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	QueryString(context.Context, string, ...filter.Filter) (spr.StandardPlacesResults, error)
	Close(context.Context) error
}

type FullTextDatabaseInitializeFunc func(ctx context.Context, uri string) (FullTextDatabase, error)

var fulltext_databases roster.Roster

func ensureFullTextRoster() error {

	if fulltext_databases == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		fulltext_databases = r
	}

	return nil
}

func RegisterFullTextDatabase(ctx context.Context, scheme string, f FullTextDatabaseInitializeFunc) error {

	err := ensureFullTextRoster()

	if err != nil {
		return err
	}

	return fulltext_databases.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureFullTextRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range fulltext_databases.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewFullTextDatabase(ctx context.Context, uri string) (FullTextDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := fulltext_databases.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(FullTextDatabaseInitializeFunc)
	return f(ctx, uri)
}
