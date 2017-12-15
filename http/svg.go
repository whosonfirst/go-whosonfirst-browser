package http

import (
	"github.com/whosonfirst/go-whosonfirst-static/reader"
	"github.com/whosonfirst/go-whosonfirst-static/utils"
	"github.com/whosonfirst/go-whosonfirst-svg"
	gohttp "net/http"
)

func SVGHandler(r reader.Reader) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		f, err, status := utils.FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Content-Type", "image/svg+xml")

		opts := svg.NewDefaultOptions()
		opts.Writer = rsp

		svg.FeatureToSVG(f, opts)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
