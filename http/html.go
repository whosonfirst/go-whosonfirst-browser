package http

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-render/assets/html"
	"github.com/whosonfirst/go-whosonfirst-render/reader"
	"html/template"
	"log"
	gohttp "net/http"
)

func HTMLHandler(r reader.Reader) (gohttp.Handler, error) {

	tpl, err := html.Asset("templates/html/spr.html")

	if err != nil {
		return nil, err
	}

	t := template.New("name")

	t, err = t.Parse(string(tpl))

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		path := req.URL.Path

		log.Println("PATH", path)
		
		fh, err := r.Read(path)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		f, err := feature.LoadFeatureFromReader(fh)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		s, err := f.SPR()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		err = t.Execute(rsp, s)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
