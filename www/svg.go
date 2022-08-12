package www

import (
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-sanitize"
	"github.com/whosonfirst/go-whosonfirst-svg"
	"log"
	"net/http"
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

func SVGHandler(r reader.Reader, handler_opts *SVGOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := ParseURIFromRequest(req, r)

		if err != nil {

			log.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			http.Error(rsp, err.Error(), status)
			return
		}

		f := uri.Feature

		sn_opts := sanitize.DefaultOptions()

		sz := "lg"

		query := req.URL.Query()
		query_sz := query.Get("size")

		req_sz, err := sanitize.SanitizeString(query_sz, sn_opts)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		if req_sz != "" {
			sz = req_sz
		}

		sz_info, ok := handler_opts.Sizes[sz]

		if !ok {
			http.Error(rsp, "Invalid output size", http.StatusBadRequest)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Content-Type", "image/svg+xml")

		opts := svg.NewDefaultOptions()
		opts.Height = float64(sz_info.MaxHeight)
		opts.Width = float64(sz_info.MaxWidth)
		opts.Writer = rsp

		// to do: support for custom styles:
		// https://github.com/whosonfirst/go-whosonfirst-browser/issues/19

		opts.StyleFunction = func(f []byte) (map[string]string, error) {

			attrs := make(map[string]string)

			type_rsp := gjson.GetBytes(f, "geometry.type")

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

	h := http.HandlerFunc(fn)
	return h, nil
}
