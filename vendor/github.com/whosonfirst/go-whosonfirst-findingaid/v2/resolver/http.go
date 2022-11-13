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
	endpoint   string
	user_agent string
}

func init() {
	ctx := context.Background()
	RegisterResolver(ctx, "http", NewHTTPResolver)
	RegisterResolver(ctx, "https", NewHTTPResolver)
}

// NewHTTPResolver will return a new `Resolver` instance for resolving repository names
// and IDs stored in a database exposed via an HTTP endpoint..
func NewHTTPResolver(ctx context.Context, uri string) (Resolver, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	r := &HTTPResolver{
		endpoint: uri,
	}

	q := u.Query()
	ua := q.Get("user-agent")

	if ua != "" {
		r.user_agent = ua
	}

	return r, nil
}

// GetRepo returns the name of the repository associated with this ID in a Who's On First finding aid.
func (r *HTTPResolver) GetRepo(ctx context.Context, id int64) (string, error) {

	u, err := url.Parse(r.endpoint)

	if err != nil {
		return "", fmt.Errorf("Failed to parse endpoint URL, %w", err)
	}

	str_id := strconv.FormatInt(id, 10)
	u.Path = filepath.Join(u.Path, str_id)

	url := u.String()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", fmt.Errorf("Failed to create new request, %w", err)
	}

	if r.user_agent != "" {
		req.Header.Set("User-Agent", r.user_agent)
	}

	cl := &http.Client{}
	rsp, err := cl.Do(req)

	if err != nil {
		return "", fmt.Errorf("Failed to execute request, %w", err)
	}

	if rsp.StatusCode != 200 {
		return "", fmt.Errorf("Unexpected status code: %s", rsp.Status)
	}

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		return "", fmt.Errorf("Failed to read body, %w", err)
	}

	return string(body), nil
}
