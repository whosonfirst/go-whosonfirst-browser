package api

import (
	"bytes"
	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v6/http"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-writer/v3"
	"io"
	"log"
	"net/http"
)

type UpdateGeometryHandlerOptions struct {
	Reader        reader.Reader
	Writer        writer.Writer
	Exporter      export.Exporter
	Authenticator auth.Authenticator
	Logger        *log.Logger
}

func UpdateGeometryHandler(opts *UpdateGeometryHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

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

		// TO DO: PIP STUFF HERE (here?)

		has_changes, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if has_changes {

			exp_body, err := opts.Exporter.Export(ctx, new_body)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			br := bytes.NewReader(exp_body)

			_, err = opts.Writer.Write(ctx, uri.URI, br)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
			}
		}

		// TBD: return updated body here?
		return
	}

	return http.HandlerFunc(fn), nil
}
