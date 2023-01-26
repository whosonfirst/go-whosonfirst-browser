package www

import (
	"errors"
	"html/template"
	_ "log"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type SearchHandlerOptions struct {
	Templates    *template.Template
	Database     fulltext.FullTextDatabase
	MapProvider  string
	URIs         *browser_uris.URIs
	Capabilities *browser_capabilities.Capabilities
}

type SearchVars struct {
	Paths        *browser_uris.URIs
	Capabilities *browser_capabilities.Capabilities
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
			Paths:        opts.URIs,
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
