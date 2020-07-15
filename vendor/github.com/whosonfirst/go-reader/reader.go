package reader

import (
	"context"
	"github.com/aaronland/go-roster"
	"io"
	"net/url"
)

var reader_roster roster.Roster

type ReaderInitializationFunc func(ctx context.Context, uri string) (Reader, error)

type Reader interface {
	Read(context.Context, string) (io.ReadCloser, error)
	URI(string) string
}

func NewService(ctx context.Context, uri string) (Reader, error) {

	err := ensureReaderRoster()

	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := parsed.Scheme

	i, err := reader_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ReaderInitializationFunc)
	return init_func(ctx, uri)
}

func RegisterReader(ctx context.Context, scheme string, init_func ReaderInitializationFunc) error {

	err := ensureReaderRoster()

	if err != nil {
		return err
	}

	return reader_roster.Register(ctx, scheme, init_func)
}

func ensureReaderRoster() error {

	if reader_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		reader_roster = r
	}

	return nil
}

func NewReader(ctx context.Context, uri string) (Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := reader_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ReaderInitializationFunc)
	return init_func(ctx, uri)
}

func Readers() []string {
	ctx := context.Background()
	return reader_roster.Drivers(ctx)
}
