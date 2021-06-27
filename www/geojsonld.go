package www

import (
	"github.com/sfomuseum/go-geojsonld"
	"github.com/whosonfirst/go-reader"
	"net/http"
)

func GeoJSONLDHandler(r reader.Reader) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := ParseURIFromRequest(req, r)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		f := uri.Feature
		body := f.Bytes()

		ctx := req.Context()

		body, err = geojsonld.AsGeoJSONLD(ctx, body)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		rsp.Header().Set("Content-Type", "application/geo+json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(body)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
