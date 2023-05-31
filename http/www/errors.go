// package www implements HTTP handlers for the whosonfirst-browser web application.
package www

import (
	"github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
)

type ErrorVars struct {
	Error error
	URIs  *uris.URIs
}

type NotFoundVars struct {
	URIs *uris.URIs
}
