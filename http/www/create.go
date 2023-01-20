package www

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

type CreateFeatureHandlerOptions struct {
	Reader        reader.Reader
	Authenticator auth.Authenticator
	Logger        *log.Logger
	Template      *template.Template
	MapProvider   string
	Paths         *Paths
	Capabilities  *Capabilities
}

type CreateFeatureVars struct {
	MapProvider  string
	Placetypes   []*placetypes.WOFPlacetype
	Paths        *Paths
	Capabilities *Capabilities
	// To do: Support alternate geometries
}

func CreateFeatureHandler(opts *CreateFeatureHandlerOptions) (http.Handler, error) {

	pt, err := placetypes.GetPlacetypeByName("planet")

	if err != nil {
		return nil, fmt.Errorf("Failed to load placetype for planet, %w", err)
	}

	roles := []string{
		"common",
		"common_optional",
		"optional",
	}

	all_placetypes := placetypes.DescendantsForRoles(pt, roles)

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

		vars := CreateFeatureVars{
			Paths:        opts.Paths,
			Capabilities: opts.Capabilities,
			MapProvider:  opts.MapProvider,
			Placetypes:   all_placetypes,
		}

		RenderTemplate(rsp, opts.Template, vars)
	}

	return http.HandlerFunc(fn), nil
}
