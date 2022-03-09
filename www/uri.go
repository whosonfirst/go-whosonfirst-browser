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

	q := req.URL.Query()
	path := q.Get("id")

	if path == "" {
		path = req.URL.Path
	}

	return ParseURIFromPath(ctx, path, r)
}

func ParseURIFromPath(ctx context.Context, path string, r reader.Reader) (*URI, error, int) {

	wofid, uri_args, err := uri.ParseURI(path)

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
