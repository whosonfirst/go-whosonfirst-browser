package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jtacoma/uritemplates"
	"github.com/whosonfirst/go-whosonfirst-findingaid"
	"github.com/whosonfirst/go-whosonfirst-uri"
	_ "log"
	"net/http"
	"net/url"
)

const FINDINGAID_URI_TEMPLATE string = "https://data.whosonfirst.org/findingaid/{id}"

// RepoResolver is a struct that implements the findingaid.Resolver interface for information about Who's On First repositories by retrieving information from an HTTP endpoint that returns JSON-encoded FindingAidResponse responses. For example, a remote server running the `application/lookupd` tool.
type HTTPResolver struct {
	findingaid.Resolver
	template *uritemplates.UriTemplate
}

func init() {

	ctx := context.Background()
	err := findingaid.RegisterResolver(ctx, "repo-http", NewHTTPResolver)

	if err != nil {
		panic(err)
	}
}

// NewRepoResolver returns a findingaid.Resolver instance for exposing information about Who's On First repositories by retrieving information from an HTTP endpoint that returns JSON-encoded FindingAidResponse responses.
func NewHTTPResolver(ctx context.Context, uri string) (findingaid.Resolver, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	fa_template_uri := q.Get("findingaid_uri_template")

	if fa_template_uri == "" {
		fa_template_uri = FINDINGAID_URI_TEMPLATE
	}

	fa_template, err := uritemplates.Parse(fa_template_uri)

	if err != nil {
		return nil, err
	}

	fa := &HTTPResolver{
		template: fa_template,
	}

	return fa, nil
}

// ResolveURI will return 'repo.FindingAidResponse' for 'str_response' if it present in the finding aid.
func (fa *HTTPResolver) ResolveURI(ctx context.Context, str_uri string) (interface{}, error) {

	id, _, err := uri.ParseURI(str_uri)

	if err != nil {
		return nil, err
	}

	values := map[string]interface{}{
		"id": id,
	}

	uri, err := fa.template.Expand(values)

	if err != nil {
		return nil, err
	}

	rsp, err := http.Get(uri)

	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to resolve, %s", rsp.Status)
	}

	var fa_rsp *FindingAidResponse

	dec := json.NewDecoder(rsp.Body)
	err = dec.Decode(&fa_rsp)

	if err != nil {
		return nil, err
	}

	return fa_rsp, nil
}
