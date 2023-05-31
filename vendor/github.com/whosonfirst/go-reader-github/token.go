package reader

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sfomuseum/runtimevar"
)

// EnsureGitHubAccessToken ensures that 'writer_uri' contains a '?access_token=VALUE' parameter. This only
// applies if the scheme is `githubapi://`. If not the method returns the original 'target_uri' value.
// If the 'target_uri' contains an empty `access_token` parameter or the value is "{access_token}"
// then the method will replace parameter with the value derived from 'token_uri' which is expected
// to be a valid `gocloud.dev/runtimevar` URI.
func EnsureGitHubAccessToken(ctx context.Context, target_uri string, token_uri string) (string, error) {

	u, err := url.Parse(target_uri)

	if err != nil {
		return "", fmt.Errorf("Failed to parse writer URI, %w", err)
	}

	switch u.Scheme {
	case "githubapi":

		q := u.Query()

		token := q.Get("access_token")

		switch token {
		case "", "{access_token}":
			// continue
		default:
			return target_uri, nil
		}

		token, err = runtimevar.StringVar(ctx, token_uri)

		if err != nil {
			return "", fmt.Errorf("Failed to expand token URI, %w", err)
		}

		q.Del("access_token")
		q.Set("access_token", token)

		u.RawQuery = q.Encode()

	case "findingaid":

		q := u.Query()

		template_uri := q.Get("template")

		template_u, err := url.Parse(template_uri)

		if err != nil {
			return "", fmt.Errorf("Failed to parse template URI from findingaid, %w", err)
		}

		template_q := template_u.Query()

		token := template_q.Get("access_token")

		switch token {
		case "", "{access_token}":
			// continue
		default:
			return target_uri, nil
		}

		token, err = runtimevar.StringVar(ctx, token_uri)

		if err != nil {
			return "", fmt.Errorf("Failed to expand token URI, %w", err)
		}

		template_q.Del("access_token")
		template_q.Set("access_token", token)

		template_u.RawQuery = template_q.Encode()

		q.Set("template", template_u.String())
		u.RawQuery = q.Encode()

	default:

		return target_uri, nil
	}

	return u.String(), nil
}
