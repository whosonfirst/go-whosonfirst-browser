package http

import (
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/whosonfirst/go-whosonfirst-static/assets/html"
	"github.com/whosonfirst/go-whosonfirst-static/utils"
	"html/template"
	"log"
	gohttp "net/http"
	"time"
)

type HTMLOptions struct {
	MapzenAPIKey string
}

type HTMLVars struct {
	SPR          spr.StandardPlacesResult
	LastModified string
	MapzenAPIKey string
}

func NewDefaultHTMLOptions() HTMLOptions {

	opts := HTMLOptions{
		MapzenAPIKey: "mapzen-xxxxxx",
	}

	return opts
}

func HTMLHandler(r reader.Reader, opts HTMLOptions) (gohttp.Handler, error) {

	tpl, err := html.Asset("templates/html/id.html")

	if err != nil {
		return nil, err
	}

	t := template.New("name")

	t, err = t.Parse(string(tpl))

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		log.Println("REQUEST", req.URL.Path)

		f, err, status := utils.FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		s, err := f.SPR()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		now := time.Now()
		lastmod := now.Format(time.RFC3339)

		vars := HTMLVars{
			SPR:          s,
			LastModified: lastmod,
			MapzenAPIKey: opts.MapzenAPIKey,
		}

		err = t.Execute(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
