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

type DeprecateFeatureHandlerOptions struct {
	Reader        reader.Reader
	WriterURIs    []string
	Exporter      export.Exporter
	Authenticator auth.Authenticator
	Logger        *log.Logger
	Cache         cache.Cache
}

func DeprecateFeatureHandler(opts *DeprecateFeatureHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		switch req.Method {
		case "POST":
			// pass
		default:
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		acct, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			switch err.(type) {
			case auth.NotLoggedIn:
				opts.Logger.Printf("Failed to determine account for request, %v", err)
				http.Error(rsp, "Not authorized", http.StatusUnauthorized)
				return
			default:
				opts.Logger.Printf("Failed to determine account for request, %v", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		uri, err, _ := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		body := uri.Feature

		new_body, err := export.DeprecateRecord(ctx, opts.Exporter, body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}

		publish_opts := &publishFeatureOptions{
			Logger:     opts.Logger,
			WriterURIs: opts.WriterURIs,
			Exporter:   opts.Exporter,
			URI:        uri,
			Cache:      opts.Cache,
			Account:    acct,
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
