package browser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/api"
)

func apiSearchHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWhosOnFirstReaderOnce.Do(setupWhosOnFirstReader)

	if setupWhosOnFirstReaderError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstReader)
	}

	setupSearchOnce.Do(setupSearch)

	if setupSearchError != nil {
		return nil, fmt.Errorf("Failed to configure search setup, %w", setupSearchError)
	}

	search_opts := api.SearchAPIHandlerOptions{
		Database:      search_database,
		EnableGeoJSON: true,
		GeoJSONReader: wof_reader,
	}

	search_api_handler, err := api.SearchAPIHandler(search_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create search handler, %w", err)
	}

	if cors_wrapper != nil {
		search_api_handler = cors_wrapper.Handler(search_api_handler)
	}

	return search_api_handler, nil
}

func apiDeprecateHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWhosOnFirstReaderOnce.Do(setupWhosOnFirstReader)

	if setupWhosOnFirstReaderError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstReader)
	}

	setupWhosOnFirstWriterOnce.Do(setupWhosOnFirstWriter)

	if setupWhosOnFirstWriterError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstWriter)
	}

	setupAuthenticatorOnce.Do(setupAuthenticator)

	if setupAuthenticatorError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupAuthenticator)
	}

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

	setupWhosOnFirstReaderOnce.Do(setupWhosOnFirstReader)

	if setupWhosOnFirstReaderError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstReader)
	}

	setupWhosOnFirstWriterOnce.Do(setupWhosOnFirstWriter)

	if setupWhosOnFirstWriterError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstWriter)
	}

	setupAuthenticatorOnce.Do(setupAuthenticator)

	if setupAuthenticatorError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupAuthenticator)
	}

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

	setupWhosOnFirstReaderOnce.Do(setupWhosOnFirstReader)

	if setupWhosOnFirstReaderError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstReader)
	}

	setupWhosOnFirstWriterOnce.Do(setupWhosOnFirstWriter)

	if setupWhosOnFirstWriterError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstWriter)
	}

	setupPointInPolygonOnce.Do(setupPointInPolygon)

	if setupPointInPolygonError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupPointInPolygonError)
	}

	setupAuthenticatorOnce.Do(setupAuthenticator)

	if setupPointInPolygonError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupAuthenticator)
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

func apiCreateFeatureHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWhosOnFirstReaderOnce.Do(setupWhosOnFirstReader)

	if setupWhosOnFirstReaderError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstReader)
	}

	setupWhosOnFirstWriterOnce.Do(setupWhosOnFirstWriter)

	if setupWhosOnFirstWriterError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstWriter)
	}

	setupPointInPolygonOnce.Do(setupPointInPolygon)

	if setupPointInPolygonError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupPointInPolygonError)
	}

	setupAuthenticatorOnce.Do(setupAuthenticator)

	if setupAuthenticatorError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupAuthenticator)
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
