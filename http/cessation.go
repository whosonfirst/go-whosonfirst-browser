package http

import (
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/update"
	"github.com/whosonfirst/go-writer"
	_ "log"
	gohttp "net/http"
	"regexp"
)

type CessationHandlerOptions struct {
	AllowedPaths *regexp.Regexp // multiple regexps?
}

func CessationHandler(r reader.Reader, wr writer.Writer, opts *CessationHandlerOptions) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		switch req.Method {
		case "POST":
			// pass
		default:
			gohttp.Error(rsp, "Method not allowed.", gohttp.StatusMethodNotAllowed)
			return
		}

		f, err, status := FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		ctx := req.Context()
		body := f.Bytes()

		updated_body, err := update.CessateFeature(ctx, body, opts.AllowedPaths)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		exported_body, err := update.ExportFeature(ctx, wr, updated_body)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		WriteGeoJSONHeaders(rsp)

		rsp.Write(exported_body)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
