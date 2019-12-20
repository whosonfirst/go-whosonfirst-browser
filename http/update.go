package http

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export"
	"github.com/whosonfirst/go-whosonfirst-export/options"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/go-writer"
	"io/ioutil"
	"log"
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
			gohttp.Error(rsp, "NOPE", gohttp.StatusInternalServerError)
			return
		}

		f, err, _ := FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		body := f.Bytes()

		max := int64(1024 * 1024 * 10) // sudo make me an option
		err = req.ParseMultipartForm(max)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		if len(req.PostForm) == 0 {
			gohttp.Error(rsp, "No updates", gohttp.StatusInternalServerError)
			return
		}

		updates := 0

		// MOVE THIS IN TO A GENERIC update/update.go PACKAGE

		// can we (should we) do this concurrently?

		for path, value := range req.PostForm {

			log.Println("PATH", path)

			if !opts.AllowedPaths.MatchString(path) {

				gohttp.Error(rsp, "Invalid path", gohttp.StatusBadRequest)
				return
			}

			// TO DO : SANITIZE value HERE
			// TO DO: CHECK WHETHER PROPERTY/VALUE IS A SINGLETON - FOR EXAMPLE:
			// curl -s -X POST -F 'properties.wof:name=SPORK' -F 'properties.wof:name=BOB' http://localhost:8080/update/101736545 | \
			// python -mjson.tool | grep 'wof:name'
			// "wof:name": [

			var new_value interface{}

			switch len(value) {
			case 0:
				gohttp.Error(rsp, "Invalid value", gohttp.StatusBadRequest)
				return
			case 1:
				new_value = value[0]
			default:
				new_value = value
			}

			old_rsp := gjson.GetBytes(body, path)

			if old_rsp.Exists() {

				old_enc, err := json.Marshal(old_rsp.Value())

				if err != nil {
					gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
					return
				}

				new_enc, err := json.Marshal(new_value)

				if err != nil {
					gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
					return
				}

				if bytes.Compare(new_enc, old_enc) == 0 {
					log.Println("SAME SAME")
					continue
				}
			}

			log.Println("SET", path, new_value)
			body, err = sjson.SetBytes(body, path, new_value)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			updates += 1
		}

		if updates == 0 {
			log.Println("NO UPDATES")
			WriteGeoJSONHeaders(rsp)
			rsp.Write(body)
			return
		}

		var buf bytes.Buffer
		bw := bufio.NewWriter(&buf)

		err = export.Export(body, ex_opts, bw)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		bw.Flush()

		ex_body := buf.Bytes()

		id_rsp := gjson.GetBytes(ex_body, "properties.wof:id")

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

		ctx := req.Context()

		br := bytes.NewReader(ex_body)
		fh := ioutil.NopCloser(br)

		err = wr.Write(ctx, rel_path, fh)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		WriteGeoJSONHeaders(rsp)

		rsp.Write(ex_body)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
