package browser

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-maps/provider"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

const CacheURIFlag string = "cache-uri"

var cache_uri string

const ConfigURIFlag string = "config-uri"

var config_uri string

const CustomEditPropertiesFlag string = "custom-edit-property"

var custom_edit_properties multi.MultiCSVString

const CustomEditValidationWasmDirFlag string = "custom-edit-validation-wasm-dir"

var custom_edit_validation_wasm_dir string

const CustomEditValidationWasmPathFlag string = "custom-edit-validation-wasm-path"

var custom_edit_validation_wasm_path string

const EnableAllFlag string = "enable-all"

var enable_all bool

const EnableCustomEditValidationWasmFlag string = "enable-custom-edit-validation-wasm"

var enable_custom_edit_validation_wasm bool

const EnableDataFlag string = "enable-data"

var enable_data bool

const EnableGraphicsFlag string = "enable-graphics"

var enable_graphics bool

const EnableHTMLFlag string = "enable-html"

var enable_html bool

const EnableIdFlag string = "enable-id"

var enable_id bool

const EnableIndexFlag string = "enable-index"

var enable_index bool

const EnablePointInPolygonFlag string = "enable-point-in-polygon"

var enable_point_in_polygon bool

const EnablePointInPolygonAPIFlag string = "enable-point-in-polygon-api"

var enable_point_in_polygon_api bool

const DisableGeoJSONFlag string = "disable-geojson"

var disable_geojson bool

const DisableGeoJSONLDFlag string = "disable-geojsonld"

var disable_geojsonld bool

const DisableIdFlag string = "disable-id"

var disable_id bool

const DisableIndexFlag string = "disable-index"

var disable_index bool

const DisableNavPlaceFlag string = "disable-navplace"

var disable_navplace bool

const DisablePNGFlag string = "disable-png"

var disable_png bool

const DisableSearchFlag string = "disable-search"

var disable_search bool

const DisableSelectFlag string = "disable-select"

var disable_select bool

const DisableSPRFlag string = "disable-spr"

var disable_spr bool

const DisableSVGFlag string = "disable-svg"

var disable_svg bool

const DisableWebFingerFlag string = "disable-webfinger"

var disable_webfinger bool

const ExporterURIFlag string = "exporter-uri"

var exporter_uri string

const PathPointInPolygonFlag string = "path-point-in-polygon"

var path_point_in_polygon string

const PathPointInPolygonAPIFlag string = "path-point-in-polygon-api"

var path_point_in_polygon_api string

const PlacetypesFooURIFlag string = "placetypes-foo-uri"

var placetypes_foo_uri string

const ReaderURIFlag string = "reader-uri"

var reader_uris multi.MultiCSVString

const ServerURIFlag string = "server-uri"

var server_uri string

const StaticPrefixFlag string = "static-prefix"

var static_prefix string

const VerboseFlag string = "verbose"

var verbose bool

const WriterURIFlag string = "writer-uri"

var writer_uris multi.MultiCSVString

/* unsorted */

const EnablePNGFlag string = "enable-png"

var enable_png bool

const EnableSVGFlag string = "enable-svg"

var enable_svg bool

const EnableGeoJSONFlag string = "enable-geojson"

var enable_geojson bool

const EnableGeoJSONLDFlag string = "enable-geojson-ld"

var enable_geojsonld bool

const EnableNavPlaceFlag string = "enable-navplace"

var enable_navplace bool

const EnableSPRFlag string = "enable-spr"

var enable_spr bool

const EnableSelect string = "enable-select"

var enable_select bool

const EnableWebFingerFlag string = "enable-webfinger"

var enable_webfinger bool

const SelectPatternFlag string = "select-pattern"

var select_pattern string

const EnableSearchAPIFlag string = "enable-search-api"

var enable_search_api bool

const EnableSearchAPIGeoJSONFlag string = "enable-search-api-geojson"

var enable_search_api_geojson bool

const EnableSearchHTMLFlag string = "enable-search-html"

var enable_search_html bool

const EnableSearchFlag string = "enable-search"

var enable_search bool

const EnableEditFlag string = "enable-edit"

var enable_edit bool

const EnableEditAPIFlag string = "enable-edit-api"

var enable_edit_api bool

const EnableEditUIFlag string = "enable-edit-ui"

var enable_edit_ui bool

const PathAPIDeprecateFeatureFlag string = "path-api-deprecate-feature"

var path_api_deprecate_feature string

const PathAPICessateFeatureFlag string = "path-api-cessate-feature"

var path_api_cessate_feature string

const PathAPIEditGeometryFlag string = "path-api-edit-geometry"

var path_api_edit_geometry string

const PathAPICreateFeatureFlag string = "path-api-create-feature"

var path_api_create_feature string

const GitHubAccessTokenURIFlag string = "github-accesstoken-uri"

