package browser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
)

func jsURIsHandlerFunc(ctx context.Context) (http.Handler, error) {

	uris_t := js_t.Lookup("uris")

	if uris_t == nil {
		return nil, fmt.Errorf("Failed to load 'uris' javascript template")
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
