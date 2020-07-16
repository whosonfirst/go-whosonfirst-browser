package options

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-id"
)

type DefaultOptions struct {
	Options
	id_provider id.Provider
}

func NewDefaultOptions() (Options, error) {

	provider, err := id.NewProvider(context.Background())

	if err != nil {
		return nil, err
	}

	return NewDefaultOptionsWithProvider(provider)
}

func NewDefaultOptionsWithProvider(provider id.Provider) (Options, error) {

	opts := DefaultOptions{
		id_provider: provider,
	}

	return &opts, nil	
}

func (opts *DefaultOptions) IDProvider() id.Provider {
	return opts.id_provider
}
