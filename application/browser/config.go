package browser

type Config struct {
	// Placeholder for a JSON-encoded config representation that can be used in place of command-line flags

	CacheURI         string
	ReaderURIs       []string
	WriterURIs       []string
	AuthenticatorURI string
	ServerURI        string
	URIPrefix        string
	ExporterURI      string
	EnableAll        bool
	EnableGraphics   bool
	EnableData       bool
	EnablePNG        bool
	EnableSVG        bool
	EnableSelect     bool
	EnableSPR        bool
	SelectPattern    string
	EnableHTML       bool
	EnableIndex      bool

	PathPing       string
	PathPNG        string
	PathPNGAlt     []string
	PathSVG        string
	PathSVGAlt     []string
	PathGeoJSON    string
	PathGeoJSONAlt []string

	PathGeoJSONLD    string
	PathGeoJSONLDAlt []string

	PathNavPlace    string
	PathNavPlaceAlt []string

	PathSelect    string
	PathSelectAlt []string

	PathWebFinger    string
	PathWebFingerAlt []string

	PathId string

	PathEditGeometry    string
	PathEditGeometryAPI string

	PathCreateFeature    string
	PathCreateFeatureAPI string

	PathDeprecateFeatureAPI string
	PathCessateFeatureAPI   string `json:"path_cessate_feature_api"`

	NavPlaceMaxFeatures int `json:"navplace_max_features"`

	EnableCORs           bool     `json:"enable_cors"`
	CORSOrigins          []string `json:"cors_origins"`
	CORSAllowCredentials bool     `json:"cors_allow_credentials"`

	GitHubAccessTokenURI       string `json:"github_accesstoken_uri"`
	GitHubReaderAccessTokenURI string `json:"github_reader_accesstoken_uri"`
	GitHubWriterAccessTokenURI string `json:"github_writer_accesstoken_uri"`

	WebFingerHostname string `json:"webfinger_hostname"`

	EnableEdit    bool `json:"enable_edit"`
	EnableEditAPI bool `json:"enable_edit_api"`
	EnableEditUI  bool `json:"enable_edit_ui"`

	SpatialDatabaseURI string `json:"spatial_database_uri"`

	CustomChromeURI string `json:"custom_chrome_uri"`

	Verbose bool `json:"verbose"`
}
