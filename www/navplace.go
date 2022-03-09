package www

// https://preview.iiif.io/api/navplace_extension/api/extension/navplace/

import (
	"github.com/whosonfirst/go-reader"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// NavPlaceHandler will return a given record as a FeatureCollection for use by the IIIF navPlace extension,
// specifically as navPlace "reference" objects.
func NavPlaceHandler(r reader.Reader) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		
		path := req.URL.Path

		base := filepath.Base(path)
		base = strings.TrimLeft(base, "/")
		base = strings.TrimRight(base, "/")

		ids := strings.Split(base, ",")

		uris := make([]*URI, len(ids))
		
		for idx, id := range ids {
		
			uri, err, status := ParseURIFromPath(ctx, id, r)

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
		
		rsp.Header().Set("Content-Type", "application/geo+json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write([]byte(`{"type":"FeatureCollection", "features":[`))

		for i, uri := range uris {

			f := uri.Feature
			body := f.Bytes()
			
			rsp.Write(body)

			if i+1 < count {
				rsp.Write([]byte(`,\n`))
			}
		}
		
		rsp.Write([]byte(`]}`))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
