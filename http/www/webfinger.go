package www

import (
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v5/http"
	"github.com/whosonfirst/go-whosonfirst-browser/v5/webfinger"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"log"
	"net/http"
	"strconv"
)

type WebfingerHandlerOptions struct {
	Reader reader.Reader
	Logger *log.Logger
}

func WebfingerHandler(opts *WebfingerHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		// Not sure this will work because of query parameters

		uri, err, status := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {

			opts.Logger.Printf("Failed to parse URI from request %s, %v", req.URL, err)

			http.Error(rsp, err.Error(), status)
			return
		}

		pt, err := properties.Placetype(uri.Feature)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}

		name, err := properties.Name(uri.Feature)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}

		lastmod := properties.LastModified(uri.Feature)
		str_lastmod := strconv.FormatInt(lastmod, 10)

		subject := fmt.Sprintf("acct:%d@whosonfirst.org", uri.Id)

		props := map[string]string{
			"http://whosonfirst.org/properties/wof/placetype":    pt,
			"http://whosonfirst.org/properties/wof/name":         name,
			"http://whosonfirst.org/properties/wof/lastmodified": str_lastmod,
		}

		aliases := []string{
			uri.URI,
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
