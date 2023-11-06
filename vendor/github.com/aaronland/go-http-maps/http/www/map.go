package www

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/aaronland/go-http-maps/provider"
)

type MapHandlerOptions struct {
	Templates        *template.Template
	InitialLatitude  float64
	InitialLongitude float64
	InitialZoom      int
	MapProvider      provider.Provider
}

type MapHandlerVars struct {
	InitialLatitude  float64
	InitialLongitude float64
	InitialZoom      int
	MapProvider      string
}

func MapHandler(opts *MapHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("map")

	if t == nil {
		return nil, errors.New("Missing 'map' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		vars := MapHandlerVars{
			MapProvider:      opts.MapProvider.Scheme(),
			InitialLatitude:  opts.InitialLatitude,
			InitialLongitude: opts.InitialLongitude,
			InitialZoom:      opts.InitialZoom,
		}

		rsp.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := t.Execute(rsp, vars)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
