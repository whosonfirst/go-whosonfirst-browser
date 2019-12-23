package http

// curl -s -X POST -d '{"properties":{"wof:name":"SPORK"}}' http://localhost:8080/update/101736545 | python -mjson.tool | grep 'wof:name'
// "wof:name": "SPORK",

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/update"
	"github.com/whosonfirst/go-whosonfirst-export"
	"github.com/whosonfirst/go-whosonfirst-export/options"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/go-writer"
	"io/ioutil"
	_ "log"
	gohttp "net/http"
	"regexp"
)

type UpdateHandlerOptions struct {
	AllowedPaths *regexp.Regexp // multiple regexps?
}

func UpdateHandler(r reader.Reader, wr writer.Writer, opts *UpdateHandlerOptions) (gohttp.Handler, error) {

	ex_opts, err := options.NewDefaultOptions()

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		switch req.Method {
		case "POST":
			// pass
		default:
			gohttp.Error(rsp, "Method not allowed.", gohttp.StatusMethodNotAllowed)
			return
		}

		f, err, status := FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		var update_req *update.Update

		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&update_req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		ctx := req.Context()
		body := f.Bytes()

		updated_body, updates, err := update.UpdateFeature(ctx, body, update_req, opts.AllowedPaths)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		if updates == 0 {
			WriteGeoJSONHeaders(rsp)
			rsp.Write(body)
			return
		}

		var buf bytes.Buffer
		bw := bufio.NewWriter(&buf)

		err = export.Export(updated_body, ex_opts, bw)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		bw.Flush()

		exported_body := buf.Bytes()

		id_rsp := gjson.GetBytes(exported_body, "properties.wof:id")

		if !id_rsp.Exists() {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		id := id_rsp.Int()

		rel_path, err := uri.Id2RelPath(id)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		br := bytes.NewReader(exported_body)
		fh := ioutil.NopCloser(br)

		err = wr.Write(ctx, rel_path, fh)

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
