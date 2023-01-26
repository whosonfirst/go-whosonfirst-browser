package api

import (
	"encoding/json"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spr-geojson"
)

type SearchAPIHandlerOptions struct {
	Database        fulltext.FullTextDatabase
	EnableGeoJSON   bool
	GeoJSONReader   reader.Reader
	SPRPathResolver geojson.SPRPathResolver
}

func SearchAPIHandler(opts SearchAPIHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		term, err := sanitize.GetString(req, "term")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if term == "" {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		format, err := sanitize.GetString(req, "format")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		switch format {
		case "geojson":

			if !opts.EnableGeoJSON {
				http.Error(rsp, "GeoJSON output is not enabled.", http.StatusBadRequest)
				return
			}

		default:
			// pass
		}

		ctx := req.Context()

		results, err := opts.Database.QueryString(ctx, term)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/json")

		switch format {
		case "geojson":

			as_opts := &geojson.AsFeatureCollectionOptions{
				Reader:          opts.GeoJSONReader,
				Writer:          rsp,
				SPRPathResolver: opts.SPRPathResolver,
			}

			err := geojson.AsFeatureCollection(ctx, results, as_opts)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

		default:

			enc := json.NewEncoder(rsp)
			err = enc.Encode(results)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
