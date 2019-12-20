package http

import (
	gohttp "net/http"
)

func WriteGeoJSONHeaders(rsp gohttp.ResponseWriter) {
	rsp.Header().Set("Content-Type", "application/json")
	rsp.Header().Set("Access-Control-Allow-Origin", "*")
}

func WriteSVGHeaders(rsp gohttp.ResponseWriter) {
	rsp.Header().Set("Content-Type", "application/json")
	rsp.Header().Set("Content-Type", "image/svg+xml")
}
