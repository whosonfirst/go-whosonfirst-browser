package api

import (
	"bytes"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
	wof_http "github.com/whosonfirst/go-whosonfirst-browser/v6/http"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-writer/v3"
	"log"
	"net/http"
)

type CessateFeatureHandlerOptions struct {
	Reader        reader.Reader
	Writer        writer.Writer
	Exporter      export.Exporter
	Authenticator auth.Authenticator
	Logger        *log.Logger
}

func CessateFeatureHandler(opts *CessateFeatureHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

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

		br := bytes.NewReader(new_body)

		_, err = opts.Writer.Write(ctx, uri.URI, br)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
