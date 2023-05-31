package browser

import (
	"context"
	"fmt"
	"sync"

	"github.com/whosonfirst/go-whosonfirst-browser/v7/pointinpolygon"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
)

var setupPointInPolygonOnce sync.Once
var setupPointInPolygonError error

func setupPointInPolygon() {

	ctx := context.Background()

	spatial_db, err := database.NewSpatialDatabase(ctx, cfg.SpatialDatabaseURI)

	if err != nil {
		setupPointInPolygonError = fmt.Errorf("Failed to create spatial database, %w", err)
		return
	}

	pt_definition, err := placetypes.NewDefinition(ctx, cfg.PlacetypesDefinitionURI)

	if err != nil {
		setupPointInPolygonError = fmt.Errorf("Failed to create placetypes definition, %w", err)
		return
	}

	pip_options := &pointinpolygon.PointInPolygonServiceOptions{
		SpatialDatabase: spatial_db,
		// FIX ME
		// ParentReader:         ...
		PlacetypesDefinition: pt_definition,
		Logger:               logger,
		SkipPlacetypeFilter:  cfg.PointInPolygonSkipPlacetypeFilter,
	}

	pointinpolygon_service, err = pointinpolygon.NewPointInPolygonServiceWithOptions(ctx, pip_options)

	if err != nil {
		setupPointInPolygonError = fmt.Errorf("Failed to create point in polygon service, %w", err)
		return
	}

	// To do: Set custom PIP options derived from RunOptions
}
