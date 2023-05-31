package browser

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-browser/v7/chrome"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/custom"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/properties"
)

type AssetHandlerFunc func(*http.ServeMux, string) error
type MiddlewareHandlerFunc func(http.Handler) http.Handler

type RunOptions struct {
	Logger                      *log.Logger
	Config                      *Config
	Templates                   []fs.FS
	CustomChrome                chrome.Chrome
	CustomWWWHandlers           map[string]http.Handler
	CustomAPIHandlers           map[string]http.Handler
	CustomMiddlewareHandlers    map[string][]MiddlewareHandlerFunc
	CustomAssetHandlerFunctions []AssetHandlerFunc
	CustomEditProperties        []properties.CustomProperty
	CustomEditValidationFunc    custom.CustomValidationFunc
	CustomEditValidationWasm    *custom.CustomValidationWasm
}
