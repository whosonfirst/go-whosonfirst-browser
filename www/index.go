package www

import (
	"errors"
	"html/template"
	_ "log"
	"net/http"
)

type IndexHandlerOptions struct {
	Templates *template.Template
	Endpoints *Endpoints
}

type IndexVars struct {
	Endpoints *Endpoints
}

func IndexHandler(opts IndexHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("index")

	if t == nil {
		return nil, errors.New("Missing index template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		vars := IndexVars{
			Endpoints: opts.Endpoints,
		}

		RenderTemplate(rsp, t, vars)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
