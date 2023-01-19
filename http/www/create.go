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
	Endpoints     *Endpoints
}

type CreateFeatureVars struct {
	MapProvider string
	Endpoints   *Endpoints
	Placetypes  []*placetypes.WOFPlacetype
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
			http.Error(rsp, err.Error(), http.StatusUnauthorized)
			return
		}

		vars := CreateFeatureVars{
			Endpoints:   opts.Endpoints,
			MapProvider: opts.MapProvider,
			Placetypes:  all_placetypes,
		}

		RenderTemplate(rsp, opts.Template, vars)
	}

	return http.HandlerFunc(fn), nil
}
