package www

import (
	"errors"
	"html/template"
	_ "log"
	"net/http"
)

type IndexHandlerOptions struct {
	Templates    *template.Template
	Paths        *Paths
	Capabilities *Capabilities
}

type IndexVars struct {
	Paths        *Paths
	Capabilities *Capabilities
}

func IndexHandler(opts IndexHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("index")

	if t == nil {
		return nil, errors.New("Missing index template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		vars := IndexVars{
			Paths:        opts.Paths,
			Capabilities: opts.Capabilities,
		}

		RenderTemplate(rsp, t, vars)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
