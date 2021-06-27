package www

import (
	"encoding/json"
	"github.com/whosonfirst/go-reader"
	"net/http"
)

func SPRHandler(r reader.Reader) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := ParseURIFromRequest(req, r)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		f := uri.Feature
		s, err := f.SPR()

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(s)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(body)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
