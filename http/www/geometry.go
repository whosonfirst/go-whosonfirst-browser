package www

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	browser_http "github.com/whosonfirst/go-whosonfirst-browser/v7/http"
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type EditGeometryHandlerOptions struct {
	Reader        reader.Reader
	Authenticator auth.Authenticator
	Logger        *log.Logger
	Template      *template.Template
	MapProvider   string
	URIs          *browser_uris.URIs
	Capabilities  *browser_capabilities.Capabilities
}

type EditGeometryVars struct {
	MapProvider  string
	Id           int64
	Paths        *browser_uris.URIs
	Capabilities *browser_capabilities.Capabilities
	// To do: Support alternate geometries
}

func EditGeometryHandler(opts *EditGeometryHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		_, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			switch err.(type) {
			case auth.NotLoggedIn:

				signin_handler := opts.Authenticator.SigninHandler()
				signin_handler.ServeHTTP(rsp, req)
				return

			default:
				opts.Logger.Printf("Failed to determine account for request, %v", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		path, err, status := browser_http.DerivePathFromRequest(req)

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
			Paths:        opts.URIs,
			Capabilities: opts.Capabilities,
			MapProvider:  opts.MapProvider,
			Id:           id,
		}

		RenderTemplate(rsp, opts.Template, vars)
	}

	return http.HandlerFunc(fn), nil
}
