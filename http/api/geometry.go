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

type UpdateGeometryHandlerOptions struct {
	Reader                reader.Reader
	Cache                 cache.Cache
	WriterURIs            []string
	Exporter              export.Exporter
	Authenticator         auth.Authenticator
	Logger                *log.Logger
	PointInPolygonService *pointinpolygon.PointInPolygonService
}

func UpdateGeometryHandler(opts *UpdateGeometryHandlerOptions) (http.Handler, error) {

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

		update, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := geojson.UnmarshalFeature(update)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		geojson_geometry := geojson.NewGeometry(f.Geometry)

		updates := map[string]interface{}{
			"geometry": geojson_geometry,
		}

		body := uri.Feature

		has_changes, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if !has_changes {
			return
		}

		has_changes, new_body, err = opts.PointInPolygonService.Update(ctx, new_body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if !has_changes {
			return
		}

		publish_opts := &publishFeatureOptions{
			Logger:     opts.Logger,
			WriterURIs: opts.WriterURIs,
			Exporter:   opts.Exporter,
			Cache:      opts.Cache,
			URI:        uri,
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