// A valid gocloud.dev/runtimevar URI that resolves to a GitHub API access token, required if you are using a githubapi:// reader URI.
var github_accesstoken_uri string

const GitHubReaderAccessTokenURIFlag string = "github-reader-accesstoken-uri"

var github_reader_accesstoken_uri string

const GitHubWriterAccessTokenURIFlag string = "github-writer-accesstoken-uri"

var github_writer_accesstoken_uri string

const PathPingFlag string = "path-ping"

var path_ping string

// The path that PNG requests should be served from.
var path_png string

// Zero or more alternate paths that PNG requests should be served from.
var path_png_alt multi.MultiCSVString

// The path that SVG requests should be served from.
var path_svg string

// Zero or more alternate paths that SVG requests should be served from.
var path_svg_alt multi.MultiCSVString

// The path that GeoJSON requests should be served from.
var path_geojson string

// Zero or more alternate paths that GeoJSON requests should be served from.
var path_geojson_alt multi.MultiCSVString

// The path that GeoJSON-LD requests should be served from.
var path_geojsonld string

// Zero or more alternate paths that GeoJSON-LD requests should be served from.
var path_geojsonld_alt multi.MultiCSVString

// The path that IIIF navPlace requests should be served from.
var path_navplace string

// Zero or more alternate paths that IIIF navPlace requests should be served from.
var path_navplace_alt multi.MultiCSVString

// The path that SPR requests should be served from.
var path_spr string

// Zero or more alternate paths that SPR requests should be served from.
var path_spr_alt multi.MultiCSVString

// The path that 'select' requests should be served from.
var path_select string

// Zero or more alternate paths that 'select' requests should be served from.
var path_select_alt multi.MultiCSVString

// The path that 'webfinger' requests should be served from.
var path_webfinger string

// Zero or more alternate paths that 'webfinger' requests should be served from.
var path_webfinger_alt multi.MultiCSVString

var path_protomaps_tiles string

var path_search_api string
var path_search_html string

const PathEditGeometryFlag string = "path-edit-geometry"

var path_edit_geometry string

const PathCreateFeatureFlag string = "path-create-feature"

var path_create_feature string

var path_id string

var navplace_max_features int

const EnableCORSFlag string = "enable-cors"

var enable_cors bool

const CORSOriginFlag string = "cors-origin"

var cors_origins multi.MultiCSVString

const CORSAllowCredentialsFlag string = "cors-allow-credentials"

var cors_allow_credentials bool

const AuthenticatorURI string = "authenticator-uri"

var authenticator_uri string

const WebFingerHostname string = "webfinger-hostname"

// An optional hostname to use for WebFinger URLs.
var webfinger_hostname string

const SpatialDatabaseURIFlag string = "spatial-database-uri"

var spatial_database_uri string

const CustomChromeURIFlag string = "custom-chrome-uri"

var custom_chrome_uri string

const SearchDatabaseURIFlag string = "search-database-uri"

var search_database_uri string

