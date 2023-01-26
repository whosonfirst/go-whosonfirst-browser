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
	Index               bool
	Id                  bool
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
