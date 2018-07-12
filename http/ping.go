package http

import (
	gohttp "net/http"
)

func PingHandler() (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {
		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Write([]byte("PONG"))
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
