package http

// curl -s -X POST -d '{"properties":{"wof:name":"SPORK"}}' http://localhost:8080/update/101736545 | python -mjson.tool | grep 'wof:name'
// "wof:name": "SPORK",

import (
	"encoding/json"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/editor"
	"github.com/whosonfirst/go-whosonfirst-browser/export"
	"github.com/whosonfirst/go-writer"
	_ "log"
	gohttp "net/http"
)

func UpdateHandler(r reader.Reader, wr writer.Writer, ed *editor.Editor) (gohttp.Handler, error) {

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

		var update_req *editor.UpdateRequest

		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&update_req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		ctx := req.Context()
		body := f.Bytes()

		updated_body, update_rsp, err := ed.UpdateFeature(ctx, body, update_req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		if update_rsp.Count() == 0 {
			WriteGeoJSONHeaders(rsp)
			rsp.Write(body)
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
