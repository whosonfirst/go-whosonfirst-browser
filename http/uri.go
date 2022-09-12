package http

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	go_http "net/http"
	"path/filepath"
	"strings"
)

type URI struct {
	Id          int64
	URI         string
	Feature     []byte
	URIArgs     *wof_uri.URIArgs
	IsAlternate bool
}

func ParseURIFromRequest(req *go_http.Request, r reader.Reader) (*URI, error, int) {

	ctx := req.Context()

	q := req.URL.Query()
	path := q.Get("id")

	if path == "" {
		path = req.URL.Path
	}

	return ParseURIFromPath(ctx, path, r)
}

func ParseURIFromPath(ctx context.Context, path string, r reader.Reader) (*URI, error, int) {

	wofid, uri_args, err := wof_uri.ParseURI(path)

	rel_path, err := wof_uri.Id2RelPath(wofid, uri_args)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive relative path from %d (%s), %w", wofid, path, err), go_http.StatusBadRequest // StatusInternalServerError
	}

	fh, err := r.Read(ctx, rel_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to read %s, %w", rel_path, err), go_http.StatusBadRequest // StatusInternalServerError
	}

	f, err := io.ReadAll(fh)

	if err != nil {
		return nil, fmt.Errorf("Failed to read feature for %s, %w", rel_path, err), go_http.StatusInternalServerError
	}

	fname, err := wof_uri.Id2Fname(wofid, uri_args)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive filename from %d (%s), %w", wofid, path, err), go_http.StatusInternalServerError
	}

	ext := filepath.Ext(fname)
	fname = strings.Replace(fname, ext, "", 1)

	parsed_uri := &URI{
		Id:          wofid,
		URI:         fname,
		URIArgs:     uri_args,
		Feature:     f,
		IsAlternate: uri_args.IsAlternate,
	}

	return parsed_uri, nil, 0
}
