package browser

import (
	"context"
	"fmt"
	"net/http"
)

func sprHandlerFunc(ctx context.Context) (http.Handler, error) {

	spr_opts := &www.SPRHandlerOptions{
		Reader: settings.Reader,
		Logger: logger,
	}

	spr_handler, err := www.SPRHandler(spr_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create SPR handler, %w", err)
	}

	if cors_handler != nil {
		spr_handler = cors_handler.Handler(spr_handler)
	}

	return spr_handler, nil
}

func geojsonHandlerFunc(ctx context.Context) (http.Handler, error) {

	geojson_opts := &www.GeoJSONHandlerOptions{
		Reader: settings.Reader,
		Logger: logger,
	}

	geojson_handler, err := www.GeoJSONHandler(geojson_opts)

	if err != nil {
		return fmt.Errorf("Failed to create GeoJSON handler, %w", err)
	}

	if cors_wrapper != nil {
		geojson_handler = cors_wrapper.Handler(geojson_handler)
	}

	return geojson_handler, nil
}

func geojsonLDHandlerFunc(ctx context.Context) (http.Handler, error) {

	geojsonld_opts := &www.GeoJSONLDHandlerOptions{
		Reader: settings.Reader,
		Logger: logger,
	}

	geojsonld_handler, err := www.GeoJSONLDHandler(geojsonld_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create GeoJSON-LD handler, %w", err)
	}

	if cors_wrapper != nil {
		geojsonld_handler = cors_wrapper.Handler(geojsonld_handler)
	}

	return geojsonld_handler, nil
}

func navPlaceHandlerFunc(ctx context.Context) (http.Handler, error) {

	navplace_opts := &www.NavPlaceHandlerOptions{
		Reader:      settings.Reader,
		MaxFeatures: settings.NavPlaceMaxFeatures,
		Logger:      logger,
	}

	navplace_handler, err := www.NavPlaceHandler(navplace_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create IIIF navPlace handler, %w", err)
	}

	if cors_wrapper != nil {
		navplace_handler = cors_wrapper.Handler(navplace_handler)
	}

	return navplace_handler, nil
}

func selectHandlerFunction(ctx context.Context) (http.Handler, error) {

	select_opts := &www.SelectHandlerOptions{
		Pattern: settings.SelectPattern,
		Reader:  settings.Reader,
		Logger:  logger,
	}

	select_handler, err := www.SelectHandler(select_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create Select handler, %w", err)
	}

	if cors_wrapper != nil {
		select_handler = cors_wrapper.Handler(select_handler)
	}

	return select_handler, nil
}

func webFingerHandlerFunc(ctx context.Context) (http.Handler, error) {

	webfinger_opts := &www.WebfingerHandlerOptions{
		Reader:       settings.Reader,
		Logger:       logger,
		URIs:         uris_table,
		Capabilities: settings.Capabilities,
		Hostname:     settings.WebFingerHostname,
	}

	webfinger_handler, err := www.WebfingerHandler(webfinger_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create WebFinger handler, %w", err)
	}

	if cors_wrapper != nil {
		webfinger_handler = cors_wrapper.Handler(webfinger_handler)
	}

	return webfinger_handler, nil
}