// DefaultFlagSet returns a `flag.FlagSet` instance with flags and defaults values assigned for use with `app`.
func DefaultFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("browser")

	fs.StringVar(&authenticator_uri, "authenticator-uri", "none://", "A valid sfomuseum/go-http-auth URI.")

	fs.StringVar(&cache_uri, CacheURIFlag, "gocache://", "A valid go-cache Cache URI string.")
	fs.StringVar(&config_uri, ConfigURIFlag, ConfigURIDefault, "...")
	fs.BoolVar(&cors_allow_credentials, CORSAllowCredentialsFlag, CORSAllowCredentialsDefault, "...")
	fs.Var(&cors_origins, CORSOriginFlag, "One or more hosts to restrict CORS support to on the API endpoint. If no origins are defined (and -cors is enabled) then the server will default to all hosts.")
	fs.StringVar(&custom_chrome_uri, CustomChromeURIFlag, "none://", "...")
	fs.Var(&custom_edit_properties, CustomEditPropertiesFlag, "...")
	fs.StringVar(&custom_edit_validation_wasm_dir, CustomEditValidationWasmDirFlag, CustomEditValidationWasmDirDefault, "...")
	fs.StringVar(&custom_edit_validation_wasm_path, CustomEditValidationWasmPathFlag, CustomEditValidationWasmPathDefault, "...")

	fs.BoolVar(&disable_geojson, DisableGeoJSONFlag, DisableGeoJSONDefault, "...")
	fs.BoolVar(&disable_geojsonld, DisableGeoJSONLDFlag, DisableGeoJSONLDDefault, "...")
	fs.BoolVar(&disable_id, DisableIdFlag, DisableIdDefault, "...")
	fs.BoolVar(&disable_index, DisableIndexFlag, DisableIndexDefault, "...")
	fs.BoolVar(&disable_navplace, DisableNavPlaceFlag, DisableNavPlaceDefault, "...")
	fs.BoolVar(&disable_png, DisablePNGFlag, DisablePNGDefault, "...")
	fs.BoolVar(&disable_search, DisableSearchFlag, DisableSearchDefault, "...")
	fs.BoolVar(&disable_spr, DisableSPRFlag, DisableSPRDefault, "...")
	fs.BoolVar(&disable_svg, DisableSVGFlag, DisableSVGDefault, "...")
	fs.BoolVar(&disable_webfinger, DisableWebFingerFlag, DisableWebFingerDefault, "...")

	fs.BoolVar(&enable_all, EnableAllFlag, false, "Enable all the available output handlers EXCEPT the search handlers which need to be explicitly enable using the -enable-search* flags.")
	fs.BoolVar(&enable_cors, EnableCORSFlag, EnableCORSDefault, "A boolean flag to enable CORS headers")
	fs.BoolVar(&enable_custom_edit_validation_wasm, EnableCustomEditValidationWasmFlag, EnableCustomEditValidationWasmDefault, "...")
	fs.BoolVar(&enable_data, EnableDataFlag, false, "Enable the 'geojson' and 'spr' and 'select' output handlers.")
	fs.BoolVar(&enable_edit, EnableEditFlag, EnableEditDefault, "Enable the API endpoints")
	fs.BoolVar(&enable_edit_api, EnableEditAPIFlag, EnableEditAPIDefault, "...")
	fs.BoolVar(&enable_edit_ui, EnableEditUIFlag, EnableEditUIDefault, "...")
	fs.BoolVar(&enable_geojson, "enable-geojson", true, "Enable the 'geojson' output handler.")
	fs.BoolVar(&enable_geojsonld, "enable-geojson-ld", true, "Enable the 'geojson-ld' output handler.")
	fs.BoolVar(&enable_graphics, EnableGraphicsFlag, false, "Enable the 'png' and 'svg' output handlers.")
	fs.BoolVar(&enable_html, EnableHTMLFlag, EnableHTMLDefault, "Enable the 'geojson' and 'spr' and 'select' output handlers.")
	fs.BoolVar(&enable_id, EnableIdFlag, EnableIdDefault, "Enable the 'html' (or human-friendly) output handlers.")
	fs.BoolVar(&enable_index, EnableIndexFlag, EnableIndexDefault, "Enable the 'index' (or human-friendly) index handler.")
	fs.BoolVar(&enable_navplace, "enable-navplace", true, "Enable the IIIF 'navPlace' output handler.")
	fs.BoolVar(&enable_png, EnablePNGFlag, false, "Enable the 'png' output handler.")
	fs.BoolVar(&enable_point_in_polygon, EnablePointInPolygonFlag, EnablePointInPolygonDefault, "...")
	fs.BoolVar(&enable_point_in_polygon_api, EnablePointInPolygonAPIFlag, EnablePointInPolygonAPIDefault, "...")
	fs.BoolVar(&enable_search, "enable-search", false, "Enable both the API and human-friendly search handlers.")
	fs.BoolVar(&enable_search_api, "enable-search-api", false, "Enable the (API) search handlers.")
	fs.BoolVar(&enable_select, "enable-select", false, "Enable the 'select' output handler.")
	fs.BoolVar(&enable_spr, "enable-spr", true, "Enable the 'spr' (or \"standard places response\") output handler.")
	fs.BoolVar(&enable_svg, EnableSVGFlag, false, "Enable the 'svg' output handler.")
	fs.BoolVar(&enable_webfinger, EnableWebFingerFlag, false, "Enable the 'webfinger' output handler.")
	fs.StringVar(&exporter_uri, ExporterURIFlag, "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export/v2 URI.")

	fs.StringVar(&github_accesstoken_uri, GitHubAccessTokenURIFlag, GitHubAccessTokenURIDefault, "A valid gocloud.dev/runtimevar URI that resolves to a GitHub API access token.")
	fs.StringVar(&github_reader_accesstoken_uri, GitHubReaderAccessTokenURIFlag, GitHubReaderAccessTokenURIDefault, "...")
	fs.StringVar(&github_writer_accesstoken_uri, GitHubWriterAccessTokenURIFlag, GitHubWriterAccessTokenURIDefault, "...")

	fs.IntVar(&navplace_max_features, "navplace-max-features", 3, "The maximum number of features to allow in a /navplace/{ID} URI string.")

	fs.StringVar(&path_api_deprecate_feature, PathAPIDeprecateFeatureFlag, PathAPIDeprecateFeatureDefault, "...")
	fs.StringVar(&path_api_cessate_feature, PathAPICessateFeatureFlag, PathAPICessateFeatureDefault, "...")
	fs.StringVar(&path_api_edit_geometry, PathAPIEditGeometryFlag, "/api/geometry/", "...")
	fs.StringVar(&path_api_create_feature, PathAPICreateFeatureFlag, PathAPICreateFeatureDefault, "...")
	fs.StringVar(&path_create_feature, PathCreateFeatureFlag, PathCreateFeatureDefault, "...")
	fs.StringVar(&path_edit_geometry, PathEditGeometryFlag, PathEditGeometryDefault, "...")
	fs.StringVar(&path_geojson, "path-geojson", "/geojson/", "The path that GeoJSON requests should be served from.")
	fs.Var(&path_geojson_alt, "path-geojson-alt", "Zero or more alternate paths that GeoJSON requests should be served from.")
	fs.StringVar(&path_geojsonld, "path-geojson-ld", "/geojson-ld/", "The path that GeoJSON-LD requests should be served from.")
	fs.Var(&path_geojsonld_alt, "path-geojson-ld-alt", "Zero or more alternate paths that GeoJSON-LD requests should be served from.")
	fs.StringVar(&path_id, "path-id", "/id/", "The URL that Who's On First documents should be served from.")
	fs.StringVar(&path_navplace, "path-navplace", "/navplace/", "The path that IIIF navPlace requests should be served from.")
	fs.Var(&path_navplace_alt, "path-navplace-alt", "Zero or more alternate paths that IIIF navPlace requests should be served from.")
	fs.StringVar(&path_ping, PathPingFlag, PathPingDefault, "...")
	fs.StringVar(&path_png, "path-png", "/png/", "The path that PNG requests should be served from.")
	fs.Var(&path_png_alt, "path-png-alt", "Zero or more alternate paths that PNG requests should be served from.")
	fs.StringVar(&path_point_in_polygon, PathPointInPolygonFlag, PathPointInPolygonDefault, "...")
	fs.StringVar(&path_point_in_polygon_api, PathPointInPolygonAPIFlag, PathPointInPolygonAPIDefault, "...")
	fs.StringVar(&path_protomaps_tiles, "path-protomaps-tiles", "/tiles/", "The root path from which Protomaps tiles will be served.")
	fs.StringVar(&path_search_api, "path-search-api", "/search/spr/", "The path that API 'search' requests should be served from.")
	fs.StringVar(&path_search_html, "path-search-html", "/search/", "The path that API 'search' requests should be served from.")
	fs.StringVar(&path_select, "path-select", "/select/", "The path that 'select' requests should be served from.")
	fs.Var(&path_select_alt, "path-select-alt", "Zero or more alternate paths that 'select' requests should be served from.")
	fs.StringVar(&path_spr, "path-spr", "/spr/", "The path that SPR requests should be served from.")
	fs.Var(&path_spr_alt, "path-spr-alt", "Zero or more alternate paths that SPR requests should be served from.")
	fs.StringVar(&path_svg, "path-svg", "/svg/", "The path that SVG requests should be served from.")
	fs.Var(&path_svg_alt, "path-svg-alt", "Zero or more alternate paths that SVG requests should be served from.")
	fs.StringVar(&path_webfinger, "path-webfinger", "/.well-known/webfinger/", "The path that 'webfinger' requests should be served from.")
	fs.Var(&path_webfinger_alt, "path-webfinger-alt", "Zero or more alternate paths that 'webfinger' requests should be served from.")
	fs.StringVar(&placetypes_foo_uri, PlacetypesFooURIFlag, PlacetypesFooURIDefault, "A valid go-whosonfirst-placetypes.Foo URI.")

	fs.Var(&reader_uris, ReaderURIFlag, "One or more valid go-reader Reader URI strings.")

	fs.StringVar(&search_database_uri, SearchDatabaseURIFlag, SearchDatabaseURIDefault, "A valid whosonfirst/go-whosonfist-search/fulltext URI.")
	fs.StringVar(&select_pattern, "select-pattern", "properties(?:.[a-zA-Z0-9-_]+){1,}", "A valid regular expression for sanitizing select parameters.")
	fs.StringVar(&server_uri, ServerURIFlag, "http://localhost:8080", "A valid aaronland/go-http-server URI.")
	fs.StringVar(&spatial_database_uri, SpatialDatabaseURIFlag, "", "...")
	fs.StringVar(&static_prefix, StaticPrefixFlag, "", "Prepend this prefix to URLs for static assets.")

	fs.BoolVar(&verbose, VerboseFlag, VerboseDefault, "...")

	fs.StringVar(&webfinger_hostname, WebFingerHostname, "", "An optional hostname to use for WebFinger URLs.")
	fs.Var(&writer_uris, "writer-uri", "One or more valid go-writer Writer URI strings.")

	err := provider.AppendProviderFlags(fs)

	if err != nil {
		return nil, fmt.Errorf("Failed to append map provider flags, %v", err)
	}

	return fs, nil
}
