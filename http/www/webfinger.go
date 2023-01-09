package www

import (
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v5/http"
	"github.com/whosonfirst/go-whosonfirst-browser/v5/webfinger"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type WebfingerHandlerOptions struct {
	Reader       reader.Reader
	Logger       *log.Logger
	Paths        *Paths
	Capabilities *Capabilities
}

func WebfingerHandler(opts *WebfingerHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		wf_scheme := "https"

		if req.TLS == nil {
			wf_scheme = "http"
		}

		wf_host := req.Host

		wof_uri, err, status := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {

			opts.Logger.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			http.Error(rsp, err.Error(), status)
			return
		}

		pt, err := properties.Placetype(wof_uri.Feature)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		name, err := properties.Name(wof_uri.Feature)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		lastmod := properties.LastModified(wof_uri.Feature)
		str_lastmod := strconv.FormatInt(lastmod, 10)

		rel_path, err := uri.Id2RelPath(wof_uri.Id, wof_uri.URIArgs)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		path_geojson, err := url.JoinPath(opts.Paths.GeoJSON, rel_path)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		geojson_uri := url.URL{}
		geojson_uri.Scheme = wf_scheme
		geojson_uri.Host = wf_host
		geojson_uri.Path = path_geojson

		subject := fmt.Sprintf("acct:%d@%s", wof_uri.Id, wf_host)

		props := map[string]string{
			"http://whosonfirst.org/properties/wof/placetype":    pt,
			"http://whosonfirst.org/properties/wof/name":         name,
			"http://whosonfirst.org/properties/wof/lastmodified": str_lastmod,
		}

		links := make([]webfinger.Link, 0)

		aliases := []string{
			geojson_uri.String(),
		}

		if opts.Capabilities.GeoJSON {

			l := webfinger.Link{
				HRef: geojson_uri.String(),
				Type: "application/geo+json",
				Rel:  "x-whosonfirst-rel#geojson",
			}

			links = append(links, l)
		}

		if opts.Capabilities.GeoJSONLD {

			path_geojsonld, err := url.JoinPath(opts.Paths.GeoJSONLD, rel_path)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			uri := url.URL{}
			uri.Scheme = wf_scheme
			uri.Host = wf_host
			uri.Path = path_geojsonld

			l := webfinger.Link{
				HRef: uri.String(),
				Type: "application/geo+json",
				Rel:  "x-whosonfirst-rel#geojson-ld",
			}

			links = append(links, l)

		}

		if opts.Capabilities.SVG {

			path_svg, err := url.JoinPath(opts.Paths.SVG, rel_path)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			uri := url.URL{}
			uri.Scheme = wf_scheme
			uri.Host = wf_host
			uri.Path = path_svg

			l := webfinger.Link{
				HRef: uri.String(),
				Type: "image/svg+xml",
				Rel:  "x-whosonfirst-rel#svg",
			}

			links = append(links, l)
		}

		if opts.Capabilities.PNG {

			path_png, err := url.JoinPath(opts.Paths.PNG, rel_path)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			uri := url.URL{}
			uri.Scheme = wf_scheme
			uri.Host = wf_host
			uri.Path = path_png

			l := webfinger.Link{
				HRef: uri.String(),
				Type: "image/png",
				Rel:  "x-whosonfirst-rel#png",
			}

			links = append(links, l)
		}

		if opts.Capabilities.Select {

			path_select, err := url.JoinPath(opts.Paths.Select, rel_path)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			uri := url.URL{}
			uri.Scheme = wf_scheme
			uri.Host = wf_host
			uri.Path = path_select

			l := webfinger.Link{
				HRef: uri.String(),
				Type: "application/json",
				Rel:  "x-whosonfirst-rel#select",
			}

			links = append(links, l)
		}

		if opts.Capabilities.NavPlace {

			path_navplace, err := url.JoinPath(opts.Paths.NavPlace, rel_path)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			uri := url.URL{}
			uri.Scheme = wf_scheme
			uri.Host = wf_host
			uri.Path = path_navplace

			l := webfinger.Link{
				HRef: uri.String(),
				Type: "application/geo+json",
				Rel:  "x-whosonfirst-rel#navplace",
			}

			links = append(links, l)
		}

		if opts.Capabilities.SPR {

			path_spr, err := url.JoinPath(opts.Paths.SPR, rel_path)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			uri := url.URL{}
			uri.Scheme = wf_scheme
			uri.Host = wf_host
			uri.Path = path_spr

			l := webfinger.Link{
				HRef: uri.String(),
				Type: "application/json",
				Rel:  "x-whosonfirst-rel#spr",
			}

			links = append(links, l)
		}

		r := webfinger.Resource{
			Subject:    subject,
			Properties: props,
			Aliases:    aliases,
			Links:      links,
		}

		rsp.Header().Set("Content-type", webfinger.ContentType)

		enc := json.NewEncoder(rsp)
		err = enc.Encode(&r)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
