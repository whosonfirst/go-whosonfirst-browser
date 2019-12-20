package options

import (
	"github.com/aaronland/go-artisanal-integers"
	brooklyn_integers "github.com/aaronland/go-brooklynintegers-api"
	"github.com/whosonfirst/go-whosonfirst-export/uid"
)

type DefaultOptions struct {
	Options
	uid_provider uid.Provider
}

func NewDefaultOptions() (Options, error) {

	bi_client := brooklyn_integers.NewAPIClient()
	return NewDefaultOptionsWithArtisanalIntegerClient(bi_client)
}

func NewDefaultOptionsWithArtisanalIntegerClient(client artisanalinteger.Client) (Options, error) {

	provider, err := uid.NewArtisanalUIDProvider(client)

	if err != nil {
		return nil, err
	}

	opts := DefaultOptions{
		uid_provider: provider,
	}

	return &opts, nil
}

func (o *DefaultOptions) UIDProvider() uid.Provider {
	return o.uid_provider
}
