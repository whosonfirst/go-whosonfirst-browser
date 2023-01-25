package browser

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/aaronland/go-http-maps/provider"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/runtimevar"
)

type Config struct {
	// Placeholder for a JSON-encoded config representation that can be used in place of command-line flags

	CacheURI         string   `json:"cache_uri"`
	ReaderURIs       []string `json:"reader_uris"`
	WriterURIs       []string `json:"writer_uris,omitempty"`
	AuthenticatorURI string   `json:"authenticator_uri"`
	ServerURI        string   `json:"server_uri"`
	MapProviderURI   string   `json:"map_provider_uri"`

	URIPrefix       string `json:"uri_prefix,omitempty"`
	ExporterURI     string `json:"exporter_uri,omitempty"`
	EnableAll       bool   `json:"enable_all,omitempty"`
	EnableGraphics  bool   `json:"enable_graphics"`
	EnableData      bool   `json:"enable_data,omitempty"`
	EnableGeoJSON   bool   `json:"enable_geojson,omitempty"`
	EnableGeoJSONLD bool   `json:"enable_geojsonld,omitempty"`
	EnablePNG       bool   `json:"enable_png,omitempty"`
	EnableSVG       bool   `json:"enable_svg,omitempty"`
	EnableSelect    bool   `json:"enable_select,omitempty"`
	EnableSPR       bool   `json:"enable_spr,omitempty"`
	EnableHTML      bool   `json:"enable_html,omitempty"`
	EnableNavPlace  bool   `json:"enable_navplace,omitempty"`
	EnableWebFinger bool   `json:"enable_webfinger,omitempty"`
	EnableSearch    bool   `json:"enable_search,omitempty"`
	EnableIndex     bool   `json:"enable_index,omitempty"`
	EnableCORS      bool   `json:"enable_cors,omitempty"`
	EnableEdit      bool   `json:"enable_edit,omitempty"`
	EnableEditAPI   bool   `json:"enable_edit_api,omitempty"`
	EnableEditUI    bool   `json:"enable_edit_ui,omitempty"`

	PathPing       string   `json:"path_ping,omitempty"`
	PathPNG        string   `json:"path_png,omitempty"`
	PathPNGAlt     []string `json:"path_png_alt,omitempty"`
	PathSVG        string   `json:"path_svg,omitempty"`
	PathSVGAlt     []string `json:"path_svg_alt,omitempty"`
	PathGeoJSON    string   `json:"path_geojson,omitempty"`
	PathGeoJSONAlt []string `json:"path_geojson_alt,omitempty"`

	PathGeoJSONLD    string   `json:"path_geojsonld,omitempty"`
	PathGeoJSONLDAlt []string `json:"path_geojsonld_alt,omitempty"`

	PathNavPlace    string   `json:"path_navplace,omitempty"`
	PathNavPlaceAlt []string `json:"path_navplace_alt,omitempty"`

	PathSelect    string   `json:"path_select,omitempty"`
	PathSelectAlt []string `json:"path_select_alt,omitempty"`

	PathWebFinger    string   `json:"path_webfinger,omitempty"`
	PathWebFingerAlt []string `json:"path_webfinger_alt,omitempty"`

	PathId string `json:"path_id,omitempty"`

	PathEditGeometry    string `json:"path_edit_geometry,omitempty"`
	PathEditGeometryAPI string `json:"path_edit_geometry_api,omitempty"`

	PathCreateFeature    string `json:"path_create_feature,omitempty"`
	PathCreateFeatureAPI string `json:"path_create_feature_api,omitempty"`

	PathDeprecateFeatureAPI string `json:"path_deprecate_feature_api,omitempty"`
	PathCessateFeatureAPI   string `json:"path_cessate_feature_api,omitempty"`

	SelectPattern string `json:"select_pattern,omitempty"`

	NavPlaceMaxFeatures int `json:"navplace_max_features,omitempty"`

	CORSOrigins          []string `json:"cors_origins,omitempty"`
	CORSAllowCredentials bool     `json:"cors_allow_credentials,omitempty"`

	GitHubAccessTokenURI       string `json:"github_accesstoken_uri,omitempty"`
	GitHubReaderAccessTokenURI string `json:"github_reader_accesstoken_uri,omitempty"`
	GitHubWriterAccessTokenURI string `json:"github_writer_accesstoken_uri,omitempty"`

	WebFingerHostname string `json:"webfinger_hostname,omitempty"`

	SpatialDatabaseURI string `json:"spatial_database_uri,omitempty"`

	CustomChromeURI string `json:"custom_chrome_uri,omitempty"`

	Verbose bool `json:"verbose,omitempty"`
}

func ConfigFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*Config, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "BROWSER")

	if err != nil {
		return nil, fmt.Errorf("Failed to set flags from environment variables, %w", err)
	}

	if config_uri != "" {

		var cfg *Config

		v, err := runtimevar.StringVar(ctx, config_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive config from URI, %w", err)
		}

		r := strings.NewReader(v)

		dec := json.NewDecoder(r)
		err = dec.Decode(&cfg)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode config, %w", err)
		}

		return cfg, nil
	}

	map_provider_uri, err := provider.ProviderURIFromFlagSet(fs)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive provider URI from flagset, %v", err)
	}

	cfg := &Config{
		CacheURI:         cache_uri,
		ReaderURIs:       reader_uris,
		WriterURIs:       writer_uris,
		AuthenticatorURI: authenticator_uri,
		ServerURI:        server_uri,
		ExporterURI:      exporter_uri,
		MapProviderURI:   map_provider_uri,

		URIPrefix: static_prefix,

		EnableAll:      enable_all,
		EnableGraphics: enable_graphics,
		EnableData:     enable_data,
		EnablePNG:      enable_png,
		EnableSVG:      enable_svg,
		EnableSelect:   enable_select,
		EnableSPR:      enable_spr,
		EnableHTML:     enable_html,
		EnableIndex:    enable_index,
		EnableCORS:     enable_cors,
		EnableEdit:     enable_edit,
		EnableEditAPI:  enable_edit_api,
		EnableEditUI:   enable_edit_ui,

		PathPing:                path_ping,
		PathPNG:                 path_png,
		PathPNGAlt:              path_png_alt,
		PathSVG:                 path_svg,
		PathSVGAlt:              path_svg_alt,
		PathGeoJSON:             path_geojson,
		PathGeoJSONAlt:          path_geojson_alt,
		PathGeoJSONLD:           path_geojsonld,
		PathGeoJSONLDAlt:        path_geojsonld_alt,
		PathNavPlace:            path_navplace,
		PathNavPlaceAlt:         path_navplace_alt,
		PathWebFinger:           path_webfinger,
		PathWebFingerAlt:        path_webfinger_alt,
		PathId:                  path_id,
		PathEditGeometry:        path_edit_geometry,
		PathEditGeometryAPI:     path_api_edit_geometry,
		PathCreateFeature:       path_create_feature,
		PathCreateFeatureAPI:    path_api_create_feature,
		PathDeprecateFeatureAPI: path_api_deprecate_feature,
		PathCessateFeatureAPI:   path_api_cessate_feature,

		NavPlaceMaxFeatures: navplace_max_features,

		SelectPattern: select_pattern,

		CORSOrigins:          cors_origins,
		CORSAllowCredentials: cors_allow_credentials,

		GitHubAccessTokenURI:       github_accesstoken_uri,
		GitHubReaderAccessTokenURI: github_reader_accesstoken_uri,
		GitHubWriterAccessTokenURI: github_writer_accesstoken_uri,

		WebFingerHostname: webfinger_hostname,

		SpatialDatabaseURI: spatial_database_uri,

		CustomChromeURI: custom_chrome_uri,

		Verbose: verbose,
	}

	return cfg, nil
}
