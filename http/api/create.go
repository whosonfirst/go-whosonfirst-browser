package api

import (
	"io"
	"log"
	"net/http"

	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v6/http"
	"github.com/whosonfirst/go-whosonfirst-browser/v6/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
)

type CreateFeatureHandlerOptions struct {
	Reader                reader.Reader
	Cache                 cache.Cache
	WriterURIs            []string
	Exporter              export.Exporter
	Authenticator         auth.Authenticator
	Logger                *log.Logger
	PointInPolygonService *pointinpolygon.PointInPolygonService
}

func CreateFeatureHandler(opts *CreateFeatureHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		switch req.Method {
		case "PUT":
			// pass
		default:
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := req.Context()

		_, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusUnauthorized)
			return
		}

		uri, err, _ := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = geojson.UnmarshalFeature(body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		// Sanity check / validate feature here

		_, new_body, err := opts.PointInPolygonService.Update(ctx, body)

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
