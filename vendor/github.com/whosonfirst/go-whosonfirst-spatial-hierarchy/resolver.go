package hierarchy

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	aa_log "github.com/aaronland/go-log/v2"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-sfomuseum-mapshaper"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	hierarchy_filter "github.com/whosonfirst/go-whosonfirst-spatial-hierarchy/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type PointInPolygonHierarchyResolverOptions struct {
	// Database is the `database.SpatialDatabase` instance used to perform point-in-polygon requests.
	Database database.SpatialDatabase
	// Mapshaper is an optional `mapshaper.Client` instance used to derive centroids used in point-in-polygon requests.
	Mapshaper *mapshaper.Client
	// PlacetypesDefinition is an optional `go-whosonfirst-placetypes.Definition` instance used to resolve custom or bespoke placetypes.
	PlacetypesDefinition placetypes.Definition
	Logger               *log.Logger
}

// PointInPolygonHierarchyResolver provides methods for constructing a hierarchy of ancestors
// for a given point, following rules established by the Who's On First project.
type PointInPolygonHierarchyResolver struct {
	// Database is the `database.SpatialDatabase` instance used to perform point-in-polygon requests.
	Database database.SpatialDatabase
	// Mapshaper is an optional `mapshaper.Client` instance used to derive centroids used in point-in-polygon requests.
	Mapshaper *mapshaper.Client
	// PlacetypesDefinition is an optional `go-whosonfirst-placetypes.Definition` instance used to resolve custom or bespoke placetypes.
	PlacetypesDefinition placetypes.Definition
	// Logger is a `log.Logger` instance for logging feedback. The default logger sends all messages to `io.Discard`.
	Logger *log.Logger
	// reader is the `reader.Reader` instance used to retrieve ancestor records. By default it is the same as `Database` but can be assigned
	// explicitly using the `SetReader` method.
	reader reader.Reader
}

// PointInPolygonHierarchyResolverUpdateCallback is a function definition for a custom callback to convert 'spr' in to a dictionary of properties
// containining hierarchy information. Records in 'spr' are expected to be able to be read from 'r'.
type PointInPolygonHierarchyResolverUpdateCallback func(context.Context, reader.Reader, spr.StandardPlacesResult) (map[string]interface{}, error)

// DefaultPointInPolygonHierarchyResolverUpdateCallback returns a `PointInPolygonHierarchyResolverUpdateCallback` function that will return a dictionary
// containing the following properties: wof:parent_id, wof:country, wof:hierarchy
func DefaultPointInPolygonHierarchyResolverUpdateCallback() PointInPolygonHierarchyResolverUpdateCallback {

	fn := func(ctx context.Context, r reader.Reader, parent_spr spr.StandardPlacesResult) (map[string]interface{}, error) {

		to_update := make(map[string]interface{})

		if parent_spr == nil {

			to_update = map[string]interface{}{
				"properties.wof:parent_id": -1,
			}

		} else {

			parent_id, err := strconv.ParseInt(parent_spr.Id(), 10, 64)

			if err != nil {
				return nil, fmt.Errorf("Failed to parse ID (%s), %w", parent_spr.Id(), err)
			}

			parent_f, err := wof_reader.LoadBytes(ctx, r, parent_id)

			if err != nil {
				return nil, fmt.Errorf("Failed to load body for %d, %w", parent_id, err)
			}

			parent_hierarchy := properties.Hierarchies(parent_f)
			parent_country := properties.Country(parent_f)

			to_update = map[string]interface{}{
				"properties.wof:parent_id": parent_id,
				"properties.wof:country":   parent_country,
				"properties.wof:hierarchy": parent_hierarchy,
			}
		}

		return to_update, nil
	}

	return fn
}

