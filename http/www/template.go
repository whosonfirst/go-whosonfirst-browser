package www

import (
	"html/template"
	gohttp "net/http"
)

func RenderTemplate(rsp gohttp.ResponseWriter, t *template.Template, vars interface{}) {

	err := t.Execute(rsp, vars)

	if err != nil {
		gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
		return
	}

	return
}
