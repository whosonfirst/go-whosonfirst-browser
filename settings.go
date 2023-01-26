package browser

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
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
	browser_capabilities "github.com/whosonfirst/go-whosonfirst-browser/v7/capabilities"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/chrome"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/templates/html"
	browser_uris "github.com/whosonfirst/go-whosonfirst-browser/v7/uris"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	github_writer "github.com/whosonfirst/go-writer-github/v3"
)

type Settings struct {
	Authenticator         auth.Authenticator
	Cache                 cache.Cache
	Capabilities          *browser_capabilities.Capabilities
	CORSWrapper           *cors.Cors
	CustomChrome          chrome.Chrome
	CustomHandlers        map[string]http.HandlerFunc
	Exporter              export.Exporter
	MapProvider           provider.Provider
	NavPlaceMaxFeatures   int
	URIs                  *browser_uris.URIs
	PointInPolygonService *pointinpolygon.PointInPolygonService
	Reader                reader.Reader
	SearchDatabase        fulltext.FullTextDatabase
	SelectPattern         *regexp.Regexp
	Templates             []fs.FS
	Verbose               bool
	WebFingerHostname     string
	WriterURIs            []string
}

func SettingsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*Settings, error) {

	cfg, err := ConfigFromFlagSet(ctx, fs)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive config from flagset, %w", err)
	}

	return SettingsFromConfig(ctx, cfg)
}

