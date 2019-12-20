package http

import (
	"github.com/whosonfirst/go-reader"
	gohttp "net/http"
)

func GeoJSONHandler(r reader.Reader) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		f, err, status := FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		WriteGeoJSONHeaders(rsp)

		rsp.Write(f.Bytes())
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
