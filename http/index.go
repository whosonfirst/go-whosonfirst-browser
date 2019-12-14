package http

import (
	"errors"
	"html/template"
	_ "log"
	gohttp "net/http"
)

type IndexHandlerOptions struct {
	Templates *template.Template
}

func IndexHandler(opts IndexHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("index")

	if t == nil {
		return nil, errors.New("Missing index template")
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		err := t.Execute(rsp, nil)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
