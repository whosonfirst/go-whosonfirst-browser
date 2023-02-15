package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	aa_log "github.com/aaronland/go-log"
	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v7/http"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
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
				aa_log.Error(opts.Logger, "Failed to determine account for request, %v", err)
				http.Error(rsp, "Not authorized", http.StatusUnauthorized)
				return
			default:
				aa_log.Error(opts.Logger, "Failed to determine account for request, %v", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		uri, err, _ := wof_http.ParseURIFromRequest(req, opts.Reader)

		if err != nil {
			aa_log.Error(opts.Logger, "Failed to parse URI from request %s, %v", req.URL.Path, err)
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		update, err := io.ReadAll(req.Body)

		if err != nil {
			aa_log.Error(opts.Logger, "Failed to read body from %s, %v", req.URL.Path, err)
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := geojson.UnmarshalFeature(update)

		if err != nil {
			aa_log.Error(opts.Logger, "Failed to unmarshal GeoJSON from %s, %v", req.URL.Path, err)
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
			aa_log.Error(opts.Logger, "Failed to assign properties to %s, %v", req.URL.Path, err)
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if !has_changes {
			return
		}

		has_changes, new_body, err = opts.PointInPolygonService.Update(ctx, new_body)

		if err != nil {
			aa_log.Error(opts.Logger, "Failed to assign point-in-polygon data to %s, %v", req.URL.Path, err)
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if !has_changes {
			return
		}

		name, err := properties.Name(new_body)

		if err != nil {
			aa_log.Error(opts.Logger, "Failed to derive name %s, %v", req.URL.Path, err)
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		// To do (maybe): Make these customizable with URI templates in opts
		// https://pkg.go.dev/github.com/jtacoma/uritemplates

		title := fmt.Sprintf("Update geometry for '%s' (%d)", name, uri.Id)
		description := title

		now := time.Now()
		ts := now.Unix()

		branch := fmt.Sprintf("update-geometry-%d-%d", ts, uri.Id)

		publish_opts := &publishFeatureOptions{
			Logger:      opts.Logger,
			WriterURIs:  opts.WriterURIs,
			Exporter:    opts.Exporter,
			Cache:       opts.Cache,
			URI:         uri,
			Account:     acct,
			Title:       title,
			Description: description,
			Branch:      branch,
		}

		final, err := publishFeature(ctx, publish_opts, new_body)

		if err != nil {
			aa_log.Error(opts.Logger, "Failed to publish feature for %s, %v", req.URL.Path, err)
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Write(final)

		aa_log.Debug(opts.Logger, "Wrote updated geometry for %d to %s branch", uri.Id, branch)
		return
	}

	return http.HandlerFunc(fn), nil
}
