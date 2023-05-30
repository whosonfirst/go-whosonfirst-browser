package browser

func urisHandlerFunc(ctx context.Context) (http.Handler, error) {

	uris_t := js_t.Lookup("uris")

	if uris_t == nil {
		return fmt.Errorf("Failed to load 'uris' javascript template")
	}

	uris_opts := &www.URIsHandlerOptions{
		URIs:     uris_table,
		Template: uris_t,
	}

	uris_handler, err := www.URIsHandler(uris_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create URIs handler, %w", err)
	}

	return uris_handler, nil
}
