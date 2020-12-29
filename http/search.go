package http

import (
	"errors"
	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"html/template"
	_ "log"
	gohttp "net/http"
)

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
