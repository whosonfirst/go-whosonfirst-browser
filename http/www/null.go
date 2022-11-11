package www

import (
	"net/http"
)

func NewNullHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		return
	}

	return http.HandlerFunc(fn)
}
