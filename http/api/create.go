package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	browser_properties "github.com/whosonfirst/go-whosonfirst-browser/v7/properties"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-validate"
)

type CreateFeatureHandlerOptions struct {
	Reader                reader.Reader
	Cache                 cache.Cache
	WriterURIs            []string
	Exporter              export.Exporter
	Authenticator         auth.Authenticator
	Logger                *log.Logger
	PointInPolygonService *pointinpolygon.PointInPolygonService
	CustomProperties      []browser_properties.CustomProperty
}

func CreateFeatureHandler(opts *CreateFeatureHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		switch req.Method {
		case "PUT":
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

		// START OF validation code

		body, err := validate.EnsureValidGeoJSON(req.Body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		validation_opts := validate.DefaultValidateOptions()

		err = validate.ValidateWithOptions(body, validation_opts)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		for _, pr := range opts.CustomProperties {

			ok, err := browser_properties.EnsureCustomPropertyHasValue(ctx, pr, body)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			if !ok {
				msg := fmt.Sprintf("Required property '%s' has invalid or missing value", pr.Name())
				http.Error(rsp, msg, http.StatusBadRequest)
				return
			}

		}

		// END OF validation code

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
			Account:    acct,
		}

		final, err := publishFeature(ctx, publish_opts, new_body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.WriteHeader(http.StatusCreated)
		rsp.Write(final)
		return
	}

	return http.HandlerFunc(fn), nil
}