// NewPointInPolygonHierarchyResolver returns a `PointInPolygonHierarchyResolver` instance for 'spatial_db' and 'ms_client'.
// The former is used to perform point in polygon operations and the latter is used to determine a "reverse geocoding" centroid
// to use for point-in-polygon operations.
func NewPointInPolygonHierarchyResolver(ctx context.Context, opts *PointInPolygonHierarchyResolverOptions) (*PointInPolygonHierarchyResolver, error) {

	var logger *log.Logger
	var pt_def placetypes.Definition

	if opts.Logger != nil {
		logger = opts.Logger
	} else {
		logger = log.New(io.Discard, "", 0)
	}

	if opts.PlacetypesDefinition == nil {

		def, err := placetypes.NewDefinition(ctx, "whosonfirst://")

		if err != nil {
			return nil, fmt.Errorf("Failed to create whosonfirst:// placetypes definition, %w", err)
		}

		pt_def = def

	} else {

		pt_def = opts.PlacetypesDefinition
	}

	t := &PointInPolygonHierarchyResolver{
		Database:             opts.Database,
		Mapshaper:            opts.Mapshaper,
		PlacetypesDefinition: pt_def,
		Logger:               logger,
		reader:               opts.Database,
	}

	return t, nil
}

// SetReader assigns 'r' as the internal `reader.Reader` instance used to retrieve ancestor records when resolving a hierarchy.
func (t *PointInPolygonHierarchyResolver) SetReader(r reader.Reader) {
	t.reader = r
}

// PointInPolygonAndUpdate will ...
func (t *PointInPolygonHierarchyResolver) PointInPolygonAndUpdate(ctx context.Context, inputs *filter.SPRInputs, results_cb hierarchy_filter.FilterSPRResultsFunc, update_cb PointInPolygonHierarchyResolverUpdateCallback, body []byte) (bool, []byte, error) {

	possible, err := t.PointInPolygon(ctx, inputs, body)

	if err != nil {
		return false, nil, fmt.Errorf("Failed to perform point in polygon operation, %w", err)
	}

	parent_spr, err := results_cb(ctx, t.reader, body, possible)

	if err != nil {
		return false, nil, fmt.Errorf("Results callback failed, %w", err)
	}

	to_assign, err := update_cb(ctx, t.reader, parent_spr)

	if err != nil {
		return false, nil, fmt.Errorf("Update callback failed, %w", err)
	}

	if to_assign == nil {
		return false, body, nil
	}

	has_changed, body, err := export.AssignPropertiesIfChanged(ctx, body, to_assign)

	if err != nil {
		return false, nil, fmt.Errorf("Failed to assign properties, %w", err)
	}

	return has_changed, body, nil
}

// PointInPolygon will perform a point-in-polygon (reverse geocoding) operation for 'body' using zero or more 'inputs' as query filters.
// This is known to not work as expected if the `wof:placetype` property is "common". There needs to be a way to a) retrieve placetypes
// using a custom WOFPlacetypeSpecification (go-whosonfirst-placetypes v0.6.0+) and b) specify an alternate property to retrieve placetypes
// from if `wof:placetype=custom`.
func (t *PointInPolygonHierarchyResolver) PointInPolygon(ctx context.Context, inputs *filter.SPRInputs, body []byte) ([]spr.StandardPlacesResult, error) {

	pt_def := t.PlacetypesDefinition
	pt_spec := pt_def.Specification()
	pt_prop := pt_def.Property()
	pt_uri := pt_def.URI()

	pt_path := fmt.Sprintf("properties.%s", pt_prop)

	pt_rsp := gjson.GetBytes(body, pt_path)

	if !pt_rsp.Exists() {
		return nil, fmt.Errorf("Missing %s property", pt_path)
	}

	pt_str := pt_rsp.String()

	pt, err := pt_spec.GetPlacetypeByName(pt_str)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new placetype for '%s', %v", pt_str, err)
	}

	roles := placetypes.AllRoles()

	ancestors := pt_spec.AncestorsForRoles(pt, roles)

	centroid, err := t.PointInPolygonCentroid(ctx, body)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive centroid, %w", err)
	}

	lon := centroid.X()
	lat := centroid.Y()

	// Start PIP-ing the list of ancestors - stop at the first match

	possible := make([]spr.StandardPlacesResult, 0)

	for _, a := range ancestors {

		coord, err := geo.NewCoordinate(lon, lat)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new coordinate, %w", err)
		}

		pt_name := fmt.Sprintf("%s#%s", a.Name, pt_uri)

		inputs.Placetypes = []string{pt_name}

		spr_filter, err := filter.NewSPRFilterFromInputs(inputs)

		if err != nil {
			return nil, fmt.Errorf("Failed to create SPR filter from input, %v", err)
		}

		aa_log.Debug(t.Logger, "Perform point in polygon at %f, %f for %s\n", lat, lon, pt_name)

		rsp, err := t.Database.PointInPolygon(ctx, coord, spr_filter)

		if err != nil {
			return nil, fmt.Errorf("Failed to point in polygon for %v, %v", coord, err)
		}

		// This should never happen...

		if rsp == nil {
			return nil, fmt.Errorf("Failed to point in polygon for %v, null response", coord)
		}

		results := rsp.Results()
		count := len(results)

		aa_log.Debug(t.Logger, "Point in polygon results at %f, %f for %s: %d\n", lat, lon, pt_name, count)

		if count == 0 {
			continue
		}

		possible = results
		break
	}

	return possible, nil
}

