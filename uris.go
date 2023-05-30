package browser

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type URIs struct {
	CessateFeatureAPI   string            `json:"cessate_feature_api,omitempty"`
	CreateFeature       string            `json:"create_feature,omitempty"`
	CreateFeatureAPI    string            `json:"create_feature_api,omitempty"`
	Custom              map[string]string `json:"custom,omitempty"`
	DeprecateFeatureAPI string            `json:"deprecate_feature_api,omitempty"`
	EditGeometry        string            `json:"edit_geometry,omitempty"`
	EditGeometryAPI     string            `json:"edit_geometry_api,omitempty"`
	GeoJSON             string            `json:"geojson,omitempty"`
	GeoJSONAlt          []string          `json:"geojson_alt,omitempty"`
	GeoJSONLD           string            `json:"geojsonld,omitempty"`
	GeoJSONLDAlt        []string          `json:"geojsonld_alt,omitempty"`
	Id                  string            `json:"id,omitempty"`
	Index               string            `json:"index,omitempty"`
	NavPlace            string            `json:"navplace,omitempty"`
	NavPlaceAlt         []string          `json:"navplace_alt,omitempty"`
	PNG                 string            `json:"png,omitempty"`
	PNGAlt              []string          `json:"png_alt,omitempty"`
	Ping                string            `json:"ping,omitempty"`
	PointInPolygon      string            `json:"point_in_polygon,omitempty"`
	PointInPolygonAPI   string            `json:"point_in_polygon_api,omitempty"`
	Search              string            `json:"search,omitempty"`
	SearchAPI           string            `json:"search_api,omitempty"`
	Select              string            `json:"select,omitempty"`
	SelectAlt           []string          `json:"select_alt,omitempty"`
	SVG                 string            `json:"svg,omitempty"`
	SVGAlt              []string          `json:"svg_alt,omitempty"`
	SPR                 string            `json:"spr,omitempty"`
	SPRAlt              []string          `json:"spr_alt,omitempty"`
	WebFinger           string            `json:"webfinger,omitempty"`
	WebFingerAlt        []string          `json:"webfinger_alt,omitempty"`
	URIPrefix           string            `json:"uriprefix,omitempty"`
}

func (u *URIs) AddCustomURI(key string, value string) error {

	if u.Custom == nil {
		u.Custom = make(map[string]string)
	}

	u.Custom[key] = value
	return nil
}

func (u *URIs) ApplyPrefix(prefix string) error {

	val := reflect.ValueOf(*u)

	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i)
		v := field.String()

		if v == "" {
			continue
		}

		if strings.HasPrefix(v, prefix) {
			continue
		}

		new_v, err := url.JoinPath(prefix, v)

		if err != nil {
			return fmt.Errorf("Failed to assign prefix to %s, %w", v, err)
		}

		reflect.ValueOf(u).Elem().Field(i).SetString(new_v)
	}

	return nil
}
