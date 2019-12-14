package http

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"html/template"
	_ "log"
	gohttp "net/http"
	"time"
)

type HTMLHandlerOptions struct {
	DataEndpoint string
	Templates    *template.Template
}

type HTMLVars struct {
	SPR          spr.StandardPlacesResult
	LastModified string
	DataEndpoint string
}

func HTMLHandler(r reader.Reader, opts HTMLHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("id")

	if t == nil {
		return nil, errors.New("Missing id template")
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		f, err, status := FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		s, err := f.SPR()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		now := time.Now()
		lastmod := now.Format(time.RFC3339)

		// this assumes the GeoJSONHandler being assigned to /data
		// (20180419/thisisaaronland)

		data_endpoint := fmt.Sprintf("%s//%s/data/", req.URL.Scheme, req.Host)

		if opts.DataEndpoint != "" {
			data_endpoint = opts.DataEndpoint
		}

		vars := HTMLVars{
			SPR:          s,
			LastModified: lastmod,
			DataEndpoint: data_endpoint,
		}

		err = t.Execute(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
