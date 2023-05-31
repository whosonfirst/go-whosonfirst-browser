package browser

import (
	"sync"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	wasm_exec "github.com/sfomuseum/go-http-wasm/v2"	
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

	maps_opts = maps.DefaultMapsOptions()
	maps_opts.AppendJavaScriptAtEOF = cfg.JavaScriptAtEOF
	maps_opts.RollupAssets = capabilities.RollupAssets
	maps_opts.Prefix = uris_table.URIPrefix
	maps_opts.Logger = logger
	
	wasm_exec_opts = wasm_exec.DefaultWASMOptions()
	wasm_exec_opts.AppendJavaScriptAtEOF = cfg.JavaScriptAtEOF
	wasm_exec_opts.RollupAssets = cfg.RollupAssets
	wasm_exec_opts.Prefix = uris_table.URIPrefix
	wasm_exec_opts.Logger = logger
	
}
