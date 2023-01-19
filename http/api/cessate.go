package api

import (
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v6/http"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"log"
	"net/http"
)

type CessateFeatureHandlerOptions struct {
	Reader        reader.Reader
	WriterURIs    []string
	Exporter      export.Exporter
	Authenticator auth.Authenticator
	Logger        *log.Logger
	Cache         cache.Cache
}

func CessateFeatureHandler(opts *CessateFeatureHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		switch req.Method {
		case "POST":
			// pass
		default:
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := req.Context()

		_, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		uri, err, _ := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		body := uri.Feature

		new_body, err := export.CessateRecord(ctx, opts.Exporter, body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		publish_opts := &publishFeatureOptions{
			Logger:     opts.Logger,
			WriterURIs: opts.WriterURIs,
			Exporter:   opts.Exporter,
			Cache:      opts.Cache,
			URI:        uri,
		}

		final, err := publishFeature(ctx, publish_opts, new_body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Write(final)
		return
	}

	return http.HandlerFunc(fn), nil
}
