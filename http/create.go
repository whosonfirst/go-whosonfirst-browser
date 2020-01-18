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

func CreateHandler(r reader.Reader, wr writer.Writer, ed *editor.Editor) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		switch req.Method {
		case "PUT":
			// pass
		default:
			gohttp.Error(rsp, "Method not allowed.", gohttp.StatusMethodNotAllowed)
			return
		}

		var update_req *editor.UpdateRequest

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&update_req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		ctx := req.Context()

		updated_body, _, err := ed.CreateFeature(ctx, update_req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		exported_body, err := export.ExportFeature(ctx, wr, updated_body)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		// STATUS 201...
		// PASS NEW ID IN HEADER?

		WriteGeoJSONHeaders(rsp)

		rsp.Write(exported_body)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
