package browser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aaronland/go-http-ping/v2"
)

func pingHandlerFunc(ctx context.Context) (http.Handler, error) {

	ping_handler, err := ping.PingPongHandler()

	if err != nil {
		return nil, fmt.Errorf("Failed to create ping handler, %w", err)
	}

	return ping_handler, nil
}
