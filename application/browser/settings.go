package browser

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aaronland/go-http-maps/provider"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-cachereader"
	github_reader "github.com/whosonfirst/go-reader-github"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/chrome"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/templates/html"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	github_writer "github.com/whosonfirst/go-writer-github/v3"
)

type Settings struct {
	Paths        *www.Paths
	Capabilities *www.Capabilities

	Reader reader.Reader
	Cache  cache.Cache

	WriterURIs            []string
	Exporter              export.Exporter
	PointInPolygonService *pointinpolygon.PointInPolygonService
	Authenticator         auth.Authenticator

	MapProvider provider.Provider

	Templates []fs.FS

	CustomChrome   chrome.Chrome
	CustomHandlers map[string]http.HandlerFunc

	CORSWrapper *cors.Cors

	NavPlaceMaxFeatures int

	SelectPattern *regexp.Regexp

	WebFingerHostname string

	Verbose bool
}

func SettingsFromConfig(ctx context.Context, cfg *Config, logger *log.Logger) (*Settings, error) {

	settings := &Settings{
		Templates: []fs.FS{html.FS},
		Verbose:   cfg.Verbose,
	}

	reader_uris := cfg.ReaderURIs
	writer_uris := cfg.WriterURIs
	cache_uri := cfg.CacheURI

	// Fetch and assign GitHub access tokens to reader/writer URIs

	if cfg.GitHubAccessTokenURI != "" {

		if cfg.GitHubReaderAccessTokenURI == "" {
			cfg.GitHubReaderAccessTokenURI = cfg.GitHubAccessTokenURI
		}

		if cfg.GitHubWriterAccessTokenURI == "" {
			cfg.GitHubWriterAccessTokenURI = cfg.GitHubAccessTokenURI
		}
	}

	if cfg.GitHubReaderAccessTokenURI != "" {

		for idx, r_uri := range reader_uris {

			r_uri, err := github_reader.EnsureGitHubAccessToken(ctx, r_uri, cfg.GitHubReaderAccessTokenURI)

			if err != nil {
				return nil, fmt.Errorf("Failed to ensure GitHub access token for '%s', %w", r_uri, err)
			}

			reader_uris[idx] = r_uri
		}
	}

	if cfg.GitHubReaderAccessTokenURI != "" {

		for idx, wr_uri := range writer_uris {

			wr_uri, err := github_writer.EnsureGitHubAccessToken(ctx, wr_uri, cfg.GitHubReaderAccessTokenURI)

			if err != nil {
				return nil, fmt.Errorf("Failed to ensure GitHub access token for '%s', %w", wr_uri, err)
			}

			writer_uris[idx] = wr_uri
		}

	}

	settings.WriterURIs = writer_uris

	// Set up reader and reader cache. Note that we create a "cachereader"
	// manually since we want to be able to purge records from the cache (assuming
	// the edit hooks are enabled)

	if cache_uri == "tmp://" {

		now := time.Now()
		prefix := fmt.Sprintf("go-whosonfirst-browser-%d", now.Unix())

		tempdir, err := ioutil.TempDir("", prefix)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive tmp dir, %w", err)
		}

		defer os.RemoveAll(tempdir)

		cache_uri = fmt.Sprintf("fs://%s", tempdir)
	}

	browser_reader, err := reader.NewMultiReaderFromURIs(ctx, reader_uris...)

	if err != nil {
		return nil, fmt.Errorf("Failed to create reader, %w", err)
	}

	browser_cache, err := cache.NewCache(ctx, cache_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new cache, %w", err)
	}

	cr_opts := &cachereader.CacheReaderOptions{
		Reader: browser_reader,
		Cache:  browser_cache,
	}

	cr, err := cachereader.NewCacheReaderWithOptions(ctx, cr_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create cache reader, %w", err)
	}

	settings.Reader = cr
	settings.Cache = browser_cache

	// Set up www.Paths and www.Capabilities structs for passing between handlers

	www_capabilities := &www.Capabilities{}

	www_paths := &www.Paths{
		URIPrefix: cfg.URIPrefix,
		Index:     cfg.PathIndex,
		Ping:      cfg.PathPing,
	}

	// URI prefix stuff

	if cfg.URIPrefix != "" {

		cfg.URIPrefix = strings.TrimRight(cfg.URIPrefix, "/")

		if !strings.HasPrefix(cfg.URIPrefix, "/") {
			return nil, fmt.Errorf("Invalid -static-prefix value")
		}

		path_index, err := url.JoinPath(cfg.URIPrefix, cfg.PathIndex)

		if err != nil {
			return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path_index)
		}

		path_ping, err = url.JoinPath(cfg.URIPrefix, cfg.PathPing)

		if err != nil {
			return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path_ping)
		}

		www_paths.Index = path_index
		www_paths.Ping = path_ping
	}

	if cfg.EnableGeoJSON {

		www_capabilities.GeoJSON = true
		www_paths.GeoJSON = cfg.PathGeoJSON

		if cfg.URIPrefix != "" {

			path_geojson, err := url.JoinPath(cfg.URIPrefix, cfg.PathGeoJSON)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathGeoJSON)
			}

			www_paths.GeoJSON = path_geojson
		}

	}

	if cfg.EnableGeoJSONLD {

		www_capabilities.GeoJSONLD = true
		www_paths.GeoJSONLD = cfg.PathGeoJSONLD

		if cfg.URIPrefix != "" {

			path_geojsonld, err := url.JoinPath(cfg.URIPrefix, cfg.PathGeoJSONLD)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathGeoJSONLD)
			}

			www_paths.GeoJSONLD = path_geojsonld
		}
	}

	if cfg.EnableSVG {

		www_capabilities.SVG = true
		www_paths.SVG = cfg.PathSVG

		if cfg.URIPrefix != "" {

			path_svg, err := url.JoinPath(cfg.URIPrefix, cfg.PathSVG)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathSVG)
			}

			www_paths.SVG = path_svg
		}

	}

	if cfg.EnablePNG {

		www_capabilities.PNG = true
		www_paths.PNG = cfg.PathPNG

		if cfg.URIPrefix != "" {

			path_png, err := url.JoinPath(cfg.URIPrefix, cfg.PathPNG)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathPNG)
			}

			www_paths.PNG = path_png
		}
	}

	if cfg.EnableSelect {

		if cfg.SelectPattern == "" {
			return nil, fmt.Errorf("Missing -select-pattern parameter.")
		}

		pat, err := regexp.Compile(select_pattern)

		if err != nil {
			return nil, fmt.Errorf("Failed to compile select pattern (%s), %w", cfg.SelectPattern, err)
		}

		settings.SelectPattern = pat

		www_capabilities.Select = true
		www_paths.Select = cfg.PathSelect

		if cfg.URIPrefix != "" {

			path_select, err := url.JoinPath(cfg.URIPrefix, cfg.PathSelect)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathSelect)
			}

			www_paths.Select = path_select
		}
	}

	if cfg.EnableNavPlace {

		settings.NavPlaceMaxFeatures = cfg.NavPlaceMaxFeatures

		www_capabilities.NavPlace = true
		www_paths.NavPlace = cfg.PathNavPlace

		if cfg.URIPrefix != "" {

			path_navplace, err := url.JoinPath(cfg.URIPrefix, cfg.PathNavPlace)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathNavPlace)
			}

			www_paths.NavPlace = path_navplace
		}

	}

	if cfg.EnableSPR {

		www_capabilities.SPR = true
		www_paths.SPR = cfg.PathSPR

		if cfg.URIPrefix != "" {

			path_spr, err := url.JoinPath(cfg.URIPrefix, cfg.PathSPR)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathSPR)
			}

			www_paths.SPR = path_spr
		}

	}

	if cfg.EnableHTML {

		www_capabilities.HTML = true
		www_paths.Id = cfg.PathId

		if cfg.URIPrefix != "" {

			path_id, err := url.JoinPath(cfg.URIPrefix, cfg.PathId)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathId)
			}

			www_paths.Id = path_id
		}

	}

	if cfg.EnableWebFinger {

		settings.WebFingerHostname = cfg.WebFingerHostname

		www_capabilities.WebFinger = true
		www_paths.WebFinger = cfg.PathWebFinger

		if cfg.URIPrefix != "" {

			path_webfinger, err := url.JoinPath(cfg.URIPrefix, cfg.PathWebFinger)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathWebFinger)
			}

			www_paths.WebFinger = path_webfinger
		}
	}

	if cfg.EnableEditUI {

		www_capabilities.CreateFeature = true
		www_capabilities.DeprecateFeature = true
		www_capabilities.CessateFeature = true
		www_capabilities.EditGeometry = true

		www_paths.CreateFeature = cfg.PathCreateFeature
		www_paths.EditGeometry = cfg.PathEditGeometry

		if cfg.URIPrefix != "" {

			path_create_feature, err := url.JoinPath(cfg.URIPrefix, cfg.PathCreateFeature)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathCreateFeature)
			}

			path_edit_geometry, err := url.JoinPath(cfg.URIPrefix, cfg.PathEditGeometry)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathEditGeometry)
			}

			www_paths.CreateFeature = path_create_feature
			www_paths.EditGeometry = path_edit_geometry
		}
	}

	if cfg.EnableEditAPI {

		www_capabilities.CreateFeatureAPI = true
		www_capabilities.DeprecateFeatureAPI = true
		www_capabilities.CessateFeatureAPI = true
		www_capabilities.EditGeometryAPI = true

		www_paths.CreateFeatureAPI = cfg.PathCreateFeatureAPI
		www_paths.DeprecateFeatureAPI = cfg.PathDeprecateFeatureAPI
		www_paths.CessateFeatureAPI = cfg.PathCessateFeatureAPI
		www_paths.EditGeometryAPI = cfg.PathEditGeometryAPI

		if cfg.URIPrefix != "" {

			path_api_create_feature, err := url.JoinPath(cfg.URIPrefix, cfg.PathCreateFeatureAPI)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathCreateFeatureAPI)
			}

			path_api_cessate_feature, err := url.JoinPath(cfg.URIPrefix, cfg.PathCessateFeatureAPI)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathCessateFeatureAPI)
			}

			path_api_deprecate_feature, err := url.JoinPath(cfg.URIPrefix, cfg.PathDeprecateFeatureAPI)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathDeprecateFeatureAPI)
			}

			path_api_edit_geometry, err := url.JoinPath(cfg.URIPrefix, cfg.PathEditGeometryAPI)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathEditGeometryAPI)
			}

			www_paths.CreateFeatureAPI = path_api_create_feature
			www_paths.DeprecateFeatureAPI = path_api_deprecate_feature
			www_paths.CessateFeatureAPI = path_api_cessate_feature
			www_paths.EditGeometryAPI = path_api_edit_geometry
		}

	}

	settings.Paths = www_paths
	settings.Capabilities = www_capabilities

	// Auth hooks

	authenticator, err := auth.NewAuthenticator(ctx, cfg.AuthenticatorURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to create authenticator, %w", err)
	}

	settings.Authenticator = authenticator

	// Map provider

	map_provider, err := provider.NewProvider(ctx, cfg.MapProviderURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new map provider, %w", err)
	}

	err = map_provider.SetLogger(logger)

	if err != nil {
		return nil, fmt.Errorf("Failed to set logger for provider, %w", err)
	}

	settings.MapProvider = map_provider

	// Custom chrome (this is still in flux)

	custom, err := chrome.NewChrome(ctx, cfg.CustomChromeURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to create custom chrome, %w", err)
	}

	settings.CustomChrome = custom

	// CORS

	if cfg.EnableCORS {

		cors_wrapper := cors.New(cors.Options{
			AllowedOrigins:   cfg.CORSOrigins,
			AllowCredentials: cfg.CORSAllowCredentials,
		})

		settings.CORSWrapper = cors_wrapper
	}

	if cfg.EnableEditAPI {

		// Exporter

		ex, err := export.NewExporter(ctx, cfg.ExporterURI)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new exporter, %w", err)
		}

		// START OF Point in polygon service setup

		// We're doing this the long way because we need/want to pass in 'cr' and I am
		// not sure what the interface/signature changes to do that should be...

		spatial_db, err := database.NewSpatialDatabase(ctx, cfg.SpatialDatabaseURI)

		if err != nil {
			return nil, fmt.Errorf("Failed to create spatial database, %w", err)
		}

		pip_service, err := pointinpolygon.NewPointInPolygonServiceWithDatabaseAndReader(ctx, spatial_db, cr)

		if err != nil {
			return nil, fmt.Errorf("Failed to create point in polygon service, %w", err)
		}

		// END OF Point in polygon service setup

		settings.Exporter = ex
		settings.PointInPolygonService = pip_service
	}

	return settings, nil
}
