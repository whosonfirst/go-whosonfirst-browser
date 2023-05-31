package browser

import (
	html_template "html/template"
	"io/fs"
	"log"
	text_template "text/template"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	"github.com/aaronland/go-http-maps/provider"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
)

// Cloned from corresponding RunOptions variables

var cfg *Config
var logger *log.Logger
var templates_fs []fs.FS

// Set up in browser.go

var uris_table *uris.URIs
var capabilities *browser_capabilities.Capabilities

// setup_whosonfirst.go

var wof_reader reader.Reader
var wof_cache cache.Cache
var wof_exporter export.Exporter
var wof_writer_uris []string

// setup_pointinpolygon.go

var pointinpolygon_service *pointinpolygon.PointInPolygonService

// setup_search.go

var search_database fulltext.FullTextDatabase

// setup_auth.go

var null_authenticator auth.Authenticator
var authenticator auth.Authenticator

// setup_www.go

var html_t *html_template.Template
var js_t *text_template.Template

// setup_maps.go

var map_provider provider.Provider

// setup_cors.go

var cors_wrapper *cors.Cors

// setup_static.go

var www_opts *www.BrowserOptions
var bootstrap_opts *bootstrap.BootstrapOptions
var maps_opts *maps.MapsOptions
