package uid

import (
	"context"
	"github.com/aaronland/go-string/random"
)

func init() {
	ctx := context.Background()	
	pr := NewRandomProvider()
	RegisterProvider(ctx, "random", pr)
}

type RandomProvider struct {
	Provider
}

// type RandomUID is a type StringUID

func NewRandomProvider() Provider {
	pr := &RandomProvider{}
	return pr
}

func (pr *RandomProvider) Open(ctx context.Context, uri string) error {
	return nil
}

func (pr *RandomProvider) UID(...interface{}) (UID, error) {

	opts := random.DefaultOptions()
	opts.Length = 16
	opts.Chars = 16
	opts.ASCII = true
	opts.AlphaNumeric = true

	s, err := random.String(opts)

	if err != nil {
		return nil, err
	}

	return NewStringUID(s)
}
