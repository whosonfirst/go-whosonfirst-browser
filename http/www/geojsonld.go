package www

import (
	"github.com/sfomuseum/go-geojsonld"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v5/http"
	"log"
	"net/http"
)

func GeoJSONLDHandler(r reader.Reader) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := wof_http.ParseURIFromRequest(req, r)

		if err != nil {

			log.Printf("Failed to parse URI from request %s, %v", req.URL, err)

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
