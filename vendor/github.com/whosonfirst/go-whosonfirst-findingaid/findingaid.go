package findingaid

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
)

type FindingAid interface {
	Resolver
	Indexer
}

type FindingAidInitializationFunc func(ctx context.Context, uri string) (FindingAid, error)

var findingaids roster.Roster

func ensureFindingAids() error {

	if findingaids == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		findingaids = r
	}

	return nil
}

func RegisterFindingAid(ctx context.Context, name string, c FindingAidInitializationFunc) error {

	err := ensureFindingAids()

	if err != nil {
		return err
	}

	return findingaids.Register(ctx, name, c)
}

func NewFindingAid(ctx context.Context, uri string) (FindingAid, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := findingaids.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init := i.(FindingAidInitializationFunc)
	c, err := init(ctx, uri)

	if err != nil {
		return nil, err
	}

	return c, nil
}
