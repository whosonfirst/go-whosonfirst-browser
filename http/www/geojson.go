package www

import (
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v7/http"
	"log"
	"net/http"
)

type GeoJSONHandlerOptions struct {
	Reader reader.Reader
	Logger *log.Logger
}

func GeoJSONHandler(opts *GeoJSONHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {

			opts.Logger.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			http.Error(rsp, err.Error(), status)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(uri.Feature)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
