package http

import (
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/warning"
	_ "log"
	gohttp "net/http"
	"strconv"
)

func IdFromPath(path string) (int64, error) {

	wofid, _, err := IdFromURI(path)
	return wofid, err
}

func FeatureFromRequest(req *gohttp.Request, r reader.Reader) (geojson.Feature, error, int) {

	path := req.URL.Path

	wofid, uri_args, err := IdFromURI(path)

	if err != nil {

		q := req.URL.Query()
		str_id := q.Get("id")

		if str_id == "" {
			return nil, err, gohttp.StatusNotFound
		}

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			return nil, err, gohttp.StatusBadRequest
		}

		wofid = id
	}

	rel_path, err := uri.Id2RelPath(wofid, uri_args)

	if err != nil {
		return nil, err, gohttp.StatusBadRequest // StatusInternalServerError
	}

	ctx := req.Context()

	fh, err := r.Read(ctx, rel_path)

	if err != nil {
		return nil, err, gohttp.StatusBadRequest // StatusInternalServerError
	}

	f, err := feature.LoadFeatureFromReader(fh)

	if err != nil && !warning.IsWarning(err) {
		return nil, err, gohttp.StatusInternalServerError
	}

	return f, nil, 0
}

func AltFeatureFromRequest(req *gohttp.Request, r reader.Reader) (geojson.Feature, error, int) {

	path := req.URL.Path

	wofid, err := IdFromPath(path)

	if err != nil {

		q := req.URL.Query()
		str_id := q.Get("id")

		if str_id == "" {
			return nil, err, gohttp.StatusNotFound
		}

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			return nil, err, gohttp.StatusBadRequest
		}

		wofid = id
	}

	q := req.URL.Query()

	alt_source := q.Get("source")
	alt_function := q.Get("function")

	args := &uri.URIArgs{
		Alternate: true,
		Source:    alt_source,
		Function:  alt_function,
	}

	rel_path, err := uri.Id2RelPath(wofid, args)

	if err != nil {
		return nil, err, gohttp.StatusBadRequest // StatusInternalServerError
	}

	ctx := req.Context()

	fh, err := r.Read(ctx, rel_path)

	if err != nil {
		return nil, err, gohttp.StatusBadRequest // StatusInternalServerError
	}

	f, err := feature.LoadFeatureFromReader(fh)

	if err != nil && !warning.IsWarning(err) {
		return nil, err, gohttp.StatusInternalServerError
	}

	return f, nil, 0
}
