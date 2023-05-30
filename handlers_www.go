package browser

import (
	"context"
	"fmt"
	"net/http"
)

func wwwIndexHandlerFunc(ctx context.Context) (http.Handler, error) {

	index_opts := www.IndexHandlerOptions{
		Templates:    html_t,
		URIs:         uris_table,
		Capabilities: capabilities,
	}

	index_handler, err := www.IndexHandler(index_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create Index handler, %w", err)
	}

	index_handler = bootstrap.AppendResourcesHandler(index_handler, bootstrap_opts)
	index_handler = www.AppendResourcesHandler(index_handler, www_opts)
	index_handler = maps.AppendResourcesHandlerWithProvider(index_handler, settings.MapProvider, maps_opts)
	index_handler = settings.CustomChrome.WrapHandler(index_handler, "whosonfirst.browser.index")
	index_handler = settings.Authenticator.WrapHandler(index_handler)

	return index_handler, nil
}

func wwwIdHandlerFunc(ctx context.Context) (http.Handler, error) {

	id_opts := www.IDHandlerOptions{
		Templates:     t,
		URIs:          uris_table,
		Capabilities:  capabilities,
		Reader:        settings.Reader,
		Logger:        logger,
		MapProvider:   settings.MapProvider.Scheme(),
		Authenticator: settings.Authenticator,
	}

	id_handler, err := www.IDHandler(id_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create Id handler, %w", err)
	}

	id_handler = bootstrap.AppendResourcesHandler(id_handler, bootstrap_opts)
	id_handler = www.AppendResourcesHandler(id_handler, www_opts.WithIdHandlerResources())
	id_handler = maps.AppendResourcesHandlerWithProvider(id_handler, settings.MapProvider, maps_opts)
	id_handler = settings.CustomChrome.WrapHandler(id_handler, "whosonfirst.browser.id")
	id_handler = settings.Authenticator.WrapHandler(id_handler)

	return id_handler, nil
}

func wwwSearchHandlerFunc(ctx context.Context) (http.Handler, error) {

	search_opts := www.SearchHandlerOptions{
		Templates:    html_t,
		URIs:         uris_table,
		Capabilities: capabilities,
		Database:     settings.SearchDatabase,
		MapProvider:  settings.MapProvider.Scheme(),
	}

	search_handler, err := www.SearchHandler(search_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create Search handler, %w", err)
	}

	search_handler = bootstrap.AppendResourcesHandler(search_handler, bootstrap_opts)
	search_handler = www.AppendResourcesHandler(search_handler, www_opts)
	search_handler = maps.AppendResourcesHandlerWithProvider(search_handler, settings.MapProvider, maps_opts)
	search_handler = settings.CustomChrome.WrapHandler(search_handler, "whosonfirst.browser.search")
	search_handler = settings.Authenticator.WrapHandler(search_handler)

	return search_handler, nil
}
