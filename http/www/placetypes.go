package www

import (
	"io"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

// PlacetypesHandler returns a `http.Handler` that serves a JSON representation of the
// `go-whosonfirst-placetypes.DefaultWOFPlacetypeSpecification` struct. This is meant to
// be consumed by the web component that displays a list of valid wof:placetype values.
func PlacetypesHandler() (http.Handler, error) {
	
	fn := func(rsp http.ResponseWriter, req *http.Request) {

		r, err := placetypes.FS.Open("placetypes.json")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Close()
		
		rsp.Header().Set("Content-type", "text/javascript")

		_, err = io.Copy(rsp, r)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
