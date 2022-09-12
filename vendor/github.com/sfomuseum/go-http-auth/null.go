package auth

import (
	"context"
	"net/http"
)

func init() {
	ctx := context.Background()
	RegisterAuthenticator(ctx, "null", NewNullAuthenticator)
}

// type NullAuthenticator implements the Authenticator interface such that no authentication is performed.
type NullAuthenticator struct {
	Authenticator
}

// NewNullAuthenticator implements the Authenticator interface such that no authentication is performed
// configured by 'uri' which is expected to take the form of:
//
//	null://
func NewNullAuthenticator(ctx context.Context, uri string) (Authenticator, error) {
	a := &NullAuthenticator{}
	return a, nil
}

// WrapHandler returns 'h' unchanged.
func (a *NullAuthenticator) WrapHandler(h http.Handler) http.Handler {
	return h
}

// GetAccountForRequest returns a `NotLoggedIn` error.
func (a *NullAuthenticator) GetAccountForRequest(req *http.Request) (*Account, error) {

	acct := &Account{
		Id:   0,
		Name: "Null",
	}

	return acct, nil
}
