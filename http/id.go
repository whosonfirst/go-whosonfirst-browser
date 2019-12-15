package http

import (
	"errors"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"html/template"
	"log"
	gohttp "net/http"
	"net/url"
	"path/filepath"
	"time"
)

type IDHandlerOptions struct {
	DataEndpoint string
	PngEndpoint string
	Templates    *template.Template
}

type IDVars struct {
	SPR          spr.StandardPlacesResult
	LastModified string
	DataEndpoint string
	PngEndpoint string
}

func IDHandler(r reader.Reader, opts IDHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("id")

	if t == nil {
		return nil, errors.New("Missing id template")
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		path := req.URL.Path
		ext := filepath.Ext(path)

		// TODO: check for .geojson, .svg, etc and redirect accordingly
		
		log.Println(path, ext)

		f, err, status := FeatureFromRequest(req, r)

		if err != nil {
			log.Println("SAD", err)
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		s, err := f.SPR()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
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
			PngEndpoint: png_endpoint,
		}

		err = t.Execute(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