// PointInPolygonCentroid derives an *orb.Point (or "centroid") to use for point-in-polygon operations.
func (t *PointInPolygonHierarchyResolver) PointInPolygonCentroid(ctx context.Context, body []byte) (*orb.Point, error) {

	f, err := geojson.UnmarshalFeature(body)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal featur body, %w", err)
	}

	// First see whether there are exsiting reverse-geocoding properties
	// that we can use

	props := f.Properties

	to_try := []string{
		"reversegeo",
		"lbl",
		"mps",
	}

	for _, prefix := range to_try {

		key_lat := fmt.Sprintf("%s:latitude", prefix)
		key_lon := fmt.Sprintf("%s:longitude", prefix)

		lat, ok_lat := props[key_lat]
		lon, ok_lon := props[key_lon]

		if !ok_lat || ok_lon {
			continue
		}

		pt := &orb.Point{
			lat.(float64),
			lon.(float64),
		}

		return pt, nil
	}

	// Next see what kind of feature we are working with

	var candidate *geojson.Feature

	geojson_type := f.Geometry.GeoJSONType()

	switch geojson_type {
	case "Point":
		candidate = f
	case "MultiPoint":

		// not at all clear this is the best way to deal with things
		// (20210204/thisisaaronland)

		bound := f.Geometry.Bound()
		pt := bound.Center()

		candidate = geojson.NewFeature(pt)

	case "Polygon", "MultiPolygon":

		if t.Mapshaper == nil {

			bound := f.Geometry.Bound()
			pt := bound.Center()

			candidate = geojson.NewFeature(pt)

		} else {

			// this is not great but it's also not hard and making
			// the "perfect" mapshaper interface is yak-shaving right
			// now (20210204/thisisaaronland)

			fc := geojson.NewFeatureCollection()
			fc.Append(f)

			fc, err := t.Mapshaper.AppendCentroids(ctx, fc)

			if err != nil {
				return nil, fmt.Errorf("Failed to append centroids, %v", err)
			}

			f = fc.Features[0]

			candidate = geojson.NewFeature(f.Geometry)

			lat, lat_ok := f.Properties["mps:latitude"]
			lon, lon_ok := f.Properties["mps:longitude"]

			if lat_ok && lon_ok {

				pt := orb.Point{
					lat.(float64),
					lon.(float64),
				}

				candidate = geojson.NewFeature(pt)
			}
		}

	default:
		return nil, fmt.Errorf("Unsupported type '%v'", t)
	}

	pt := candidate.Geometry.(orb.Point)
	return &pt, nil
}
