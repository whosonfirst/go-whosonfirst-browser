package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	aa_log "github.com/aaronland/go-log"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	browser_custom "github.com/whosonfirst/go-whosonfirst-browser/v7/custom"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	browser_properties "github.com/whosonfirst/go-whosonfirst-browser/v7/properties"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
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
	CustomValidationFunc  browser_custom.CustomValidationFunc
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

		/*

			Ideally we'd like to be able to do something like this:

			// github.com/aaronland/go-wasi

			_, err = wasi.Run(ctx, opts.CustomValidationWASM, string(body))

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusBadRequest)
				return
			}

			But practically it is... impractical. Specifically, there is not
			broad support or parity between WASI and WASM stuff. Further in
			order to compile Go code in to a WASI binary you need to use tinygo
			which lacks support for much (most) of the standard Go library, for
			example encoding/json. While it would be hypotherically possible to
			run WASM produced with GOOS=JS here using the v8go package the problem
			with that approach is a) it's a proposed extension to the v8go package
			that requires monkey-patching a bunch of things (https://github.com/rogchap/v8go/issues/333)
			and b) I haven't been able to make it work and c) we would then need
			to check whether we're runing v8/GOOS=JS WASM or plain vanilla WASM.

			So, in the end it's easier (emphasis on the "-ier") to expect that server-
			side validation code be written in Go since that same code can be compiled
			in to WASM that can be run client-side. It's a bit of a nuisance from a
			configuration / scaffolding perspective but that's just where we're at
			with WASM today...

			(20230201/thisisaaronland)
		*/

		if opts.CustomValidationFunc != nil {

			err = opts.CustomValidationFunc(body)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// END OF validation code

		_, new_body, err := opts.PointInPolygonService.Update(ctx, body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		name, err := properties.Name(new_body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		// To do (maybe): Make these customizable with URI templates in opts
		// https://pkg.go.dev/github.com/jtacoma/uritemplates

		title := fmt.Sprintf("Create new record for '%s'", name)
		description := title

		now := time.Now()
		ts := now.Unix()

		branch := fmt.Sprintf("create-feature-%d", ts)

		publish_opts := &publishFeatureOptions{
			Logger:      opts.Logger,
			WriterURIs:  opts.WriterURIs,
			Exporter:    opts.Exporter,
			Cache:       opts.Cache,
			Account:     acct,
			Title:       title,
			Description: description,
			Branch:      branch,
		}

		final, err := publishFeature(ctx, publish_opts, new_body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.WriteHeader(http.StatusCreated)
		rsp.Write(final)

		aa_log.Debug(opts.Logger, "Wrote new feature to %s branch", branch)
		return
	}

	return http.HandlerFunc(fn), nil
}
