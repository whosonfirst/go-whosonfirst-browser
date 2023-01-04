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

		aliases := []string{
			geojson_uri.String(),
		}

		r := webfinger.Resource{
			Subject:    subject,
			Properties: props,
			Aliases:    aliases,
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
