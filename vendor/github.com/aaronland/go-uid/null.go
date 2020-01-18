package uid

import (
	"context"
)

func init() {
	ctx := context.Background()
	pr := NewNullProvider()
	RegisterProvider(ctx, "null", pr)
}

type NullProvider struct {
	Provider
}

type NullUID struct {
	UID
}

func NewNullProvider() Provider {
	pr := &NullProvider{}
	return pr
}

func (pr *NullProvider) Open(ctx context.Context, uri string) error {
	return nil
}

func (n *NullProvider) UID(...interface{}) (UID, error) {
	return NewNullUID()
}

func NewNullUID() (UID, error) {
	n := &NullUID{}
	return n, nil
}

func (n *NullUID) String() string {
	return ""
}
