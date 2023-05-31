package browser

import (
	"context"
	"sync"

	"github.com/sfomuseum/go-http-auth"
)

var setupAuthenticatorOnce sync.Once
var setupAuthenticatorError error

func setupAuthenticator() {

	var err error
	ctx := context.Background()

	authenticator, err = auth.NewAuthenticator(ctx, cfg.AuthenticatorURI)

	if err != nil {
		setupAuthenticatorError = err
		return
	}
}
