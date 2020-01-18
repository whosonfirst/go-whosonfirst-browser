package artisanalinteger

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
)

var clients roster.Roster

func ensureClients() error {

	if clients == nil {
		
		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		clients = r
	}

	return nil
}

func RegisterClient(ctx context.Context, name string, cl Client) error {

	err := ensureClients()

	if err != nil {
		return err
	}

	return clients.Register(ctx, name, cl)
}

func NewClient(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := clients.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	cl := i.(Client)
	return cl, nil
}

type Engine interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
	SetLastInt(int64) error
	SetKey(string) error
	SetOffset(int64) error
	SetIncrement(int64) error
	Close() error
}

type Service interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
}

type Server interface {
	ListenAndServe(Service) error
	Address() string
}

type Client interface {
	NextInt() (int64, error)
}

type Integer struct {
	Integer int64 `json:"integer"`
}
