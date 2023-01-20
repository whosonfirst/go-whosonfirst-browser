package www

import (
	aa_static "github.com/aaronland/go-http-static"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/static"
	"net/http"
)

func AppendStaticAssetHandlers(mux *http.ServeMux) error {
	return aa_static.AppendStaticAssetHandlers(mux, static.FS)
}

func AppendStaticAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {
	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
