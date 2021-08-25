package www

// https://preview.iiif.io/api/navplace_extension/api/extension/navplace/

import (
	"github.com/whosonfirst/go-reader"
	"net/http"
)

// NavPlaceHandler will return a given record as a FeatureCollection for use by the IIIF navPlace extension,
// specifically as navPlace "reference" objects.
func NavPlaceHandler(r reader.Reader) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := ParseURIFromRequest(req, r)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		f := uri.Feature
		body := f.Bytes()

		rsp.Header().Set("Content-Type", "application/geo+json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write([]byte(`{"type":"FeatureCollection", "features":[`))
		rsp.Write(body)
		rsp.Write([]byte(`]}`))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
