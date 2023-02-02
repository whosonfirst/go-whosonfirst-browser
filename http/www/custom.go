package www

import (
	"io"
	"net/http"

	browser_custom "github.com/whosonfirst/go-whosonfirst-browser/v7/custom"
)

type CustomValidationWasmHandlerOptions struct {
	CustomValidationWasm *browser_custom.CustomValidationWasm
}

func CustomValidationWasmHandler(opts *CustomValidationWasmHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		wasm_opts := opts.CustomValidationWasm

		r, err := wasm_opts.FS.Open(wasm_opts.Path)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Close()

		rsp.Header().Set("Content-type", "application/wasm")

		_, err = io.Copy(rsp, r)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
