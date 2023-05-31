package capabilities

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
	PointInPolygon      bool
	PointInPolygonAPI   bool
	RollupAssets        bool
}

func (c *Capabilities) HasHTMLCapabilities() bool {

	if c.Index || c.Id || c.Search || c.EditUI {
		return true
	}

	return false
}
