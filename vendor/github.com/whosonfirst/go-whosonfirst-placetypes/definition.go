package placetypes

// Type WOFPlacetypeDefinition is a struct that maps to the JSON definition files
// for individual placetypes in the `whosonfirst-placetypes` repo. Note that as of
// this writing it does not account for BCP-47 name: properties.
type WOFPlacetypeDefinition struct {
	Id           int64             `json:"wof:id"`
	Name         string            `json:"wof:name"`
	Role         string            `json:"wof:role"`
	Parent       []string          `json:"wof:parent"`
	Concordances map[string]string `json:"wof:concordances"`
}
