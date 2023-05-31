package browser

import (
	"context"
	"fmt"
	html_template "html/template"
	"log"
	"sync"
	text_template "text/template"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	"github.com/aaronland/go-http-maps/provider"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-http-auth"
	"github.com/sfomuseum/go-template/html"
	"github.com/sfomuseum/go-template/text"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-cachereader"
	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/templates/javascript"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
)

var logger *log.Logger

var uris_table *uris.URIs
var capabilities *browser_capabilities.Capabilities
var cfg *Config

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

var setupCORSOnce sync.Once
var setupCORSError error

var setupPointInPolygonOnce sync.Once
var setupPointInPolygonError error

var setupAuthenticatorOnce sync.Once
var setupAuthenticatorError error

var setupWhosOnFirstReaderOnce sync.Once
var setupWhosOnFirstReaderError error

var setupWhosOnFirstWriterOnce sync.Once
var setupWhosOnFirstWriterError error

var setupSearchOnce sync.Once
var setupSearchError error

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

func setupCORS() {

	if !cfg.EnableCORS {
		return
	}

	cors_origins := cfg.CORSOrigins

	if len(cors_origins) == 0 {
		cors_origins = []string{"*"}
	}

	cors_wrapper = cors.New(cors.Options{
		AllowedOrigins:   cors_origins,
		AllowCredentials: cfg.CORSAllowCredentials,
	})
}

func setupWWW() {

	ctx := context.Background()
	var err error

	html_t, err = html.LoadTemplates(ctx, cfg.Templates...)

	if err != nil {
		setupWWWError = fmt.Errorf("Failed to load templates, %w", err)
		return
	}

	js_t, err = text.LoadTemplatesMatching(ctx, "*.js", javascript.FS)

	if err != nil {
		setupWWWError = fmt.Errorf("Failed to load JS templates, %w", err)
		return
	}

	map_provider, err = provider.NewProvider(ctx, cfg.MapProviderURI)

	if err != nil {
		setupWWWError = fmt.Errorf("Failed to create new map provider, %w", err)
		return
	}

	setupCORSOnce.Do(setupCORS)
}

func setupPointInPolygon() {

	ctx := context.Background()

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
		SpatialDatabase: spatial_db,
		// FIX ME
		// ParentReader:         ...
		PlacetypesDefinition: pt_definition,
		Logger:               logger,
		SkipPlacetypeFilter:  cfg.PointInPolygonSkipPlacetypeFilter,
	}

	pointinpolygon_service, err = pointinpolygon.NewPointInPolygonServiceWithOptions(ctx, pip_options)

	if err != nil {
		setupPointInPolygonError = fmt.Errorf("Failed to create point in polygon service, %w", err)
		return
	}

}

func setupAuthenticator() {

	var err error
	ctx := context.Background()

	authenticator, err = auth.NewAuthenticator(ctx, cfg.AuthenticatorURI)

	if err != nil {
		setupAuthenticatorError = err
		return
	}
}

func setupWhosOnFirstReader() {

	ctx := context.Background()

	browser_reader, err := reader.NewMultiReaderFromURIs(ctx, cfg.ReaderURIs...)

	if err != nil {
		setupWhosOnFirstReaderError = fmt.Errorf("Failed to create reader, %w", err)
		return
	}

	browser_cache, err := cache.NewCache(ctx, cfg.CacheURI)

	if err != nil {
		setupWhosOnFirstReaderError = fmt.Errorf("Failed to create new cache, %w", err)
		return
	}

	cr_opts := &cachereader.CacheReaderOptions{
		Reader: browser_reader,
		Cache:  browser_cache,
	}

	cr, err := cachereader.NewCacheReaderWithOptions(ctx, cr_opts)

	if err != nil {
		setupWhosOnFirstReaderError = fmt.Errorf("Failed to create cache reader, %w", err)
		return
	}

	wof_reader = cr
	wof_cache = browser_cache
}

func setupWhosOnFirstWriter() {

	ctx := context.Background()
	var err error

	wof_exporter, err = export.NewExporter(ctx, cfg.ExporterURI)

	if err != nil {
		setupWhosOnFirstReaderError = fmt.Errorf("Failed to create new exporter, %w", err)
		return
	}

}

func setupSearch() {

	ctx := context.Background()
	var err error

	search_database, err = fulltext.NewFullTextDatabase(ctx, cfg.SearchDatabaseURI)

	if err != nil {
		setupSearchError = fmt.Errorf("Failed to create fulltext database for '%s', %w", cfg.SearchDatabaseURI, err)
		return
	}

}
