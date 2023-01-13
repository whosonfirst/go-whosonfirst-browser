package www

import (
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v6/http"
	"html/template"
	"log"
	"net/http"
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
	Geometry    string
	Id          int64
}

func EditGeometryHandler(opts *EditGeometryHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		// ctx := req.Context()

		_, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusUnauthorized)
			return
		}

		uri, err, _ := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		vars := EditGeometryVars{
			MapProvider: opts.MapProvider,
			Id:          uri.Id,
		}

		RenderTemplate(rsp, opts.Template, vars)
	}

	return http.HandlerFunc(fn), nil
}
