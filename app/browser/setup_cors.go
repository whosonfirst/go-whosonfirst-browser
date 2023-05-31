package browser

import (
	"sync"

	"github.com/rs/cors"
)

var setupCORSOnce sync.Once
var setupCORSError error

func setupCORS() {

	if !cfg.EnableCORS {
		return
	}

	cors_origins := cfg.CORSOrigins

	if len(cors_origins) == 0 {
		cors_origins = []string{"*"}
	}

	cors_wrapper = cors.New(cors.Options{
		AllowedOrigins:   cors_origins,
		AllowCredentials: cfg.CORSAllowCredentials,
	})
}
