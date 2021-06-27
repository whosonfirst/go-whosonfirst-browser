package http

import (
	"encoding/json"
	"errors"
	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-spr-geojson"
	"html/template"
	_ "log"
	gohttp "net/http"
)

type SearchAPIHandlerOptions struct {
	Database        fulltext.FullTextDatabase
	EnableGeoJSON   bool
	GeoJSONReader   reader.Reader
	SPRPathResolver geojson.SPRPathResolver
}

type SearchHandlerOptions struct {
	Templates *template.Template
	Endpoints *Endpoints
	Database  fulltext.FullTextDatabase
}

type SearchVars struct {
	Endpoints *Endpoints
	Query     string
	Results   []spr.StandardPlacesResult
}

func SearchHandler(opts SearchHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("search")

	if t == nil {
		return nil, errors.New("Missing 'search' template.")
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		vars := SearchVars{
			Endpoints: opts.Endpoints,
		}

		term, err := sanitize.GetString(req, "term")

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		if term != "" {

			ctx := req.Context()

			results, err := opts.Database.QueryString(ctx, term)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			vars.Query = term
			vars.Results = results.Results()
		}

		rsp.Header().Set("Content-type", "text/html")

		RenderTemplate(rsp, t, vars)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}

func SearchAPIHandler(opts SearchAPIHandlerOptions) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		term, err := sanitize.GetString(req, "term")

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		if term == "" {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		format, err := sanitize.GetString(req, "format")

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		switch format {
		case "geojson":

			if !opts.EnableGeoJSON {
				gohttp.Error(rsp, "GeoJSON output is not enabled.", gohttp.StatusBadRequest)
				return
			}

		default:
			// pass
		}

		ctx := req.Context()

		results, err := opts.Database.QueryString(ctx, term)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
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
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

		default:

			enc := json.NewEncoder(rsp)
			err = enc.Encode(results)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
