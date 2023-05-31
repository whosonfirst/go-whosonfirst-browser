package browser

import (
	"sync"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
)

var setupStaticOnce sync.Once
var setupStaticError error

func setupStatic() {

	www_opts = www.DefaultBrowserOptions()
	www_opts.AppendJavaScriptAtEOF = cfg.JavaScriptAtEOF
	www_opts.RollupAssets = capabilities.RollupAssets
	www_opts.Prefix = uris_table.URIPrefix
	www_opts.Logger = logger
	www_opts.DataAttributes["whosonfirst-uri-endpoint"] = uris_table.GeoJSON

	bootstrap_opts = bootstrap.DefaultBootstrapOptions()
	bootstrap_opts.AppendJavaScriptAtEOF = cfg.JavaScriptAtEOF
	bootstrap_opts.RollupAssets = capabilities.RollupAssets
	bootstrap_opts.Prefix = uris_table.URIPrefix
	bootstrap_opts.Logger = logger
}
