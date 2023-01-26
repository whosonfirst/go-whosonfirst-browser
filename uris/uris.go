package browser

type URIs struct {
	CessateFeatureAPI   string   `json:"cessate_feature_api,omitempty"`
	CreateFeature       string   `json:"create_feature,omitempty"`
	CreateFeatureAPI    string   `json:"create_feature_api,omitempty"`
	DeprecateFeatureAPI string   `json:"deprecate_feature_api,omitempty"`
	EditGeometry        string   `json:"edit_geometry,omitempty"`
	EditGeometryAPI     string   `json:"edit_geometry_api,omitempty"`
	GeoJSON             string   `json:"geojson,omitempty"`
	GeoJSONAlt          []string `json:"geojson_alt,omitempty"`
	GeoJSONLD           string   `json:"geojsonld,omitempty"`
	GeoJSONLDAlt        []string `json:"geojsonld_alt,omitempty"`
	Id                  string   `json:"id,omitempty"`
	Index               string   `json:"index,omitempty"`
	NavPlace            string   `json:"navplace,omitempty"`
	NavPlaceAlt         []string `json:"navplace_alt,omitempty"`
	PNG                 string   `json:"png,omitempty"`
	PNGAlt              []string `json:"png_alt,omitempty"`
	Ping                string   `json:"ping,omitempty"`
	Search              string   `json:"search,omitempty"`
	SearchAPI           string   `json:"search_api,omitempty"`
	Select              string   `json:"select,omitempty"`
	SelectAlt           []string `json:"select_alt,omitempty"`
	SVG                 string   `json:"svg,omitempty"`
	SVGAlt              []string `json:"svg_alt,omitempty"`
	SPR                 string   `json:"spr,omitempty"`
	SPRAlt              []string `json:"spr_alt,omitempty"`
	WebFinger           string   `json:"webfinger,omitempty"`
	WebFingerAlt        []string `json:"webfinger_alt,omitempty"`
	URIPrefix           string   `json:"uriprefix,omitempty"`
}
