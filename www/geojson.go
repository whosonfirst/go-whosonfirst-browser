package www

import (
	"github.com/whosonfirst/go-reader"
	"log"
	"net/http"
)

func GeoJSONHandler(r reader.Reader) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := ParseURIFromRequest(req, r)

		if err != nil {

			log.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			http.Error(rsp, err.Error(), status)
			return
		}

		f := uri.Feature

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(f.Bytes())
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
