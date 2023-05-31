package browser

import (
	"context"
	"fmt"
	"sync"

	"github.com/aaronland/go-http-maps/provider"
)

var setupMapsOnce sync.Once
var setupMapsError error

func setupMaps() {

	ctx := context.Background()
	var err error

	map_provider, err = provider.NewProvider(ctx, cfg.MapProviderURI)

	if err != nil {
		setupMapsError = fmt.Errorf("Failed to create new map provider, %w", err)
		return
	}

}
