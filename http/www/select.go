package www

import (
	"encoding/json"
	"github.com/aaronland/go-http-sanitize"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v6/http"
	"log"
	"net/http"
	"regexp"
)

type SelectHandlerOptions struct {
	Pattern *regexp.Regexp
	Reader  reader.Reader
	Logger  *log.Logger
}

func SelectHandler(opts *SelectHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		query, err := sanitize.GetString(req, "select")

		if err != nil {
			http.Error(rsp, "Invalid parameters", http.StatusBadRequest)
			return
		}

		if query == "" {
			http.Error(rsp, "Missing select", http.StatusBadRequest)
			return
		}

		if !opts.Pattern.MatchString(query) {
			http.Error(rsp, "Invalid select", http.StatusBadRequest)
			return
		}

		uri, err, status := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {

			opts.Logger.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			http.Error(rsp, err.Error(), status)
			return
		}

		query_rsp := gjson.GetBytes(uri.Feature, query)

		var rsp_body []byte

		if query_rsp.Exists() {

			enc, err := json.Marshal(query_rsp.Value())

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			rsp_body = enc
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(rsp_body)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
