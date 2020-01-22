package http

import (
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/editor"
	"github.com/whosonfirst/go-whosonfirst-browser/export"
	"github.com/whosonfirst/go-writer"
	_ "log"
	gohttp "net/http"
	"time"
)

func CessationHandler(r reader.Reader, wr writer.Writer, ed *editor.Editor) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		switch req.Method {
		case "POST":
			// pass
		default:
			gohttp.Error(rsp, "Method not allowed.", gohttp.StatusMethodNotAllowed)
			return
		}

		foo, err, status := FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		f := foo.Feature

		err = req.ParseMultipartForm(1024) // something something something... maybe?

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
		}

		date := req.FormValue("edtf:cessation")

		var t time.Time

		if date != "" {

			date_t, err := time.Parse("2006-01-02", date)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			t = date_t

		} else {
			t = time.Now()
		}

		ctx := req.Context()
		body := f.Bytes()

		updated_body, _, err := ed.CessateFeature(ctx, body, t)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		exported_body, err := export.ExportFeature(ctx, wr, updated_body)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		WriteGeoJSONHeaders(rsp)

		rsp.Write(exported_body)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
