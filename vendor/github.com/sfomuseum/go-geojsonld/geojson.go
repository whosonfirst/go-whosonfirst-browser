package geojsonld

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	_ "log"
	"strings"
)

// NS_GEOJSON is the default namespace for GeoJSON-LD
const NS_GEOJSON string = "https://purl.org/geojson/vocab#"

// DefaultGeoJSONLDContext return a dictionary mapping GeoJSON property keys to their GeoJSON-LD @context equivalents.
func DefaultGeoJSONLDContext() map[string]interface{} {

	bbox := map[string]string{
		"@container": "@list",
		"@id":        "geojson:bbox",
	}

	coords := map[string]string{
		"@container": "@list",
		"@id":        "geojson:coordinates",
	}

	features := map[string]string{
		"@container": "@set",
		"@id":        "geojson:features",
	}

	geojson_ctx := map[string]interface{}{
		"geojson":            NS_GEOJSON,
		"Feature":            "geojson:Feature",
		"FeatureCollection":  "geojson:FeatureCollection",
		"GeometryCollection": "geojson:GeometryCollection",
		"LineString":         "geojson:LineString",
		"MultiLineString":    "geojson:MultiLineString",
		"MultiPoint":         "geojson:MultiPoint",
		"MultiPolygon":       "geojson:MultiPolygon",
		"Point":              "geojson:Point",
		"Polygon":            "geojson:Polygon",
		"bbox":               bbox,
		"coordinates":        coords,
		"features":           features,
		"geometry":           "geojson:geometry",
		"id":                 "@id",
		"properties":         "geojson:properties",
		"type":               "@type",
	}

	return geojson_ctx
}

// AsGeoJSONLDWithReader convert GeoJSON Feature data contained in r in to GeoJSON-LD.
func AsGeoJSONLDWithReader(ctx context.Context, r io.Reader) ([]byte, error) {

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to read Feature data, %w", err)
	}

	return AsGeoJSONLD(ctx, body)
}

// AsGeoJSONLDWithReader convert GeoJSON Feature data contained in body in to GeoJSON-LD.
func AsGeoJSONLD(ctx context.Context, body []byte) ([]byte, error) {

	geojson_ctx := DefaultGeoJSONLDContext()

	props_rsp := gjson.GetBytes(body, "properties")

	if !props_rsp.Exists() {
		return nil, fmt.Errorf("Missing properties element")
	}

	extra := make(map[string]string)

	for k, _ := range props_rsp.Map() {

		parts := strings.Split(k, ":")

		var k_fq string

		if len(parts) == 2 {

			ns := parts[0]

			_, exists := extra[ns]

			if exists {
				continue
			}

			extra[ns] = fmt.Sprintf("https://github.com/whosonfirst/whosonfirst-properties/tree/master/properties/%s", ns)

		} else {

			k_fq = fmt.Sprintf("x-urn:geojson:properties#%s", k)
			extra[k] = k_fq
		}
	}

	for k, v := range extra {
		geojson_ctx[k] = v
	}

	body, err := sjson.SetBytes(body, "\\@context", geojson_ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to assign @context, %w", err)
	}

	id_rsp := gjson.GetBytes(body, "id")

	if id_rsp.Exists() {

		body, err = sjson.SetBytes(body, "id", id_rsp.String())

		if err != nil {
			return nil, fmt.Errorf("Failed to assign id, %w", err)
		}
	}

	var i interface{}
	err = json.Unmarshal(body, &i)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal geojson-ld, %w", err)
	}

	enc, err := json.Marshal(i)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal geojson-ld, %w", err)
	}

	return enc, nil
}
