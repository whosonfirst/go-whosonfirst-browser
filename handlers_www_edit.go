package browser

func wwwEditGeometryHandlerFunc(ctx context.Context) (http.Handler, error) {

	geom_t := t.Lookup("geometry")

	if geom_t == nil {
		return fmt.Errorf("Failed to load 'geometry' template")
	}

	geom_opts := &www.EditGeometryHandlerOptions{
		Authenticator: settings.Authenticator,
		MapProvider:   settings.MapProvider.Scheme(),
		URIs:          uris_table,
		Capabilities:  capabilities,
		Template:      geom_t,
		Logger:        logger,
		Reader:        settings.Reader,
	}

	geom_handler, err := www.EditGeometryHandler(geom_opts)

	if err != nil {
		return fmt.Errorf("Failed to create edit geometry handler, %w", err)
	}

	geom_handler = maps.AppendResourcesHandlerWithProvider(geom_handler, settings.MapProvider, maps_opts)
	geom_handler = bootstrap.AppendResourcesHandler(geom_handler, bootstrap_opts)
	geom_handler = www.AppendResourcesHandler(geom_handler, www_opts.WithGeometryHandlerResources())
	geom_handler = settings.CustomChrome.WrapHandler(geom_handler, "whosonfirst.browser.geometry")
	geom_handler = settings.Authenticator.WrapHandler(geom_handler)

	return geom_handler, nil
}

func wwwCreateGeometryHandlerFunc(ctx context.Context) (http.Handler, error) {

	create_t := html_t.Lookup("create")

	if create_t == nil {
		return fmt.Errorf("Failed to load 'create' template")
	}

	create_opts := &www.CreateFeatureHandlerOptions{
		Authenticator:    settings.Authenticator,
		MapProvider:      settings.MapProvider.Scheme(),
		URIs:             uris_table,
		Capabilities:     capabilities,
		Template:         create_t,
		Logger:           logger,
		Reader:           settings.Reader,
		CustomProperties: settings.CustomEditProperties,
	}

	create_handler, err := www.CreateFeatureHandler(create_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create create feature handler, %w", err)
	}

	create_handler = maps.AppendResourcesHandlerWithProvider(create_handler, settings.MapProvider, maps_opts)
	create_handler = wasm_exec.AppendResourcesHandler(create_handler, wasm_exec_opts)

	// create_handler = appendCustomMiddlewareHandlers(settings, uris_table.CreateFeature, create_handler)

	create_handler = bootstrap.AppendResourcesHandler(create_handler, bootstrap_opts)
	create_handler = www.AppendResourcesHandler(create_handler, www_opts.WithCreateHandlerResources())

	create_handler = settings.CustomChrome.WrapHandler(create_handler, "whosonfirst.browser.create")

	create_handler = settings.Authenticator.WrapHandler(create_handler)

}
