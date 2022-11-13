// package www implements HTTP handlers for the whosonfirst-browser web application.
package www

// Endpoints defines a struct containing (relative) URLs for the various whosonfirst-browser web application handlers.
type Endpoints struct {
	Index     string
	Id        string
	Data      string
	Png       string
	Svg       string
	Spr       string
	Search    string
	NavPlace  string
	GeoJSONLD string
}

type ErrorVars struct {
	Error     error
	Endpoints *Endpoints
}

type NotFoundVars struct {
	Endpoints *Endpoints
}
