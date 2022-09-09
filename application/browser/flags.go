package browser

import (
	"context"
	"flag"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var server_uri string

var static_prefix string

var reader_uris multi.MultiString
var cache_uri string

var nextzen_api_key string
var nextzen_style_url string
var nextzen_tile_url string

var proxy_tiles bool
var proxy_tiles_url string
var proxy_tiles_cache string
var proxy_tiles_timeout int

var tilepack_db string
var tilepack_uri string

var enable_all bool
var enable_graphics bool
var enable_data bool

var enable_png bool
var enable_svg bool

var enable_geojson bool
var enable_geojsonld bool
var enable_navplace bool
var enable_spr bool
var enable_select bool

var select_pattern string

var enable_html bool
var enable_index bool

var enable_search_api bool
var enable_search_api_geojson bool

var enable_search_html bool

var enable_search bool

var search_database_uri string

var path_png string
var path_svg string
var path_geojson string
var path_geojsonld string
var path_navplace string
var path_spr string
var path_select string

var path_search_api string
var path_search_html string

var path_id string

var navplace_max_features int

var enable_cors bool

var cors_origins multi.MultiCSVString

// DefaultFlagSet returns a `flag.FlagSet` instance with flags and defaults values assigned for use with `app`.
func DefaultFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("browser")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.StringVar(&static_prefix, "static-prefix", "", "Prepend this prefix to URLs for static assets.")

	fs.Var(&reader_uris, "reader-uri", "One or more valid go-reader Reader URI strings.")
	fs.StringVar(&cache_uri, "cache-uri", "gocache://", "A valid go-cache Cache URI string.")

	fs.StringVar(&nextzen_api_key, "nextzen-api-key", "", "A valid Nextzen API key (https://developers.nextzen.org/).")
	fs.StringVar(&nextzen_style_url, "nextzen-style-url", "/tangram/refill-style.zip", "A valid Tangram scene file URL.")
	fs.StringVar(&nextzen_tile_url, "nextzen-tile-url", tangramjs.NEXTZEN_MVT_ENDPOINT, "A valid Nextzen MVT tile URL.")

	fs.StringVar(&tilepack_db, "nextzen-tilepack-database", "", "The path to a valid MBTiles database (tilepack) containing Nextzen MVT tiles.")
	fs.StringVar(&tilepack_uri, "nextzen-tilepack-uri", "/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt", "The relative URI to serve Nextzen MVT tiles from a MBTiles database (tilepack).")

	fs.BoolVar(&proxy_tiles, "proxy-tiles", false, "Proxy (and cache) Nextzen tiles.")
	fs.StringVar(&proxy_tiles_url, "proxy-tiles-url", "/tiles/", "The URL (a relative path) for proxied tiles.")
	fs.StringVar(&proxy_tiles_cache, "proxy-tiles-cache", "gocache://", "A valid tile proxy DSN string.")
	fs.IntVar(&proxy_tiles_timeout, "proxy-tiles-timeout", 30, "The maximum number of seconds to allow for fetching a tile from the proxy.")

	fs.BoolVar(&enable_all, "enable-all", false, "Enable all the available output handlers EXCEPT the search handlers which need to be explicitly enable using the -enable-search* flags.")
	fs.BoolVar(&enable_graphics, "enable-graphics", false, "Enable the 'png' and 'svg' output handlers.")
	fs.BoolVar(&enable_data, "enable-data", false, "Enable the 'geojson' and 'spr' and 'select' output handlers.")

	fs.BoolVar(&enable_png, "enable-png", false, "Enable the 'png' output handler.")
	fs.BoolVar(&enable_svg, "enable-svg", false, "Enable the 'svg' output handler.")

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
	fs.StringVar(&path_svg, "path-svg", "/svg/", "The path that SVG requests should be served from.")
	fs.StringVar(&path_geojson, "path-geojson", "/geojson/", "The path that GeoJSON requests should be served from.")
	fs.StringVar(&path_geojsonld, "path-geojson-ld", "/geojson-ld/", "The path that GeoJSON-LD requests should be served from.")
	fs.StringVar(&path_navplace, "path-navplace", "/navplace/", "The path that IIIF navPlace requests should be served from.")
	fs.StringVar(&path_spr, "path-spr", "/spr/", "The path that SPR requests should be served from.")
	fs.StringVar(&path_select, "path-select", "/select/", "The path that 'select' requests should be served from.")

	fs.StringVar(&path_search_api, "path-search-api", "/search/spr/", "The path that API 'search' requests should be served from.")
	fs.StringVar(&path_search_html, "path-search-html", "/search/", "The path that API 'search' requests should be served from.")

	fs.StringVar(&path_id, "path-id", "/id/", "The URL that Who's On First documents should be served from.")

	fs.IntVar(&navplace_max_features, "navplace-max-features", 3, "The maximum number of features to allow in a /navplace/{ID} URI string.")

	fs.BoolVar(&enable_cors, "enable-cors", true, "A boolean flag to enable CORS headers")
	fs.Var(&cors_origins, "cors-origin", "One or more hosts to restrict CORS support to on the API endpoint. If no origins are defined (and -cors is enabled) then the server will default to all hosts.")

	return fs, nil
}
