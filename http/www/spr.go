package www

import (
	"encoding/json"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v7/http"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"log"
	"net/http"
)

type SPRHandlerOptions struct {
	Reader reader.Reader
	Logger *log.Logger
}

func SPRHandler(opts *SPRHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {

			opts.Logger.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			http.Error(rsp, err.Error(), status)
			return
		}

		s, err := spr.WhosOnFirstSPR(uri.Feature)

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
		rsp.Write(body)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
