package auth

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-http-auth"
	"net/http"
)

type GitHubAPIAuthenticator struct {
	auth.Authenticator
}

func NewGitHubAPIAuthenticator(ctx context.Context, uri string) (auth.Authenticator, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (a *GitHubAPIAuthenticator) WrapHandler(next http.Handler) http.Handler {
	return next
}

func (a *GitHubAPIAuthenticator) GetAccountForRequest(*http.Request) (*auth.Account, error) {
	return nil, fmt.Errorf("Not implemented")
}
