package writer

import (
	"context"
	"fmt"
	"github.com/sfomuseum/runtimevar"
	"net/url"
	"strings"
)

// EnsureGitHubAccessToken ensures that 'writer_uri' contains a '?access_token=VALUE' parameter. This only
// applies if the scheme is `githubapi://`. If not the method returns the original 'writer_uri' value.
// If the 'writer_uri' contains an empty `access_token` parameter or the value is "{access_token}"
// then the method will replace parameter with the value derived from 'token_uri' which is expected
// to be a valid `gocloud.dev/runtimevar` URI.
func EnsureGitHubAccessToken(ctx context.Context, writer_uri string, token_uri string) (string, error) {

	u, err := url.Parse(writer_uri)

	if err != nil {
		return "", fmt.Errorf("Failed to parse writer URI, %w", err)
	}

	if !strings.HasPrefix(u.Scheme, "githubapi") {
		return writer_uri, nil
	}

	q := u.Query()

	token := q.Get("access_token")

	switch token {
	case "", "{access_token}":
		// continue
	default:
		return writer_uri, nil
	}

	token, err = runtimevar.StringVar(ctx, token_uri)

	if err != nil {
		return "", fmt.Errorf("Failed to expand token URI, %w", err)
	}

	q.Del("access_token")
	q.Set("access_token", token)

	u.RawQuery = q.Encode()

	return u.String(), nil
}
