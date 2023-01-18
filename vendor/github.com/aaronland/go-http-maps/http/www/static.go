package www

import (
	aa_static "github.com/aaronland/go-http-static"	
	"github.com/aaronland/go-http-maps/static"
	_ "log"
	"net/http"
)

func AppendStaticAssetHandlers(mux *http.ServeMux) error {
	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, "")
}

func AppendStaticAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {
	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
