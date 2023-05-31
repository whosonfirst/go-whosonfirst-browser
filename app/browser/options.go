package browser

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/aaronland/go-http-server/handler"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/chrome"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/custom"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/properties"
	wof_hierarchy "github.com/whosonfirst/go-whosonfirst-spatial-hierarchy"
)

type AssetHandlerFunc func(*http.ServeMux, string) error
type MiddlewareHandlerFunc func(http.Handler) http.Handler

type PointInPolygonOptions struct {
	ResultsCallback          wof_hierarchy.PointInPolygonHierarchyResolverUpdateCallback
	UpdateCallback           wof_hierarchy.PointInPolygonHierarchyResolverUpdateCallback
	ToCopyFromParentOnUpdate []string
}

type RunOptions struct {
	Logger    *log.Logger
	Config    *Config
	Templates []fs.FS

	CustomURIs              map[string]string
	CustomChrome            chrome.Chrome
	CustomRouteHandlerFuncs map[string]handler.RouteHandlerFunc
	// CustomMiddlewareHandlers    map[string][]MiddlewareHandlerFunc
	CustomAssetHandlerFunctions []AssetHandlerFunc
	CustomEditProperties        []properties.CustomProperty
	CustomEditValidationFunc    custom.CustomValidationFunc
	CustomEditValidationWasm    *custom.CustomValidationWasm
	PointInPolygonOptions       *PointInPolygonOptions
}
