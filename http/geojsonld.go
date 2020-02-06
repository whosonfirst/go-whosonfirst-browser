package http

import (
	"github.com/whosonfirst/go-reader"
	"github.com/sfomuseum/go-geojsonld"
	gohttp "net/http"
)

func GeoJSONLDHandler(r reader.Reader) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		uri, err, status := ParseURIFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		f := uri.Feature
		body := f.Bytes()

		ctx := req.Context()
		
		body, err = geojsonld.AsGeoJSONLD(ctx, body)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}
		
		rsp.Header().Set("Content-Type", "application/geo+json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(body)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
