package pointinpolygon

// Is this general enough to put in a common WOF or other package? Not sure yet...

import (
	_ "github.com/whosonfirst/go-whosonfirst-spatial-pmtiles"
)

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial-hierarchy"
	hierarchy_filter "github.com/whosonfirst/go-whosonfirst-spatial-hierarchy/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	spatial_filter "github.com/whosonfirst/go-whosonfirst-spatial/filter"
)

type PointInPolygonService struct {
	resolver             *hierarchy.PointInPolygonHierarchyResolver
	parent_reader        reader.Reader
	ResultsCallback      hierarchy_filter.FilterSPRResultsFunc
	UpdateCallback       hierarchy.PointInPolygonHierarchyResolverUpdateCallback
	PlacetypesDefinition placetypes.Definition
}

type PointInPolygonServiceOptions struct {
	SpatialDatabase      database.SpatialDatabase
	ParentReader         reader.Reader
	PlacetypesDefinition placetypes.Definition
	// Mapshaper ...

}

func NewPointInPolygonService(ctx context.Context, spatial_database_uri string, parent_reader_uri string) (*PointInPolygonService, error) {

	parent_reader, err := reader.NewReader(ctx, parent_reader_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create parent reader, %w", err)
	}

	spatial_db, err := database.NewSpatialDatabase(ctx, spatial_database_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create spatial database, %w", err)
	}

	return NewPointInPolygonServiceWithDatabaseAndReader(ctx, spatial_db, parent_reader)
}

func NewPointInPolygonServiceWithDatabaseAndReader(ctx context.Context, spatial_db database.SpatialDatabase, parent_reader reader.Reader) (*PointInPolygonService, error) {

	service_opts := &PointInPolygonServiceOptions{
		SpatialDatabase: spatial_db,
		ParentReader:    parent_reader,
	}

	return NewPointInPolygonServiceWithOptions(ctx, service_opts)
}

func NewPointInPolygonServiceWithOptions(ctx context.Context, opts *PointInPolygonServiceOptions) (*PointInPolygonService, error) {

	resolver_opts := &hierarchy.PointInPolygonHierarchyResolverOptions{
		Database:             opts.SpatialDatabase,
		PlacetypesDefinition: opts.PlacetypesDefinition,
		Mapshaper:            nil,
	}

	resolver, err := hierarchy.NewPointInPolygonHierarchyResolver(ctx, resolver_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create hierarchy resolver, %w", err)
	}

	// This is particularly important if we are using a PMTiles/Protomaps spatial database
	// because the spatial hierarchy resolver may need to retrieve properties that have not
	// been encoded in the PMTiles database.

	resolver.SetReader(opts.ParentReader)

	results_cb := hierarchy_filter.FirstButForgivingSPRResultsFunc
	update_cb := hierarchy.DefaultPointInPolygonHierarchyResolverUpdateCallback()

	s := &PointInPolygonService{
		resolver:        resolver,
		parent_reader:   opts.ParentReader,
		ResultsCallback: results_cb,
		UpdateCallback:  update_cb,
	}

	return s, nil
}

func (s *PointInPolygonService) Update(ctx context.Context, body []byte) (bool, []byte, error) {

	inputs := &spatial_filter.SPRInputs{
		IsCurrent: []int64{1},
	}

	return s.UpdateWithInputs(ctx, body, inputs)
}

func (s *PointInPolygonService) UpdateWithInputs(ctx context.Context, body []byte, inputs *spatial_filter.SPRInputs) (bool, []byte, error) {

	has_changes, new_body, err := s.resolver.PointInPolygonAndUpdate(ctx, inputs, s.ResultsCallback, s.UpdateCallback, body)

	if err != nil {
		return false, nil, fmt.Errorf("Failed to update feature, %w", err)
	}

	if !has_changes {
		return false, body, nil
	}

	parent_rsp := gjson.GetBytes(new_body, "properties.wof:parent_id")

	if parent_rsp.Exists() && parent_rsp.Int() != int64(-1) {

		parent_id := parent_rsp.Int()
		parent_f, err := wof_reader.LoadBytes(ctx, s.parent_reader, parent_id)

		if err == nil {

			to_copy := []string{
				"properties.sfomuseum:post_security",
				"properties.sfo:level",
			}

			updates := make(map[string]interface{})

			for _, path := range to_copy {

				rsp := gjson.GetBytes(parent_f, path)

				if !rsp.Exists() {
					continue
				}

				updates[path] = rsp.Value()
			}

			parent_changes := false

			parent_changes, new_body, err = export.AssignPropertiesIfChanged(ctx, new_body, updates)

			if err != nil {
				return false, nil, fmt.Errorf("Failed to assign properties from parent, %w", err)
			}

			if !has_changes && parent_changes {
				has_changes = true
			}
		}
	}

	return has_changes, new_body, nil
}
