package browser

import (
	"context"
	"flag"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

const ServerURIFlag string = "server-uri"

var server_uri string

const StaticPrefixFlag string = "static-prefix"

var static_prefix string

const ReaderURIFlag string = "reader-uri"

var reader_uris multi.MultiCSVString

const CacheURIFlag string = "cache-uri"

var cache_uri string

const ExporterURIFlag string = "exporter-uri"

var exporter_uri string

const WriterURIFlag string = "writer-uri"

var writer_uris multi.MultiCSVString

var map_provider string

var nextzen_api_key string
var nextzen_style_url string
var nextzen_tile_url string

var proxy_tiles bool
var proxy_tiles_url string
var proxy_tiles_cache string
var proxy_tiles_timeout int

var tilepack_db string
var tilepack_uri string

var protomaps_bucket_uri string
var protomaps_cache_size int
var protomaps_tiles_database string

const EnableAllFlag string = "enable-all"

var enable_all bool

const EnableGraphicsFlag string = "enable-graphics"

var enable_graphics bool

const EnableDataFlag string = "enable-data"

var enable_data bool

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

const EnableHTMLFlag string = "enable-html"

var enable_html bool

const EnableIndexFlag string = "enable-index"

var enable_index bool

const EnableSearchAPIFlag string = "enable-search-api"

var enable_search_api bool

const EnableSearchAPIGeoJSONFlag string = "enable-search-api-geojson"

var enable_search_api_geojson bool

const EnableSearchHTMLFlag string = "enable-search-html"

var enable_search_html bool

const EnableSearchFlag string = "enable-search"

var enable_search bool

const EnableAPIFlag string = "enable-api"

var enable_api bool

var path_api_deprecate string
var path_api_cessate string

var search_database_uri string

// A valid gocloud.dev/runtimevar URI that resolves to a GitHub API access token, required if you are using a githubapi:// reader URI.
var github_accesstoken_uri string

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

var path_id string

var navplace_max_features int

const EnableCORSFlag string = "enable-cors"

var enable_cors bool

const CORSOriginFlag string = "cors-origin"

var cors_origins multi.MultiCSVString

const AuthenticatorURI string = "authenticator-uri"

var authenticator_uri string

const WebFingerHostname string = "webfinger-hostname"

// An optional hostname to use for WebFinger URLs.
var webfinger_hostname string


// DefaultFlagSet returns a `flag.FlagSet` instance with flags and defaults values assigned for use with `app`.
func DefaultFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("browser")

	fs.StringVar(&server_uri, ServerURIFlag, "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.StringVar(&static_prefix, StaticPrefixFlag, "", "Prepend this prefix to URLs for static assets.")

	fs.Var(&reader_uris, ReaderURIFlag, "One or more valid go-reader Reader URI strings.")
	fs.StringVar(&cache_uri, CacheURIFlag, "gocache://", "A valid go-cache Cache URI string.")

	fs.StringVar(&exporter_uri, ExporterURIFlag, "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export/v2 URI.")

	// START OF replace/reconcile with aaronland/go-http-maps

	fs.StringVar(&map_provider, "map-provider", "nextzen", "Valid options are: nextzen, protomaps")

	fs.StringVar(&nextzen_api_key, "nextzen-api-key", "", "A valid Nextzen API key (https://developers.nextzen.org/).")
	fs.StringVar(&nextzen_style_url, "nextzen-style-url", "/tangram/refill-style.zip", "A valid Tangram scene file URL.")
	fs.StringVar(&nextzen_tile_url, "nextzen-tile-url", tangramjs.NEXTZEN_MVT_ENDPOINT, "A valid Nextzen MVT tile URL.")

	fs.StringVar(&tilepack_db, "nextzen-tilepack-database", "", "The path to a valid MBTiles database (tilepack) containing Nextzen MVT tiles.")
	fs.StringVar(&tilepack_uri, "nextzen-tilepack-uri", "/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt", "The relative URI to serve Nextzen MVT tiles from a MBTiles database (tilepack).")

	fs.BoolVar(&proxy_tiles, "proxy-tiles", false, "Proxy (and cache) Nextzen tiles.")
	fs.StringVar(&proxy_tiles_url, "proxy-tiles-url", "/tiles/", "The URL (a relative path) for proxied tiles.")
	fs.StringVar(&proxy_tiles_cache, "proxy-tiles-cache", "gocache://", "A valid `whosonfirst/go-cache` URI.")
	fs.IntVar(&proxy_tiles_timeout, "proxy-tiles-timeout", 30, "The maximum number of seconds to allow for fetching a tile from the proxy.")

	fs.StringVar(&protomaps_bucket_uri, "protomaps-bucket-uri", "", "A valid gocloud.dev/blob.Bucket URI containing Protomaps tile databases.")
	fs.IntVar(&protomaps_cache_size, "protomaps-cache-size", 64, "The size in MB of the Protomaps tile cache.")
	fs.StringVar(&protomaps_tiles_database, "protomaps-tiles-database", "", "The name of the Protomaps tiles database to use.")

	// END OF replace/reconcile with aaronland/go-http-maps

	fs.BoolVar(&enable_all, EnableAllFlag, false, "Enable all the available output handlers EXCEPT the search handlers which need to be explicitly enable using the -enable-search* flags.")
	fs.BoolVar(&enable_graphics, EnableGraphicsFlag, false, "Enable the 'png' and 'svg' output handlers.")
	fs.BoolVar(&enable_data, EnableDataFlag, false, "Enable the 'geojson' and 'spr' and 'select' output handlers.")

	fs.BoolVar(&enable_png, EnablePNGFlag, false, "Enable the 'png' output handler.")
	fs.BoolVar(&enable_svg, EnableSVGFlag, false, "Enable the 'svg' output handler.")

	fs.BoolVar(&enable_webfinger, EnableWebFingerFlag, false, "Enable the 'webfinger' output handler.")

	fs.BoolVar(&enable_geojson, "enable-geojson", true, "Enable the 'geojson' output handler.")
	fs.BoolVar(&enable_geojsonld, "enable-geojson-ld", true, "Enable the 'geojson-ld' output handler.")
	fs.BoolVar(&enable_navplace, "enable-navplace", true, "Enable the IIIF 'navPlace' output handler.")
	fs.BoolVar(&enable_spr, "enable-spr", true, "Enable the 'spr' (or \"standard places response\") output handler.")
	fs.BoolVar(&enable_select, "enable-select", false, "Enable the 'select' output handler.")
	fs.StringVar(&select_pattern, "select-pattern", "properties(?:.[a-zA-Z0-9-_]+){1,}", "A valid regular expression for sanitizing select parameters.")

	fs.BoolVar(&enable_html, "enable-html", true, "Enable the 'html' (or human-friendly) output handlers.")
	fs.BoolVar(&enable_index, "enable-index", true, "Enable the 'index' (or human-friendly) index handler.")

	fs.BoolVar(&enable_search_api, "enable-search-api", false, "Enable the (API) search handlers.")
	fs.BoolVar(&enable_search_api_geojson, "enable-search-api-geojson", false, "Enable the (API) search handlers to return results as GeoJSON.")
	fs.BoolVar(&enable_search_html, "enable-search-html", false, "Enable the (human-friendly) search handlers.")
	fs.BoolVar(&enable_search, "enable-search", false, "Enable both the API and human-friendly search handlers.")
	fs.StringVar(&search_database_uri, "search-database-uri", "", "A valid whosonfirst/go-whosonfist-search/fulltext URI.")

	fs.StringVar(&path_png, "path-png", "/png/", "The path that PNG requests should be served from.")
	fs.Var(&path_png_alt, "path-png-alt", "Zero or more alternate paths that PNG requests should be served from.")

	fs.StringVar(&path_svg, "path-svg", "/svg/", "The path that SVG requests should be served from.")
	fs.Var(&path_svg_alt, "path-svg-alt", "Zero or more alternate paths that SVG requests should be served from.")

	fs.StringVar(&path_geojson, "path-geojson", "/geojson/", "The path that GeoJSON requests should be served from.")
	fs.Var(&path_geojson_alt, "path-geojson-alt", "Zero or more alternate paths that GeoJSON requests should be served from.")

	fs.StringVar(&path_geojsonld, "path-geojson-ld", "/geojson-ld/", "The path that GeoJSON-LD requests should be served from.")
	fs.Var(&path_geojsonld_alt, "path-geojson-ld-alt", "Zero or more alternate paths that GeoJSON-LD requests should be served from.")

	fs.StringVar(&path_navplace, "path-navplace", "/navplace/", "The path that IIIF navPlace requests should be served from.")
	fs.Var(&path_navplace_alt, "path-navplace-alt", "Zero or more alternate paths that IIIF navPlace requests should be served from.")

	fs.StringVar(&path_spr, "path-spr", "/spr/", "The path that SPR requests should be served from.")
	fs.Var(&path_spr_alt, "path-spr-alt", "Zero or more alternate paths that SPR requests should be served from.")

	fs.StringVar(&path_select, "path-select", "/select/", "The path that 'select' requests should be served from.")
	fs.Var(&path_select_alt, "path-select-alt", "Zero or more alternate paths that 'select' requests should be served from.")

	fs.StringVar(&path_webfinger, "path-webfinger", "/.well-known/webfinger/", "The path that 'webfinger' requests should be served from.")
	fs.Var(&path_webfinger_alt, "path-webfinger-alt", "Zero or more alternate paths that 'webfinger' requests should be served from.")

	fs.StringVar(&path_protomaps_tiles, "path-protomaps-tiles", "/tiles/", "The root path from which Protomaps tiles will be served.")

	fs.StringVar(&path_search_api, "path-search-api", "/search/spr/", "The path that API 'search' requests should be served from.")
	fs.StringVar(&path_search_html, "path-search-html", "/search/", "The path that API 'search' requests should be served from.")

	fs.StringVar(&path_id, "path-id", "/id/", "The URL that Who's On First documents should be served from.")

	fs.IntVar(&navplace_max_features, "navplace-max-features", 3, "The maximum number of features to allow in a /navplace/{ID} URI string.")

	fs.BoolVar(&enable_cors, "enable-cors", true, "A boolean flag to enable CORS headers")
	fs.Var(&cors_origins, "cors-origin", "One or more hosts to restrict CORS support to on the API endpoint. If no origins are defined (and -cors is enabled) then the server will default to all hosts.")

	fs.StringVar(&authenticator_uri, "authenticator-uri", "null://", "A valid sfomuseum/go-http-auth URI.")

	fs.StringVar(&github_accesstoken_uri, "github-accesstoken-uri", "", "A valid gocloud.dev/runtimevar URI that resolves to a GitHub API access token, required if you are using a githubapi:// reader URI.")

	fs.StringVar(&webfinger_hostname, WebFingerHostname, "", "An optional hostname to use for WebFinger URLs.")
	
	/*
		fs.BoolVar(&enable_api, "enable-api", false, "Enable the API endpoints")
		fs.StringVar(&path_api_deprecate, "path-api-deprecate", "/api/deprecate/", "...")
		fs.StringVar(&path_api_cessate, "path-api-cessate", "/api/cessate/", "...")
		fs.Var(&writer_uris, "writer-uri", "One or more valid go-writer Writer URI strings.")
	*/

	return fs, nil
}
