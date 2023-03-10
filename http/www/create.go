package www

import (
	"html/template"
	"log"
	"net/http"

	aa_log "github.com/aaronland/go-log/v2"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	browser_properties "github.com/whosonfirst/go-whosonfirst-browser/v7/properties"
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
)

type CreateFeatureHandlerOptions struct {
	Reader           reader.Reader
	Authenticator    auth.Authenticator
	Logger           *log.Logger
	Template         *template.Template
	MapProvider      string
	URIs             *browser_uris.URIs
	Capabilities     *browser_capabilities.Capabilities
	CustomProperties []browser_properties.CustomProperty
}

type CreateFeatureVars struct {
	MapProvider      string
	Paths            *browser_uris.URIs
	Capabilities     *browser_capabilities.Capabilities
	CustomProperties []browser_properties.CustomProperty
}

func CreateFeatureHandler(opts *CreateFeatureHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		_, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			switch err.(type) {
			case auth.NotLoggedIn:

				signin_handler := opts.Authenticator.SigninHandler()
				signin_handler.ServeHTTP(rsp, req)
				return

			default:
				aa_log.Error(opts.Logger, "Failed to determine account for request, %v", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		vars := CreateFeatureVars{
			Paths:            opts.URIs,
			Capabilities:     opts.Capabilities,
			MapProvider:      opts.MapProvider,
			CustomProperties: opts.CustomProperties,
		}

		RenderTemplate(rsp, opts.Template, vars)
	}

	return http.HandlerFunc(fn), nil
}
