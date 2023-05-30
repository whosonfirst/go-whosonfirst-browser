// package www implements HTTP handlers for the whosonfirst-browser web application.
package www

import (
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7"
)

type ErrorVars struct {
	Error error
	URIs  *browser.URIs
}

type NotFoundVars struct {
	URIs *browser.URIs
}
