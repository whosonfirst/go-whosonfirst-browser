package http

import (
	"github.com/protomaps/go-pmtiles/pmtiles"
	"log"
	gohttp "net/http"
	"time"
)

// TileHandler provides a `http.Handler` for serving Protomaps tile requests using 'loop'.
func TileHandler(loop *pmtiles.Loop, logger *log.Logger) gohttp.Handler {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		start := time.Now()

		status_code, headers, body := loop.Get(req.Context(), req.URL.Path)

		for k, v := range headers {
			rsp.Header().Set(k, v)
		}

		rsp.WriteHeader(status_code)
		rsp.Write(body)

		go logger.Printf("[%d] served %s in %s", status_code, req.URL.Path, time.Since(start))
	}

	return gohttp.HandlerFunc(fn)
}
