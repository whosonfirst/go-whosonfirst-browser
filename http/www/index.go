package www

import (
	"errors"
	"html/template"
	_ "log"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
)

type IndexHandlerOptions struct {
	Templates    *template.Template
	URIs         *uris.URIs
	Capabilities *capabilities.Capabilities
}

type IndexVars struct {
	Paths        *uris.URIs
	Capabilities *capabilities.Capabilities
}

func IndexHandler(opts IndexHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("index")

	if t == nil {
		return nil, errors.New("Missing index template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		vars := IndexVars{
			Paths:        opts.URIs,
			Capabilities: opts.Capabilities,
		}

		RenderTemplate(rsp, t, vars)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
