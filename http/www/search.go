package www

import (
	"encoding/json"
	"errors"
	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spr-geojson"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"html/template"
	_ "log"
	"net/http"
)

type SearchAPIHandlerOptions struct {
	Database        fulltext.FullTextDatabase
	EnableGeoJSON   bool
	GeoJSONReader   reader.Reader
	SPRPathResolver geojson.SPRPathResolver
}

type SearchHandlerOptions struct {
	Templates    *template.Template
	Paths        *Paths
	Capabilities *Capabilities
	Database     fulltext.FullTextDatabase
	MapProvider  string
}

type SearchVars struct {
	Paths        *Paths
	Capabilities *Capabilities
	Query        string
	Results      []spr.StandardPlacesResult
}

func SearchHandler(opts SearchHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("search")

	if t == nil {
		return nil, errors.New("Missing 'search' template.")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		vars := SearchVars{
			Paths:        opts.Paths,
			Capabilities: opts.Capabilities,
		}

		term, err := sanitize.GetString(req, "term")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if term != "" {

			ctx := req.Context()

			results, err := opts.Database.QueryString(ctx, term)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			vars.Query = term
			vars.Results = results.Results()
		}

		rsp.Header().Set("Content-type", "text/html")

		RenderTemplate(rsp, t, vars)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
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
