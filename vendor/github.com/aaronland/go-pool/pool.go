package pool

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
	"strconv"
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

func Register(ctx context.Context, name string, pr Pool) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return providers.Register(ctx, name, pr)
}

type Pool interface {
	Open(context.Context, string) error
	Length() int64
	Push(Item)
	Pop() (Item, bool)
}

type Item interface {
	String() string
	Int() int64
}

type Int struct {
	Item
	int int64
}

type String struct {
	Item
	string string
}

func NewIntItem(i int64) Item {
	return &Int{int: i}
}

func NewStringItem(s string) Item {
	return &String{string: s}
}

func (i Int) String() string {
	return strconv.FormatInt(i.int, 10)
}

func (i Int) Int() int64 {
	return i.int
}

func (s String) String() string {
	return s.string
}

func (s String) Int() int64 {
	return int64(0)
}

func NewPool(ctx context.Context, uri string) (Pool, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := providers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	pr := i.(Pool)

	err = pr.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return pr, nil
}
