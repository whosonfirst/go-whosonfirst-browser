package utils

// this should go in the general-purpose go-whosonfirst-reader package
// that doesn't exist yet... (20171214/thisisaaronland)

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-render/reader"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"net/http"
)

func FeatureFromRequest(req *http.Request, r reader.Reader) (geojson.Feature, error, int) {

	path := req.URL.Path

	wofid, err := uri.IdFromPath(path)

	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	rel_path, err := uri.Id2RelPath(wofid)

	if err != nil {
		return nil, err, http.StatusBadRequest // StatusInternalServerError
	}

	fh, err := r.Read(rel_path)

	if err != nil {
		return nil, err, http.StatusBadRequest // StatusInternalServerError
	}

	f, err := feature.LoadFeatureFromReader(fh)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return f, nil, 0
}
