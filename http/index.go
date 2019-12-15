package http

import (
	"errors"
	"html/template"
	_ "log"
	gohttp "net/http"
)

type IndexHandlerOptions struct {
	Templates  *template.Template
	IdEndpoint string
}

type IndexVars struct {
	IdEndpoint string
}

func IndexHandler(opts IndexHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("index")

	if t == nil {
		return nil, errors.New("Missing index template")
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		vars := IndexVars{
			IdEndpoint: "",
		}

		err := t.Execute(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
