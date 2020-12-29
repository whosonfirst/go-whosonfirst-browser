package fulltext

import (
	"context"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/whosonfirst/go-whosonfirst-search/filter"	
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"	
)

type NullFullTextDatabase struct {
	FullTextDatabase
}

func init() {
	ctx := context.Background()
	RegisterFullTextDatabase(ctx, "null", NewNullFullTextDatabase)
}

func NewNullFullTextDatabase(ctx context.Context, str_uri string) (FullTextDatabase, error) {

	ftdb := &NullFullTextDatabase{}
	return ftdb, nil
}

func (ftdb *NullFullTextDatabase) Close(ctx context.Context) error {
	return nil
}

func (ftdb *NullFullTextDatabase) IndexFeature(ctx context.Context, f wof_geojson.Feature) error {
	return nil
}

func (ftdb *NullFullTextDatabase) QueryString(ctx context.Context, q string, filters ...filter.Filter) (spr.StandardPlacesResults, error) {
	return nil, errors.New("Not implemented")
}
