package www

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
)

type CreateFeatureHandlerOptions struct {
	Reader        reader.Reader
	Authenticator auth.Authenticator
	Logger        *log.Logger
	Template      *template.Template
	MapProvider   string
	Endpoints     *Endpoints
}

type CreateFeatureVars struct {
	MapProvider string
	Endpoints   *Endpoints
	// To do: Support alternate geometries
}

func CreateFeatureHandler(opts *CreateFeatureHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		_, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusUnauthorized)
			return
		}

		vars := CreateFeatureVars{
			Endpoints:   opts.Endpoints,
			MapProvider: opts.MapProvider,
		}

		RenderTemplate(rsp, opts.Template, vars)
	}

	return http.HandlerFunc(fn), nil
}
