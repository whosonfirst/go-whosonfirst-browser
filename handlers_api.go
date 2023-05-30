package browser

import (
	"context"
	"fmt"
	"net/http"
)

func apiSearchHandlerFunc(ctx context.Context) (http.Handler, error) {

	search_opts := api.SearchAPIHandlerOptions{
		Database:      search_database,
		EnableGeoJSON: true,
		GeoJSONReader: wof_reader,
	}

	search_api_handler, err := api.SearchAPIHandler(search_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create search handler, %w", err)
	}

	if settings.CORSWrapper != nil {
		search_api_handler = settings.CORSWrapper.Handler(search_api_handler)
	}

	return search_api_handler, nil
}

func apiDeprecateHandlerFunc(ctx context.Context) (http.Handler, error) {

	// Writers are created at runtime using the http/api/publish.go#publishFeature
	// method which in turn calls writer/writer.go#NewWriter

	// Deprecate a record

	deprecate_opts := &api.DeprecateFeatureHandlerOptions{
		Reader:        wof_reader,
		Cache:         wof_cache,
		Logger:        logger,
		Authenticator: authenticator,
		Exporter:      wof_exporter,
		WriterURIs:    wof_writer_uris,
	}

	deprecate_handler, err := api.DeprecateFeatureHandler(deprecate_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create deprecate feature handler, %w", err)
	}

	deprecate_handler = authenticator.WrapHandler(deprecate_handler)

	return deprecate_handler, nil
}

func apiCessateHandlerFunc(ctx context.Context) (http.Handler, error) {

	cessate_opts := &api.CessateFeatureHandlerOptions{
		Reader:        wof_reader,
		Cache:         wof_cache,
		Logger:        logger,
		Authenticator: authenticator,
		Exporter:      wof_exporter,
		WriterURIs:    wof_writer_uris,
	}

	cessate_handler, err := api.CessateFeatureHandler(cessate_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create cessate feature handler, %w", err)
	}

	cessate_handler = authenticator.WrapHandler(cessate_handler)

	return cessate_handler, nil
}

func apiEditGeometryHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupPointInPolygonOnce.Do(setupPointInPolygon)

	if setupPointInPolygonError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupPointInPolygonError)
	}

	geom_opts := &api.UpdateGeometryHandlerOptions{
		Reader:                wof_reader,
		Cache:                 wof_cache,
		Logger:                logger,
		Authenticator:         authenticator,
		Exporter:              wof_exporter,
		WriterURIs:            wof_writer_uris,
		PointInPolygonService: pointinpolygon_service,
	}

	geom_handler, err := api.UpdateGeometryHandler(geom_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create update geometry handler, %w", err)
	}

	geom_handler = authenticator.WrapHandler(geom_handler)
	return geom_handler, nil
}

func apiCreateFeatureHandlerOption(ctx context.Context) (http.Handler, error) {

	setupPointInPolygonOnce.Do(setupPointInPolygon)

	if setupPointInPolygonError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupPointInPolygonError)
	}

	create_opts := &api.CreateFeatureHandlerOptions{
		Reader:                wof_reader,
		Cache:                 wof_cache,
		Logger:                logger,
		Authenticator:         authenticator,
		Exporter:              wof_exporter,
		WriterURIs:            wof_writer_uris,
		PointInPolygonService: pointinpolygon_service,
		// CustomProperties:      settings.CustomEditProperties,
		// CustomValidationFunc:  settings.CustomEditValidationFunc,
	}

	create_handler, err := api.CreateFeatureHandler(create_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create create feature handler, %w", err)
	}

	create_handler = authenticator.WrapHandler(create_handler)
	return create_handler, nil
}
