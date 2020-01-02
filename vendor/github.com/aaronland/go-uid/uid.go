package uid

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
)

var providers roster.Roster

func ensureRoster() error {

	if providers == nil {
		
		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		providers = r
	}

	return nil
}

func RegisterProvider(ctx context.Context, name string, pr Provider) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return providers.Register(ctx, name, pr)
}

func NewProvider(ctx context.Context, uri string) (Provider, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := providers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	pr := i.(Provider)
	
	err = pr.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return pr, nil
}

type Provider interface {
	Open(context.Context, string) error
	UID(...interface{}) (UID, error) // NOT SURE ABOUT THIS...
}

type UID interface {
	String() string
}
