package www

import (
	"github.com/aaronland/go-http-maps/static"
	aa_static "github.com/aaronland/go-http-static"
	_ "log"
	"net/http"
)

func AppendStaticAssetHandlers(mux *http.ServeMux) error {
	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, "")
}

func AppendStaticAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {
	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
