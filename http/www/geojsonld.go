package www

import (
	"log"
	"net/http"

	"github.com/sfomuseum/go-geojsonld"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v7/http"
)

type GeoJSONLDHandlerOptions struct {
	Reader reader.Reader
	Logger *log.Logger
}

func GeoJSONLDHandler(opts *GeoJSONLDHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {

			opts.Logger.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			http.Error(rsp, err.Error(), status)
			return
		}

		ctx := req.Context()

		body, err := geojsonld.AsGeoJSONLD(ctx, uri.Feature)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		rsp.Header().Set("Content-Type", "application/geo+json")
		rsp.Write(body)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
