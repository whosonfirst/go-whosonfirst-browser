package www

type Paths struct {
	PNG                 string `json:"png,omitempty"`
	SVG                 string `json:"svg,omitempty"`
	GeoJSON             string `json:"geojson,omitempty"`
	GeoJSONLD           string `json:"geojsonld,omitempty"`
	NavPlace            string `json:"navplace,omitempty"`
	SPR                 string `json:"spr,omitempty"`
	Select              string `json:"select,omitempty"`
	URIPrefix           string `json:"uriprefix,omitempty"`
	Id                  string `json:"id,omitempty"`
	Index               string `json:"index,omitempty"`
	Search              string `json:"search,omitempty"`
	Ping                string `json:"ping,omitempty"`
	CreateFeature       string `json:"create_feature,omitempty"`
	CreateFeatureAPI    string `json:"create_feature_api,omitempty"`
	DeprecateFeatureAPI string `json:"deprecate_feature_api,omitempty"`
	CessateFeatureAPI   string `json:"cessate_feature_api,omitempty"`
	EditGeometry        string `json:"edit_geometry,omitempty"`
	EditGeometryAPI     string `json:"edit_geometry_api,omitempty"`
}
