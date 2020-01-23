package http

import (
	"encoding/json"
	"github.com/aaronland/go-http-sanitize"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	gohttp "net/http"
	"regexp"
)

type SelectHandlerOptions struct {
	Pattern *regexp.Regexp
}

func SelectHandler(r reader.Reader, opts *SelectHandlerOptions) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		query, err := sanitize.GetString(req, "select")

		if err != nil {
			gohttp.Error(rsp, "Invalid parameters", gohttp.StatusBadRequest)
			return
		}

		if query == "" {
			gohttp.Error(rsp, "Missing select", gohttp.StatusBadRequest)
			return
		}

		if !opts.Pattern.MatchString(query) {
			gohttp.Error(rsp, "Invalid select", gohttp.StatusBadRequest)
			return
		}

		uri, err, status := ParseURIFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		f := uri.Feature

		query_rsp := gjson.GetBytes(f.Bytes(), query)

		var rsp_body []byte

		if query_rsp.Exists() {

			enc, err := json.Marshal(query_rsp.Value())

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			rsp_body = enc
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(rsp_body)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
