package www

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
)

type URIsHandlerOptions struct {
	URIs     *uris.URIs
	Template *template.Template
}

type URIsHandlerVars struct {
	URIs string
}

func URIsHandler(opts *URIsHandlerOptions) (http.Handler, error) {

	enc_uris, err := json.Marshal(opts.URIs)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal URIs, %w", err)
	}

	vars := URIsHandlerVars{
		URIs: string(enc_uris),
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Header().Set("Content-type", "text/javascript")

		err := opts.Template.Execute(rsp, vars)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
