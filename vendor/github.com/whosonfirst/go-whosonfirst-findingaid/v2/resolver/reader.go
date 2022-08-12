package resolver

import (
	"context"
	"fmt"
	"io"
	_ "log"
	"net/url"
	"github.com/whosonfirst/go-reader"
	"github.com/tidwall/gjson"	
	"github.com/whosonfirst/go-whosonfirst-uri"	
)

// type ReaderResolver implements the `Resolver` interface for data that can be resolved using a whosonfirst/go-reader.Reader instance.
type ReaderResolver struct {
	Resolver
	reader reader.Reader
	strategy string
}

func init() {
	ctx := context.Background()
	RegisterResolver(ctx, "reader", NewReaderResolver)
}

// NewReaderResolver will return a new `Resolver` instance for resolving repository names
// that can be resolved using a whosonfirst/go-reader.Reader instance derived from 'uri'.
// 'uri' takes the form of:
//
//	reader://?reader={READER_URI}&strategy={STRATEGY}
//
// Where:
// * {READER_URI} is a valid whosonfirst/go-reader.Reader URI
// * {STRATEGY} is a string describing the strategy use to expand the 'id' parameter passed to the
//   `GetRepo` method to a URI. Valid options are:
//   ** 'fname' which will expand '101736545' to '101736545.geojson' (default)
//   ** 'uri' which will expand '101736545' to '101/736/545/101736545.geojson'
func NewReaderResolver(ctx context.Context, uri string) (Resolver, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	q := u.Query()

	reader_uri := q.Get("reader")

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse reader URI '%s', %w", reader_uri, err)
	}

	s := q.Get("strategy")

	if s == "" {
		s = "fname"
	}
	
	f := &ReaderResolver{
		reader: r,
		strategy: s,
	}

	return f, nil
}

// GetRepo returns the name of the repository associated with this ID.
func (r *ReaderResolver) GetRepo(ctx context.Context, id int64) (string, error) {

	var path string

	switch r.strategy {
	case "uri":
		
		rel_path, err := uri.Id2RelPath(id)
		
		if err != nil {
			return "", fmt.Errorf("Failed to derive rel path, %w", err)
		}

		path = rel_path

	case "fname":
		
		fname, err := uri.Id2Fname(id)

		if err != nil {
			return "", fmt.Errorf("Failed to derive filename, %w", err)
		}

		path = fname
		
	default:

		return "", fmt.Errorf("Invalid strategy")
	}
	
	fh, err := r.reader.Read(ctx, path)

	if err != nil {
		return "", fmt.Errorf("Failed to read %s, %w", path, err)
	}
	
	defer fh.Close()

	body, err := io.ReadAll(fh)

	if err != nil {
		return "", fmt.Errorf("Failed to read body, %w", err)
	}

	repo_rsp := gjson.GetBytes(body, "properties.wof:repo")

	if !repo_rsp.Exists(){
		return "", fmt.Errorf("Body missing wof:repo property")
	}

	return repo_rsp.String(), nil
}
