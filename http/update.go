package http

import (
	"bufio"
	"bytes"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-reader"
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
	AllowedPath *regexp.Regexp
}

func UpdateHandler(r reader.Reader, wr writer.Writer, opts UpdateHandlerOptions) (gohttp.Handler, error) {

	ex_opts, err := options.NewDefaultOptions()

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		ctx := req.Context()

		err := req.ParseForm()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
		}

		f, err, _ := FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
		}

		body := f.Bytes()

		// MOVE THIS IN TO A GENERIC update/update.go PACKAGE

		for path, value := range req.Form {

			if !opts.AllowedPath.MatchString(path) {

				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			// FIX ME : SANITIZE value HERE

			body, err = sjson.SetBytes(body, path, value)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			}
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

		br := bytes.NewReader(ex_body)
		fh := ioutil.NopCloser(br)

		err = wr.Write(ctx, rel_path, fh)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(ex_body)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
