package browser

import (
	"html/template"
	"sync"
	"text/template"

	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader-cachereader"	
	"github.com/aaronland/go-http-maps"	
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"	
)

var logger *log.Logger

var uris_table *URIs
var capabilities *Capabilities
var config *Config

var wof_reader reader.Reader
var wof_cache cache.Cache
var wof_exporter export.Exporter
var wof_writer_uris []string

var pointinpolygon_service *pointinpolygon.PointInPolygonService
var search_database fulltext.FullTextDatabase

var authenticator auth.Authenticator
var map_provider provider.Provider

var html_t *html_template.Template
var js_t *text_template.Template

var cors_wrapper *cors.Cors

var www_opts *www.BrowserOptions
var bootstrap_opts *bootstrap.BootstrapOptions
var maps_opts *maps.MapsOptions

var setupStaticOnce sync.Once
var setupStaticError error

var setupWWWOnce sync.Once
var setupWWWError error

var setupPointInPolygonOnce sync.Once
var setupPointInPolygonError error

func setupStatic() {

	www_opts = www.DefaultBrowserOptions()
	www_opts.AppendJavaScriptAtEOF = settings.JavaScriptAtEOF
	www_opts.RollupAssets = capabilities.RollupAssets
	www_opts.Prefix = uris_table.URIPrefix
	www_opts.Logger = logger
	www_opts.DataAttributes["whosonfirst-uri-endpoint"] = uris_table.GeoJSON

	bootstrap_opts = bootstrap.DefaultBootstrapOptions()
	bootstrap_opts.AppendJavaScriptAtEOF = settings.JavaScriptAtEOF
	bootstrap_opts.RollupAssets = capabilities.RollupAssets
	bootstrap_opts.Prefix = uris_table.URIPrefix
	bootstrap_opts.Logger = logger
}

func setupWWW() {

	var err error

	html_t, err = html.LoadTemplates(ctx, settings.Templates...)

	if err != nil {
		setupWWWError = fmt.Errorf("Failed to load templates, %w", err)
		return
	}

	js_t, err = text.LoadTemplatesMatching(ctx, "*.js", javascript.FS)

	if err != nil {
		setupWWWError = fmt.Errorf("Failed to load JS templates, %w", err)
		return
	}

}

func setupPointInPolygon() {

	spatial_db, err := database.NewSpatialDatabase(ctx, cfg.SpatialDatabaseURI)

	if err != nil {
		setupPointInPolygonError = fmt.Errorf("Failed to create spatial database, %w", err)
		return
	}

	pt_definition, err := placetypes.NewDefinition(ctx, cfg.PlacetypesDefinitionURI)

	if err != nil {
		setupPointInPolygonError = fmt.Errorf("Failed to create placetypes definition, %w", err)
		return
	}

	pip_options := &pointinpolygon.PointInPolygonServiceOptions{
		SpatialDatabase:      spatial_db,
		ParentReader:         cr,
		PlacetypesDefinition: pt_definition,
		Logger:               logger,
		SkipPlacetypeFilter:  cfg.PointInPolygonSkipPlacetypeFilter,
	}

	pointinpolygon_service, err := pointinpolygon.NewPointInPolygonServiceWithOptions(ctx, pip_options)

	if err != nil {
		setupPointInPolygonError = fmt.Errorf("Failed to create point in polygon service, %w", err)
		return
	}

}
