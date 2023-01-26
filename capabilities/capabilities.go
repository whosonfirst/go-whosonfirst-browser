package browser

type Capabilities struct {
	GeoJSON             bool
	GeoJSONLD           bool
	PNG                 bool
	SVG                 bool
	NavPlace            bool
	WebFinger           bool
	Select              bool
	SPR                 bool
	HTML                bool // To do: Rename this; this is the /id/{ID} page
	Search              bool
	SearchAPI           bool
	EditUI              bool
	EditAPI             bool
	CreateFeature       bool
	CreateFeatureAPI    bool
	DeprecateFeature    bool
	DeprecateFeatureAPI bool
	CessateFeature      bool
	CessateFeatureAPI   bool
	EditGeometry        bool
	EditGeometryAPI     bool
}
