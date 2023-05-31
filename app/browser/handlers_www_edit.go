package browser

import (
	"context"
	"fmt"
	"net/http"

	// wasm_exec "github.com/sfomuseum/go-http-wasm/v2"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-maps"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
)

func wwwEditGeometryHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		return nil, fmt.Errorf("Failed to configure www setup, %w", setupWWWError)
	}

	setupWhosOnFirstReaderOnce.Do(setupWhosOnFirstReader)

	if setupWhosOnFirstReaderError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstReader)
	}

	setupAuthenticatorOnce.Do(setupAuthenticator)

	if setupAuthenticatorError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupAuthenticator)
	}

	geom_t := html_t.Lookup("geometry")

	if geom_t == nil {
		return nil, fmt.Errorf("Failed to load 'geometry' template")
	}

	geom_opts := &www.EditGeometryHandlerOptions{
		Authenticator: authenticator,
		MapProvider:   map_provider.Scheme(),
		URIs:          uris_table,
		Capabilities:  capabilities,
		Template:      geom_t,
		Logger:        logger,
		Reader:        wof_reader,
	}

	geom_handler, err := www.EditGeometryHandler(geom_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create edit geometry handler, %w", err)
	}

	geom_handler = maps.AppendResourcesHandlerWithProvider(geom_handler, map_provider, maps_opts)
	geom_handler = bootstrap.AppendResourcesHandler(geom_handler, bootstrap_opts)
	geom_handler = www.AppendResourcesHandler(geom_handler, www_opts.WithGeometryHandlerResources())

	// FIX ME
	// geom_handler = settings.CustomChrome.WrapHandler(geom_handler, "whosonfirst.browser.geometry")

	geom_handler = authenticator.WrapHandler(geom_handler)
	return geom_handler, nil
}

func wwwCreateGeometryHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		return nil, fmt.Errorf("Failed to configure www setup, %w", setupWWWError)
	}

	setupWhosOnFirstReaderOnce.Do(setupWhosOnFirstReader)

	if setupWhosOnFirstReaderError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupWhosOnFirstReader)
	}

	setupAuthenticatorOnce.Do(setupAuthenticator)

	if setupAuthenticatorError != nil {
		return nil, fmt.Errorf("Failed to configure PIP set up, %w", setupAuthenticator)
	}

	create_t := html_t.Lookup("create")

	if create_t == nil {
		return nil, fmt.Errorf("Failed to load 'create' template")
	}

	create_opts := &www.CreateFeatureHandlerOptions{
		Authenticator: authenticator,
		MapProvider:   map_provider.Scheme(),
		URIs:          uris_table,
		Capabilities:  capabilities,
		Template:      create_t,
		Logger:        logger,
		Reader:        wof_reader,
		// FIX ME
		// CustomProperties: settings.CustomEditProperties,
	}

	create_handler, err := www.CreateFeatureHandler(create_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create create feature handler, %w", err)
	}

	create_handler = maps.AppendResourcesHandlerWithProvider(create_handler, map_provider, maps_opts)

	// FIX ME
	// create_handler = wasm_exec.AppendResourcesHandler(create_handler, wasm_exec_opts)

	// FIX ME
	// create_handler = appendCustomMiddlewareHandlers(settings, uris_table.CreateFeature, create_handler)

	create_handler = bootstrap.AppendResourcesHandler(create_handler, bootstrap_opts)
	create_handler = www.AppendResourcesHandler(create_handler, www_opts.WithCreateHandlerResources())

	// FIX ME
	// create_handler = settings.CustomChrome.WrapHandler(create_handler, "whosonfirst.browser.create")

	create_handler = authenticator.WrapHandler(create_handler)
	return create_handler, nil
}
