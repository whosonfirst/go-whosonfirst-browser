package writer

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"io"
	"net/url"
	"sort"
	"strings"
)

var writer_roster roster.Roster

type WriterInitializationFunc func(ctx context.Context, uri string) (Writer, error)

type Writer interface {
	Write(context.Context, string, io.ReadSeeker) (int64, error)
	WriterURI(context.Context, string) string
	Close(context.Context) error
}

func NewService(ctx context.Context, uri string) (Writer, error) {

	err := ensureWriterRoster()

	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := parsed.Scheme

	i, err := writer_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(WriterInitializationFunc)
	return init_func(ctx, uri)
}

func RegisterWriter(ctx context.Context, scheme string, init_func WriterInitializationFunc) error {

	err := ensureWriterRoster()

	if err != nil {
		return err
	}

	return writer_roster.Register(ctx, scheme, init_func)
}

func ensureWriterRoster() error {

	if writer_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		writer_roster = r
	}

	return nil
}

func NewWriter(ctx context.Context, uri string) (Writer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := writer_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(WriterInitializationFunc)
	return init_func(ctx, uri)
}

func Writers() []string {
	ctx := context.Background()
	return writer_roster.Drivers(ctx)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureWriterRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range writer_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