func SettingsFromConfig(ctx context.Context, cfg *Config) (*Settings, error) {

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

	capabilities := &browser_capabilities.Capabilities{}

	uris := &browser_uris.URIs{
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

		uris.Index = path_index
		uris.Ping = path_ping
	}

	if cfg.EnableGeoJSON {

		capabilities.GeoJSON = true
		uris.GeoJSON = cfg.PathGeoJSON
		uris.GeoJSONAlt = cfg.PathGeoJSONAlt

		if cfg.URIPrefix != "" {

			path_geojson, err := url.JoinPath(cfg.URIPrefix, cfg.PathGeoJSON)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathGeoJSON, err)
			}

			uris.GeoJSON = path_geojson

			alt_paths := make([]string, len(cfg.PathGeoJSONAlt))

			for idx, path := range cfg.PathGeoJSONAlt {

				path, err := url.JoinPath(cfg.URIPrefix, path)

				if err != nil {
					return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path, err)
				}

				alt_paths[idx] = path
			}

			uris.GeoJSONAlt = alt_paths
		}

	}

	if cfg.EnableGeoJSONLD {

		capabilities.GeoJSONLD = true
		uris.GeoJSONLD = cfg.PathGeoJSONLD
		uris.GeoJSONLDAlt = cfg.PathGeoJSONLDAlt

		if cfg.URIPrefix != "" {

			path_geojsonld, err := url.JoinPath(cfg.URIPrefix, cfg.PathGeoJSONLD)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathGeoJSONLD, err)
			}

			uris.GeoJSONLD = path_geojsonld

			alt_paths := make([]string, len(cfg.PathGeoJSONLDAlt))

			for idx, path := range cfg.PathGeoJSONLDAlt {

				path, err := url.JoinPath(cfg.URIPrefix, path)

				if err != nil {
					return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path, err)
				}

				alt_paths[idx] = path
			}

			uris.GeoJSONLDAlt = alt_paths
		}
	}

	if cfg.EnableSVG {

		capabilities.SVG = true
		uris.SVG = cfg.PathSVG
		uris.SVGAlt = cfg.PathSVGAlt

		if cfg.URIPrefix != "" {

			path_svg, err := url.JoinPath(cfg.URIPrefix, cfg.PathSVG)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathSVG)
			}

			uris.SVG = path_svg

			alt_paths := make([]string, len(cfg.PathSVGAlt))

			for idx, path := range cfg.PathSVGAlt {

				path, err := url.JoinPath(cfg.URIPrefix, path)

				if err != nil {
					return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path, err)
				}

				alt_paths[idx] = path
			}

			uris.SVGAlt = alt_paths

		}

	}

	if cfg.EnablePNG {

		capabilities.PNG = true
		uris.PNG = cfg.PathPNG
		uris.PNGAlt = cfg.PathPNGAlt

		if cfg.URIPrefix != "" {

			path_png, err := url.JoinPath(cfg.URIPrefix, cfg.PathPNG)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathPNG)
			}

			uris.PNG = path_png

			alt_paths := make([]string, len(cfg.PathPNGAlt))

			for idx, path := range cfg.PathPNGAlt {

				path, err := url.JoinPath(cfg.URIPrefix, path)

				if err != nil {
					return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path, err)
				}

				alt_paths[idx] = path
			}

			uris.PNGAlt = alt_paths

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

		capabilities.Select = true
		uris.Select = cfg.PathSelect

		if cfg.URIPrefix != "" {

			path_select, err := url.JoinPath(cfg.URIPrefix, cfg.PathSelect)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathSelect)
			}

			uris.Select = path_select
		}
	}

	if cfg.EnableNavPlace {

		settings.NavPlaceMaxFeatures = cfg.NavPlaceMaxFeatures

		capabilities.NavPlace = true
		uris.NavPlace = cfg.PathNavPlace
		uris.NavPlaceAlt = cfg.PathNavPlaceAlt

		if cfg.URIPrefix != "" {

			path_navplace, err := url.JoinPath(cfg.URIPrefix, cfg.PathNavPlace)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathNavPlace)
			}

			uris.NavPlace = path_navplace

			alt_paths := make([]string, len(cfg.PathNavPlaceAlt))

			for idx, path := range cfg.PathNavPlaceAlt {

				path, err := url.JoinPath(cfg.URIPrefix, path)

				if err != nil {
					return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path, err)
				}

				alt_paths[idx] = path
			}

			uris.NavPlaceAlt = alt_paths

		}

	}

	if cfg.EnableSPR {

		capabilities.SPR = true
		uris.SPR = cfg.PathSPR
		uris.SPRAlt = cfg.PathSPRAlt

		if cfg.URIPrefix != "" {

			path_spr, err := url.JoinPath(cfg.URIPrefix, cfg.PathSPR)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathSPR)
			}

			uris.SPR = path_spr

			alt_paths := make([]string, len(cfg.PathSPRAlt))

			for idx, path := range cfg.PathSPRAlt {

				path, err := url.JoinPath(cfg.URIPrefix, path)

				if err != nil {
					return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path, err)
				}

				alt_paths[idx] = path
			}

			uris.SPRAlt = alt_paths

		}

	}

	if cfg.EnableHTML {

		capabilities.HTML = true
		uris.Id = cfg.PathId

		if cfg.URIPrefix != "" {

			path_id, err := url.JoinPath(cfg.URIPrefix, cfg.PathId)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathId)
			}

			uris.Id = path_id
		}

	}

	if cfg.EnableWebFinger {

		settings.WebFingerHostname = cfg.WebFingerHostname

		capabilities.WebFinger = true
		uris.WebFinger = cfg.PathWebFinger
		uris.WebFingerAlt = cfg.PathWebFingerAlt

		if cfg.URIPrefix != "" {

			path_webfinger, err := url.JoinPath(cfg.URIPrefix, cfg.PathWebFinger)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathWebFinger)
			}

			uris.WebFinger = path_webfinger

			alt_paths := make([]string, len(cfg.PathWebFingerAlt))

			for idx, path := range cfg.PathWebFingerAlt {

				path, err := url.JoinPath(cfg.URIPrefix, path)

				if err != nil {
					return nil, fmt.Errorf("Failed to assign prefix to %s, %w", path, err)
				}

				alt_paths[idx] = path
			}

			uris.WebFingerAlt = alt_paths
		}
	}

	if cfg.EnableEditUI {

		capabilities.EditUI = true
		capabilities.CreateFeature = true
		capabilities.DeprecateFeature = true
		capabilities.CessateFeature = true
		capabilities.EditGeometry = true

		uris.CreateFeature = cfg.PathCreateFeature
		uris.EditGeometry = cfg.PathEditGeometry

		if cfg.URIPrefix != "" {

			path_create_feature, err := url.JoinPath(cfg.URIPrefix, cfg.PathCreateFeature)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathCreateFeature)
			}

			path_edit_geometry, err := url.JoinPath(cfg.URIPrefix, cfg.PathEditGeometry)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathEditGeometry)
			}

			uris.CreateFeature = path_create_feature
			uris.EditGeometry = path_edit_geometry
		}
	}

	if cfg.EnableEditAPI {

		capabilities.EditAPI = true
		capabilities.CreateFeatureAPI = true
		capabilities.DeprecateFeatureAPI = true
		capabilities.CessateFeatureAPI = true
		capabilities.EditGeometryAPI = true

		uris.CreateFeatureAPI = cfg.PathCreateFeatureAPI
		uris.DeprecateFeatureAPI = cfg.PathDeprecateFeatureAPI
		uris.CessateFeatureAPI = cfg.PathCessateFeatureAPI
		uris.EditGeometryAPI = cfg.PathEditGeometryAPI

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

			uris.CreateFeatureAPI = path_api_create_feature
			uris.DeprecateFeatureAPI = path_api_deprecate_feature
			uris.CessateFeatureAPI = path_api_cessate_feature
			uris.EditGeometryAPI = path_api_edit_geometry
		}

	}

	if cfg.EnableSearchAPI {

		search_db, err := fulltext.NewFullTextDatabase(ctx, cfg.SearchDatabaseURI)

		if err != nil {
			return nil, fmt.Errorf("Failed to create fulltext database for '%s', %w", cfg.SearchDatabaseURI, err)
		}

		settings.SearchDatabase = search_db

		capabilities.SearchAPI = true
		uris.Search = cfg.PathSearchAPI

		if cfg.URIPrefix != "" {

			path_search_api, err := url.JoinPath(cfg.URIPrefix, cfg.PathSearchAPI)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathSearchAPI, err)
			}

			uris.SearchAPI = path_search_api
		}
	}

	if cfg.EnableSearch {

		capabilities.Search = true
		uris.Search = cfg.PathSearch

		if cfg.URIPrefix != "" {

			path_search, err := url.JoinPath(cfg.URIPrefix, cfg.PathSearch)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign prefix to %s, %w", cfg.PathSearch, err)
			}

			uris.Search = path_search
		}
	}

	settings.URIs = uris
	settings.Capabilities = capabilities

	// Auth hooks

	authenticator, err := auth.NewAuthenticator(ctx, cfg.AuthenticatorURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to create authenticator, %w", err)
	}

	settings.Authenticator = authenticator

	// Map provider

	if cfg.EnableHTML {

		map_provider, err := provider.NewProvider(ctx, cfg.MapProviderURI)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new map provider, %w", err)
		}

		settings.MapProvider = map_provider
	}

	// Custom chrome (this is still in flux)

	custom, err := chrome.NewChrome(ctx, cfg.CustomChromeURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to create custom chrome, %w", err)
	}

	settings.CustomChrome = custom

	// CORS

	if cfg.EnableCORS {

		cors_origins := cfg.CORSOrigins

		if len(cors_origins) == 0 {
			cors_origins = []string{"*"}
		}

		cors_wrapper := cors.New(cors.Options{
			AllowedOrigins:   cors_origins,
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
