package http

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/warning"
	gohttp "net/http"
	"path/filepath"
	"regexp"
	"strconv"
)

func FeatureFromRequest(req *gohttp.Request, r reader.Reader) (geojson.Feature, error, int) {

	path := req.URL.Path

	wofid, err := IdFromPath(path)

	if err != nil {
		return nil, err, gohttp.StatusNotFound
	}

	rel_path, err := uri.Id2RelPath(wofid)

	if err != nil {
		return nil, err, gohttp.StatusBadRequest // StatusInternalServerError
	}

	fh, err := r.Read(rel_path)

	if err != nil {
		return nil, err, gohttp.StatusBadRequest // StatusInternalServerError
	}

	f, err := feature.LoadFeatureFromReader(fh)

	if err != nil && !warning.IsWarning(err) {
		return nil, err, gohttp.StatusInternalServerError
	}

	return f, nil, 0
}
