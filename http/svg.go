package http

import (
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-sanitize"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-svg"
	"log"
	gohttp "net/http"
)

type SVGSize struct {
	Label     string
	MaxHeight int
	MaxWidth  int
}

type SVGOptions struct {
	Sizes map[string]SVGSize
}

func NewDefaultSVGOptions() (*SVGOptions, error) {

	sm := SVGSize{
		Label:     "sm",
		MaxHeight: 300,
		MaxWidth:  300,
	}

	med := SVGSize{
		Label:     "med",
		MaxHeight: 640,
		MaxWidth:  640,
	}

	lg := SVGSize{
		Label:     "lg",
		MaxHeight: 1024,
		MaxWidth:  1024,
	}

	sz := map[string]SVGSize{
		"sm":  sm,
		"med": med,
		"lg":  lg,
	}

	opts := SVGOptions{
		Sizes: sz,
	}

	return &opts, nil
}

func SVGHandler(r reader.Reader, handler_opts *SVGOptions) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		foo, err, status := FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		f := foo.Feature

		sn_opts := sanitize.DefaultOptions()

		sz := "lg"

		query := req.URL.Query()
		query_sz := query.Get("size")

		req_sz, err := sanitize.SanitizeString(query_sz, sn_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		if req_sz != "" {
			sz = req_sz
		}

		sz_info, ok := handler_opts.Sizes[sz]

		if !ok {
			gohttp.Error(rsp, "Invalid output size", gohttp.StatusBadRequest)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Content-Type", "image/svg+xml")

		opts := svg.NewDefaultOptions()
		opts.Height = float64(sz_info.MaxHeight)
		opts.Width = float64(sz_info.MaxWidth)
		opts.Writer = rsp

		opts.StyleFunction = func(f geojson.Feature) (map[string]string, error) {

			attrs := make(map[string]string)

			type_rsp := gjson.GetBytes(f.Bytes(), "geometry.type")

			if !type_rsp.Exists() {
				return nil, errors.New("Missing geometry.type")
			}

			geom_type := type_rsp.String()
			// log.Println(geom_type)

			switch geom_type {
			case "LineString":
				attrs["fill-opacity"] = "0.0"
				attrs["stroke-width"] = "1.0"
				attrs["stroke-opacity"] = "2.0"
				attrs["stroke"] = "#000"
			case "Point", "MultiPoint":
				// something something something
				// https://github.com/whosonfirst/go-whosonfirst-browser/issues/18
			default:
				// pass
			}

			return attrs, nil
		}

		svg.FeatureToSVG(f, opts)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
