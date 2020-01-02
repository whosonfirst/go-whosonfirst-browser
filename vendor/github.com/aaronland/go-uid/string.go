package uid

import (
	"context"
	"errors"
	"net/url"
)

func init() {
	ctx := context.Background()
	pr := NewStringProvider()
	RegisterProvider(ctx, "string", pr)
}

type StringProvider struct {
	Provider
	string string
}

type StringUID struct {
	UID
	string string
}

func NewStringProvider() Provider {
	pr := &StringProvider{}
	return pr
}

func (pr *StringProvider) Open(ctx context.Context, uri string) error {

	u, err := url.Parse(uri)

	if err != nil {
		return err
	}

	q := u.Query()
	s := q.Get("string")

	if s == "" {
		return errors.New("Empty string")
	}

	pr.string = s
	return nil
}

func (pr *StringProvider) UID(...interface{}) (UID, error) {
	return NewStringUID(pr.string)
}

func NewStringUID(s string) (UID, error) {

	u := StringUID{
		string: s,
	}

	return &u, nil
}

// where is UID() ?

func (u *StringUID) String() string {
	return u.string
}
