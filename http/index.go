package http

import (
	"errors"
	"html/template"
	_ "log"
	gohttp "net/http"
)

type IndexHandlerOptions struct {
	Templates *template.Template
	Endpoints *Endpoints
}

type IndexVars struct {
	Endpoints *Endpoints
}

func IndexHandler(opts IndexHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("index")

	if t == nil {
		return nil, errors.New("Missing index template")
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		vars := IndexVars{
			Endpoints: opts.Endpoints,
		}

		RenderTemplate(rsp, t, vars)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
