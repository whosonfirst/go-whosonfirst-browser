package http

import (
	"errors"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"html/template"
	_ "log"
	gohttp "net/http"
	_ "net/url"
	"path/filepath"
	"time"
)

type AltHandlerOptions struct {
	Templates *template.Template
	Endpoints *Endpoints
}

type AltVars struct {
	SPR          spr.StandardPlacesResult
	LastModified string
	Endpoints    *Endpoints
}

func AltHandler(r reader.Reader, opts AltHandlerOptions) (gohttp.Handler, error) {

	alt_t := opts.Templates.Lookup("alt")

	if alt_t == nil {
		return nil, errors.New("Missing alt template")
	}

	error_t := opts.Templates.Lookup("error")

	if error_t == nil {
		return nil, errors.New("Missing error template")
	}

	notfound_t := opts.Templates.Lookup("alt")

	if notfound_t == nil {
		return nil, errors.New("Missing notfound template")
	}

	handle_other := func(rsp gohttp.ResponseWriter, req *gohttp.Request, f geojson.Feature, endpoint string) {

		if endpoint == "" {

			vars := NotFoundVars{
				Endpoints: opts.Endpoints,
			}

			RenderTemplate(rsp, notfound_t, vars)
			return
		}

		url := filepath.Join(endpoint, f.Id())
		gohttp.Redirect(rsp, req, url, gohttp.StatusMovedPermanently)
		return
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		f, err, _ := AltFeatureFromRequest(req, r)

		if err != nil {

			vars := ErrorVars{
				Error:     err,
				Endpoints: opts.Endpoints,
				// status...
			}

			RenderTemplate(rsp, error_t, vars)
			return
		}

		path := req.URL.Path
		ext := filepath.Ext(path)

		switch ext {
		case ".geojson":
			handle_other(rsp, req, f, opts.Endpoints.Data)
			return
		case ".png":
			handle_other(rsp, req, f, opts.Endpoints.Png)
			return
		case ".spr":
			handle_other(rsp, req, f, opts.Endpoints.Spr)
			return
		case ".svg":
			handle_other(rsp, req, f, opts.Endpoints.Svg)
			return
		default:
			// pass
		}

		s, err := f.SPR()

		if err != nil {

			vars := ErrorVars{
				Error:     err,
				Endpoints: opts.Endpoints,
			}

			RenderTemplate(rsp, error_t, vars)
			return
		}

		now := time.Now()
		lastmod := now.Format(time.RFC3339)

		vars := AltVars{
			SPR:          s,
			LastModified: lastmod,
			Endpoints:    opts.Endpoints,
		}

		RenderTemplate(rsp, alt_t, vars)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
