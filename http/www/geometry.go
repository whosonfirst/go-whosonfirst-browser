package www

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v6/http"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type EditGeometryHandlerOptions struct {
	Reader        reader.Reader
	Authenticator auth.Authenticator
	Logger        *log.Logger
	Template      *template.Template
	MapProvider   string
}

type EditGeometryVars struct {
	MapProvider string
	Id          int64
	// To do: Support alternate geometries
}

func EditGeometryHandler(opts *EditGeometryHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		_, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusUnauthorized)
			return
		}

		path, err, status := wof_http.DerivePathFromRequest(req)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		id, _, err := uri.ParseURI(path)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		vars := EditGeometryVars{
			MapProvider: opts.MapProvider,
			Id:          id,
		}

		RenderTemplate(rsp, opts.Template, vars)
	}

	return http.HandlerFunc(fn), nil
}
