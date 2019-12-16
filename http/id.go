package http

import (
	"errors"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"html/template"
	_ "log"
	gohttp "net/http"
	"net/url"
	"path/filepath"
	"time"
)

type IDHandlerOptions struct {
	Templates    *template.Template
	DataEndpoint string
	SprEndpoint  string
	SvgEndpoint  string
	PngEndpoint  string
}

type IDVars struct {
	SPR          spr.StandardPlacesResult
	LastModified string
	DataEndpoint string
	PngEndpoint  string
}

func IDHandler(r reader.Reader, opts IDHandlerOptions) (gohttp.Handler, error) {

	id_t := opts.Templates.Lookup("id")

	if id_t == nil {
		return nil, errors.New("Missing id template")
	}

	error_t := opts.Templates.Lookup("error")

	if error_t == nil {
		return nil, errors.New("Missing error template")
	}

	notfound_t := opts.Templates.Lookup("id")

	if notfound_t == nil {
		return nil, errors.New("Missing notfound template")
	}
		
	handle_other := func(rsp gohttp.ResponseWriter, req *gohttp.Request, f geojson.Feature, endpoint string) {

		if endpoint == "" {

			RenderTemplate(rsp, notfound_t, nil)
			return
			
			// gohttp.Error(rsp, "Not found", gohttp.StatusNotFound)
			// return
		}

		url := filepath.Join(endpoint, f.Id())
		gohttp.Redirect(rsp, req, url, gohttp.StatusMovedPermanently)
		return
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		f, err, _ := FeatureFromRequest(req, r)

		if err != nil {

			RenderTemplate(rsp, error_t, err)
			return
			
			// gohttp.Error(rsp, err.Error(), status)
			// return
		}

		path := req.URL.Path
		ext := filepath.Ext(path)

		switch ext {
		case ".geojson":
			handle_other(rsp, req, f, opts.DataEndpoint)
			return
		case ".png":
			handle_other(rsp, req, f, opts.PngEndpoint)
			return
		case ".spr":
			handle_other(rsp, req, f, opts.SprEndpoint)
			return
		case ".svg":
			handle_other(rsp, req, f, opts.SvgEndpoint)
			return
		default:
			// pass
		}

		s, err := f.SPR()

		if err != nil {
			RenderTemplate(rsp, error_t, err)
			return
		}

		now := time.Now()
		lastmod := now.Format(time.RFC3339)

		data_url := new(url.URL)
		data_url.Scheme = req.URL.Scheme
		data_url.Host = req.URL.Host
		data_url.Path = opts.DataEndpoint

		data_endpoint := data_url.String()

		png_url := new(url.URL)
		png_url.Scheme = req.URL.Scheme
		png_url.Host = req.URL.Host
		png_url.Path = opts.DataEndpoint

		png_endpoint := png_url.String()

		vars := IDVars{
			SPR:          s,
			LastModified: lastmod,
			DataEndpoint: data_endpoint,
			PngEndpoint:  png_endpoint,
		}

		RenderTemplate(rsp, id_t, vars)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
