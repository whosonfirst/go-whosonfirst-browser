package www

import (
	"errors"
	"html/template"
	_ "log"
	"net/http"

	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
)

type IndexHandlerOptions struct {
	Templates    *template.Template
	URIs         *browser_uris.URIs
	Capabilities *browser_capabilities.Capabilities
}

type IndexVars struct {
	Paths        *browser_uris.URIs
	Capabilities *browser_capabilities.Capabilities
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
