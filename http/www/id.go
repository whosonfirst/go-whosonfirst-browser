package www

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/whosonfirst/go-reader"
	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	browser_http "github.com/whosonfirst/go-whosonfirst-browser/v7/http"
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type IDHandlerOptions struct {
	Templates    *template.Template
	Reader       reader.Reader
	Logger       *log.Logger
	MapProvider  string
	URIs         *browser_uris.URIs
	Capabilities *browser_capabilities.Capabilities
}

type IDVars struct {
	SPR          spr.StandardPlacesResult
	URI          string
	URIArgs      *uri.URIArgs
	IsAlternate  bool
	LastModified string
	Paths        *browser_uris.URIs
	Capabilities *browser_capabilities.Capabilities
	MapProvider  string
}

func IDHandler(opts IDHandlerOptions) (http.Handler, error) {

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

	handle_other := func(rsp http.ResponseWriter, req *http.Request, f []byte, endpoint string) {

		if endpoint == "" {

			vars := NotFoundVars{
				URIs: opts.URIs,
			}

			RenderTemplate(rsp, notfound_t, vars)
			return
		}

		id, err := properties.Id(f)

		if err != nil {

			vars := ErrorVars{
				Error: err,
				URIs:  opts.URIs,
				// status...
			}

			RenderTemplate(rsp, error_t, vars)
		}

		str_id := strconv.FormatInt(id, 10)
		url := filepath.Join(endpoint, str_id)

		http.Redirect(rsp, req, url, http.StatusMovedPermanently)
		return
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, _ := browser_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {

			opts.Logger.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			vars := ErrorVars{
				Error: err,
				URIs:  opts.URIs,
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
			handle_other(rsp, req, f, opts.URIs.GeoJSON)
			return
		case ".png":
			handle_other(rsp, req, f, opts.URIs.PNG)
			return
		case ".spr":
			handle_other(rsp, req, f, opts.URIs.SPR)
			return
		case ".svg":
			handle_other(rsp, req, f, opts.URIs.SVG)
			return
		default:
			// pass
		}

		s, err := spr.WhosOnFirstSPR(f)

		if err != nil {

			vars := ErrorVars{
				Error: err,
				URIs:  opts.URIs,
			}

			RenderTemplate(rsp, error_t, vars)
			return
		}

		now := time.Now()
		lastmod := now.Format(time.RFC3339)

		vars := IDVars{
			SPR:          s,
			URI:          uri.URI,
			URIArgs:      uri.URIArgs,
			IsAlternate:  uri.IsAlternate,
			LastModified: lastmod,
			Paths:        opts.URIs,
			Capabilities: opts.Capabilities,
			MapProvider:  opts.MapProvider,
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
