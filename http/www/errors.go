// package www implements HTTP handlers for the whosonfirst-browser web application.
package www

import (
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
)

type ErrorVars struct {
	Error error
	URIs  *browser_uris.URIs
}

type NotFoundVars struct {
	URIs *browser_uris.URIs
}
