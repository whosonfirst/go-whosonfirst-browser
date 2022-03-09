package www

// https://preview.iiif.io/api/navplace_extension/api/extension/navplace/

import (
	"github.com/whosonfirst/go-reader"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type NavPlaceHandlerOptions struct {
	Reader      reader.Reader
	MaxFeatures int
}

// NavPlaceHandler will return a given record as a FeatureCollection for use by the IIIF navPlace extension,
// specifically as navPlace "reference" objects.
func NavPlaceHandler(opts *NavPlaceHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		q := req.URL.Query()
		base := q.Get("id")

		if base == "" {
			path := req.URL.Path
			base = filepath.Base(path)

			base = strings.TrimLeft(base, "/")
			base = strings.TrimRight(base, "/")
		}

		ids := strings.Split(base, ",")

		uris := make([]*URI, len(ids))

		for idx, id := range ids {

			uri, err, status := ParseURIFromPath(ctx, id, opts.Reader)

			if err != nil {

				log.Printf("Failed to parse URI from request %s, %v", req.URL, err)

				http.Error(rsp, err.Error(), status)
				return
			}

			uris[idx] = uri
		}

		count := len(uris)

		if count == 0 {
			http.Error(rsp, "No IDs to include", http.StatusBadRequest)
			return
		}

		if count > opts.MaxFeatures {
			http.Error(rsp, "Maximum number of IDs exceeded", http.StatusBadRequest)
			return
		}

		rsp.Header().Set("Content-Type", "application/geo+json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write([]byte(`{"type":"FeatureCollection", "features":[`))

		for i, uri := range uris {

			f := uri.Feature
			body := f.Bytes()

			rsp.Write(body)

			if i+1 < count {
				rsp.Write([]byte(`,`))
			}
		}

		rsp.Write([]byte(`]}`))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
