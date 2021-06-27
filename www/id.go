package www

import (
	"errors"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"html/template"
	_ "log"
	"net/http"
	"path/filepath"
	"time"
)

type IDHandlerOptions struct {
	Templates *template.Template
	Endpoints *Endpoints
}

type IDVars struct {
	SPR          spr.StandardPlacesResult
	URI          string
	URIArgs      *uri.URIArgs
	IsAlternate  bool
	LastModified string
	Endpoints    *Endpoints
}

func IDHandler(r reader.Reader, opts IDHandlerOptions) (http.Handler, error) {

	id_t := opts.Templates.Lookup("id")

	if id_t == nil {
		return nil, errors.New("Missing id template")
	}

	alt_t := opts.Templates.Lookup("alt")

	if alt_t == nil {
		return nil, errors.New("Missing alt template")
	}

	error_t := opts.Templates.Lookup("error")

	if error_t == nil {
		return nil, errors.New("Missing error template")
	}

	notfound_t := opts.Templates.Lookup("id")

	if notfound_t == nil {
		return nil, errors.New("Missing notfound template")
	}

	handle_other := func(rsp http.ResponseWriter, req *http.Request, f geojson.Feature, endpoint string) {

		if endpoint == "" {

			vars := NotFoundVars{
				Endpoints: opts.Endpoints,
			}

			RenderTemplate(rsp, notfound_t, vars)
			return
		}

		url := filepath.Join(endpoint, f.Id())
		http.Redirect(rsp, req, url, http.StatusMovedPermanently)
		return
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, _ := ParseURIFromRequest(req, r)

		if err != nil {

			vars := ErrorVars{
				Error:     err,
				Endpoints: opts.Endpoints,
				// status...
			}

			RenderTemplate(rsp, error_t, vars)
			return
		}

		f := uri.Feature

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

		/*
			data_url := new(url.URL)
			data_url.Scheme = req.URL.Scheme
			data_url.Host = req.URL.Host
			data_url.Path = opts.Endpoints.Data

			data_endpoint := data_url.String()

			png_url := new(url.URL)
			png_url.Scheme = req.URL.Scheme
			png_url.Host = req.URL.Host
			png_url.Path = opts.Endpoints.Data

			png_endpoint := png_url.String()
		*/

		vars := IDVars{
			SPR:          s,
			URI:          uri.URI,
			URIArgs:      uri.URIArgs,
			IsAlternate:  uri.IsAlternate,
			LastModified: lastmod,
			Endpoints:    opts.Endpoints,
		}

		t := id_t

		if uri.IsAlternate {
			t = alt_t
		}

		RenderTemplate(rsp, t, vars)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
