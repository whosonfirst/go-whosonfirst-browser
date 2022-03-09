package resolver

import (
	"context"
	"fmt"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

// type HTTPResolver implements the `Resolver` interface for data stored in a database exposed via an HTTP endpoint.
type HTTPResolver struct {
	Resolver
	endpoint string
}

func init() {
	ctx := context.Background()
	RegisterResolver(ctx, "http", NewHTTPResolver)
	RegisterResolver(ctx, "https", NewHTTPResolver)
}

// NewHTTPResolver will return a new `Resolver` instance for resolving repository names
// and IDs stored in a database exposed via an HTTP endpoint..
func NewHTTPResolver(ctx context.Context, uri string) (Resolver, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	f := &HTTPResolver{
		endpoint: uri,
	}

	return f, nil
}

// GetRepo returns the name of the repository associated with this ID in a Who's On First finding aid.
func (r *HTTPResolver) GetRepo(ctx context.Context, id int64) (string, error) {

	u, err := url.Parse(r.endpoint)

	if err != nil {
		return "", fmt.Errorf("Failed to parse endpoint URL, %w", err)
	}
	
	str_id := strconv.FormatInt(id, 10)
	u.Path = filepath.Join(u.Path, str_id)
	
	rsp, err := http.Get(u.String())

	if err != nil {
		return "", fmt.Errorf("Failed to request %s, %w", u.String(), err)
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		return "", fmt.Errorf("Failed to read body, %w", err)
	}

	return string(body), nil
}
