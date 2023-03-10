package browser

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/aaronland/go-http-maps/provider"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/sfomuseum/runtimevar"
)

type Config struct {
	AuthenticatorURI                  string   `json:"authenticator_uri"`
	CacheURI                          string   `json:"cache_uri"`
	CORSOrigins                       []string `json:"cors_origins,omitempty"`
	CORSAllowCredentials              bool     `json:"cors_allow_credentials,omitempty"`
	CustomChromeURI                   string   `json:"custom_chrome_uri,omitempty"`
	CustomEditProperties              []string `json:"custom_edit_properties,omitempty"`
	CustomEditValidationWasmDir       string   `json:"custom_edit_validation_wasm_dir"`
	CustomEditValidationWasmPath      string   `json:"custom_edit_validation_wasm_path"`
	DisableGeoJSON                    bool     `json:"disable_geojson,omitempty"`
	DisableGeoJSONLD                  bool     `json:"disable_geojsonld,omitempty"`
	DisableId                         bool     `json:"disable_id,omitempty"`
	DisableIndex                      bool     `json:"disable_index,omitempty"`
	DisableNavPlace                   bool     `json:"disable_navplace,omitempty"`
	DisablePNG                        bool     `json:"disable_png,omitempty"`
	DisableSearch                     bool     `json:"disable_search,omitempty"`
	DisableSelect                     bool     `json:"disable_select"`
	DisableSPR                        bool     `json:"disable_spr"`
	DisableSVG                        bool     `json:"disable_svg,omitempty"`
	DisableWebFinger                  bool     `json:"disable_webfinger"`
	EnableAll                         bool     `json:"enable_all,omitempty"`
	EnableCORS                        bool     `json:"enable_cors,omitempty"`
	EnableCustomEditValidationWasm    bool     `json:"enable_custom_edit_validation_wasm"`
	EnableData                        bool     `json:"enable_data,omitempty"`
	EnableEdit                        bool     `json:"enable_edit,omitempty"`
	EnableEditAPI                     bool     `json:"enable_edit_api,omitempty"`
	EnableEditUI                      bool     `json:"enable_edit_ui,omitempty"`
	EnableGeoJSON                     bool     `json:"enable_geojson,omitempty"`
	EnableGeoJSONLD                   bool     `json:"enable_geojsonld,omitempty"`
	EnableGraphics                    bool     `json:"enable_graphics"`
	EnableHTML                        bool     `json:"enable_html,omitempty"`
	EnableId                          bool     `json:"enable_id,omitempty"`
	EnableIndex                       bool     `json:"enable_index,omitempty"`
	EnableNavPlace                    bool     `json:"enable_navplace,omitempty"`
	EnablePointInPolygon              bool     `json:"enable_point_in_polygon,omitempty"`
	EnablePointInPolygonAPI           bool     `json:"enable_point_in_polygon_api,omitempty"`
	EnablePNG                         bool     `json:"enable_png,omitempty"`
	EnableSearch                      bool     `json:"enable_search,omitempty"`
	EnableSearchAPI                   bool     `json:"enable_search_api,omitempty"`
	EnableSelect                      bool     `json:"enable_select,omitempty"`
	EnableSPR                         bool     `json:"enable_spr,omitempty"`
	EnableSVG                         bool     `json:"enable_svg,omitempty"`
	EnableWebFinger                   bool     `json:"enable_webfinger,omitempty"`
	ExporterURI                       string   `json:"exporter_uri,omitempty"`
	JavaScriptAtEOF                   bool     `json:"javascript_at_eof"`
	GitHubAccessTokenURI              string   `json:"github_accesstoken_uri,omitempty"`
	GitHubReaderAccessTokenURI        string   `json:"github_reader_accesstoken_uri,omitempty"`
	GitHubWriterAccessTokenURI        string   `json:"github_writer_accesstoken_uri,omitempty"`
	MapProviderURI                    string   `json:"map_provider_uri"`
	NavPlaceMaxFeatures               int      `json:"navplace_max_features,omitempty"`
	PathCreateFeature                 string   `json:"path_create_feature,omitempty"`
	PathCreateFeatureAPI              string   `json:"path_create_feature_api,omitempty"`
	PathCessateFeatureAPI             string   `json:"path_cessate_feature_api,omitempty"`
	PathDeprecateFeatureAPI           string   `json:"path_deprecate_feature_api,omitempty"`
	PathEditGeometry                  string   `json:"path_edit_geometry,omitempty"`
	PathEditGeometryAPI               string   `json:"path_edit_geometry_api,omitempty"`
	PathGeoJSON                       string   `json:"path_geojson,omitempty"`
	PathGeoJSONAlt                    []string `json:"path_geojson_alt,omitempty"`
	PathGeoJSONLD                     string   `json:"path_geojsonld,omitempty"`
	PathGeoJSONLDAlt                  []string `json:"path_geojsonld_alt,omitempty"`
	PathIndex                         string   `json:"path_index,omitempty"`
	PathId                            string   `json:"path_id,omitempty"`
	PathNavPlace                      string   `json:"path_navplace,omitempty"`
	PathNavPlaceAlt                   []string `json:"path_navplace_alt,omitempty"`
	PathPing                          string   `json:"path_ping,omitempty"`
	PathPNG                           string   `json:"path_png,omitempty"`
	PathPNGAlt                        []string `json:"path_png_alt,omitempty"`
	PathPointInPolygon                string   `json:"path_point_in_polygon,omitempty"`
	PathPointInPolygonAPI             string   `json:"path_point_in_polygon_api,omitempty"`
	PathSearch                        string   `json:"path_search,omitempty"`
	PathSearchAPI                     string   `json:"path_search_api,omitempty"`
	PathSelect                        string   `json:"path_select,omitempty"`
	PathSelectAlt                     []string `json:"path_select_alt,omitempty"`
	PathSPR                           string   `json:"path_spr,omitempty"`
	PathSPRAlt                        []string `json:"path_spr,omitempty"`
	PathSVG                           string   `json:"path_svg,omitempty"`
	PathSVGAlt                        []string `json:"path_svg_alt,omitempty"`
	PathWebFinger                     string   `json:"path_webfinger,omitempty"`
	PathWebFingerAlt                  []string `json:"path_webfinger_alt,omitempty"`
	PlacetypesDefinitionURI           string   `json:"placetypes_definition_uri,omitempty"`
	PointInPolygonSkipPlacetypeFilter bool     `json:"point_in_polygon_skip_placetype_filter"`
	ReaderURIs                        []string `json:"reader_uris"`
	SearchDatabaseURI                 string   `json:"search_database_uri,omitempty"`
	SelectPattern                     string   `json:"select_pattern,omitempty"`
	ServerURI                         string   `json:"server_uri"`
	SpatialDatabaseURI                string   `json:"spatial_database_uri,omitempty"`
	URIPrefix                         string   `json:"uri_prefix,omitempty"`
	Verbose                           bool     `json:"verbose,omitempty"`
	WebFingerHostname                 string   `json:"webfinger_hostname,omitempty"`
	WriterURIs                        []string `json:"writer_uris,omitempty"`
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

	js_at_eof, err := lookup.BoolVar(fs, provider.JavaScriptAtEOFFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", provider.JavaScriptAtEOFFlag, err)
	}

	cfg := &Config{
		AuthenticatorURI:                  authenticator_uri,
		CacheURI:                          cache_uri,
		CORSAllowCredentials:              cors_allow_credentials,
		CORSOrigins:                       cors_origins,
		CustomChromeURI:                   custom_chrome_uri,
		CustomEditProperties:              custom_edit_properties,
		CustomEditValidationWasmDir:       custom_edit_validation_wasm_dir,
		CustomEditValidationWasmPath:      custom_edit_validation_wasm_path,
		DisableIndex:                      disable_index,
		EnableAll:                         enable_all,
		EnableCORS:                        enable_cors,
		EnableCustomEditValidationWasm:    enable_custom_edit_validation_wasm,
		EnableData:                        enable_data,
		EnableEdit:                        enable_edit,
		EnableEditAPI:                     enable_edit_api,
		EnableEditUI:                      enable_edit_ui,
		EnableGraphics:                    enable_graphics,
		EnableHTML:                        enable_html,
		EnableId:                          enable_id,
		EnableIndex:                       enable_index,
		EnablePNG:                         enable_png,
		EnablePointInPolygon:              enable_point_in_polygon,
		EnablePointInPolygonAPI:           enable_point_in_polygon_api,
		EnableSelect:                      enable_select,
		EnableSearch:                      enable_search,
		EnableSearchAPI:                   enable_search_api,
		EnableSPR:                         enable_spr,
		EnableSVG:                         enable_svg,
		ExporterURI:                       exporter_uri,
		GitHubAccessTokenURI:              github_accesstoken_uri,
		GitHubReaderAccessTokenURI:        github_reader_accesstoken_uri,
		GitHubWriterAccessTokenURI:        github_writer_accesstoken_uri,
		JavaScriptAtEOF:                   js_at_eof,
		MapProviderURI:                    map_provider_uri,
		NavPlaceMaxFeatures:               navplace_max_features,
		PathCreateFeature:                 path_create_feature,
		PathCreateFeatureAPI:              path_api_create_feature,
		PathCessateFeatureAPI:             path_api_cessate_feature,
		PathDeprecateFeatureAPI:           path_api_deprecate_feature,
		PathEditGeometry:                  path_edit_geometry,
		PathEditGeometryAPI:               path_api_edit_geometry,
		PathGeoJSON:                       path_geojson,
		PathGeoJSONAlt:                    path_geojson_alt,
		PathGeoJSONLD:                     path_geojsonld,
		PathGeoJSONLDAlt:                  path_geojsonld_alt,
		PathId:                            path_id,
		PathIndex:                         "/",
		PathNavPlace:                      path_navplace,
		PathNavPlaceAlt:                   path_navplace_alt,
		PathPing:                          path_ping,
		PathPNG:                           path_png,
		PathPNGAlt:                        path_png_alt,
		PathPointInPolygon:                path_point_in_polygon,
		PathPointInPolygonAPI:             path_point_in_polygon_api,
		PathSelect:                        path_select,
		PathSelectAlt:                     path_select_alt,
		PathSPR:                           path_spr,
		PathSPRAlt:                        path_spr_alt,
		PathSVG:                           path_svg,
		PathSVGAlt:                        path_svg_alt,
		PathWebFinger:                     path_webfinger,
		PathWebFingerAlt:                  path_webfinger_alt,
		PlacetypesDefinitionURI:           placetypes_definition_uri,
		PointInPolygonSkipPlacetypeFilter: point_in_polygon_skip_placetype_filter,
		ReaderURIs:                        reader_uris,
		SearchDatabaseURI:                 search_database_uri,
		SelectPattern:                     select_pattern,
		ServerURI:                         server_uri,
		SpatialDatabaseURI:                spatial_database_uri,
		WriterURIs:                        writer_uris,
		URIPrefix:                         static_prefix,
		Verbose:                           verbose,
		WebFingerHostname:                 webfinger_hostname,
	}

	return cfg, nil
}
