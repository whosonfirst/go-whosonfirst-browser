package uid

import (
	"github.com/aaronland/go-artisanal-integers"
)

type ArtisanalUIDProvider struct {
	Provider
	client       artisanalinteger.Client
	max_attempts int
}

func NewArtisanalUIDProvider(client artisanalinteger.Client) (Provider, error) {

	p := ArtisanalUIDProvider{
		client:       client,
		max_attempts: 5,
	}

	return &p, nil
}

func (p *ArtisanalUIDProvider) UID() (int64, error) {

	attempts := p.max_attempts

	var i int64
	var err error

	for attempts > 0 {

		i, err = p.client.NextInt()

		if err == nil {
			break
		} else {
			attempts -= 1
		}
	}

	return i, err
}
