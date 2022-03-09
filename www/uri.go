package www

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/warning"
	_ "log"
	"net/http"
	"path/filepath"
	_ "strconv"
	"strings"
)

type URI struct {
	Id          int64
	URI         string
	Feature     geojson.Feature
	URIArgs     *uri.URIArgs
	IsAlternate bool
}

func ParseURIFromRequest(req *http.Request, r reader.Reader) (*URI, error, int) {

	ctx := req.Context()	
	path := req.URL.Path

	return ParseURIFromPath(ctx, path, r)
}

func ParseURIFromPath(ctx context.Context, path string, r reader.Reader) (*URI, error, int) {
		
	wofid, uri_args, err := uri.ParseURI(path)

	/*
	if err != nil || wofid == -1 {

		q := req.URL.Query()
		str_id := q.Get("id")

		if str_id == "" {
			return nil, fmt.Errorf("Failed to parse %s and ?id parameter is empty, %w", err), http.StatusNotFound
		}

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse %s and ?id=%s is invalid, %w", err), http.StatusBadRequest
		}

		wofid = id

		uri_args = &uri.URIArgs{
			IsAlternate: false,
		}
	}
	*/
	
	rel_path, err := uri.Id2RelPath(wofid, uri_args)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive relative path from %d (%s), %w", wofid, path, err), http.StatusBadRequest // StatusInternalServerError
	}

	fh, err := r.Read(ctx, rel_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to read %s, %w", rel_path, err), http.StatusBadRequest // StatusInternalServerError
	}

	f, err := feature.LoadFeatureFromReader(fh)

	if err != nil && !warning.IsWarning(err) {
		return nil, fmt.Errorf("Failed to read feature for %s, %w", rel_path, err), http.StatusInternalServerError
	}

	fname, err := uri.Id2Fname(wofid, uri_args)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive filename from %d (%s), %w", wofid, path, err), http.StatusInternalServerError
	}

	ext := filepath.Ext(fname)
	fname = strings.Replace(fname, ext, "", 1)

	uri := &URI{
		Id:          wofid,
		URI:         fname,
		URIArgs:     uri_args,
		Feature:     f,
		IsAlternate: uri_args.IsAlternate,
	}

	return uri, nil, 0
}
