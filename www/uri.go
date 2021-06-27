package www

import (
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/warning"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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

	path := req.URL.Path

	wofid, uri_args, err := uri.ParseURI(path)

	log.Println("PARSE", path, wofid, uri_args, err)

	if err != nil || wofid == -1 {

		q := req.URL.Query()
		str_id := q.Get("id")

		if str_id == "" {
			return nil, err, http.StatusNotFound
		}

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			return nil, err, http.StatusBadRequest
		}

		wofid = id

		uri_args = &uri.URIArgs{
			IsAlternate: false,
		}
	}

	rel_path, err := uri.Id2RelPath(wofid, uri_args)

	log.Println("REL PATH", rel_path)
		
	if err != nil {
		return nil, err, http.StatusBadRequest // StatusInternalServerError
	}

	ctx := req.Context()

	fh, err := r.Read(ctx, rel_path)

	log.Printf("READ FROM %T, %v\n", r, err)
	
	if err != nil {
		return nil, err, http.StatusBadRequest // StatusInternalServerError
	}

	f, err := feature.LoadFeatureFromReader(fh)

	log.Println("F", err)
	
	if err != nil && !warning.IsWarning(err) {
		return nil, err, http.StatusInternalServerError
	}

	fname, err := uri.Id2Fname(wofid, uri_args)

	if err != nil {
		return nil, err, http.StatusInternalServerError
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
