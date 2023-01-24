package browser

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
)

type Settings struct {
	// Placeholder for application settings derived from a *Config instance

	Paths        *www.Paths
	Capabilities *www.Capabilities
	Reader       reader.Reader
	Cache        cache.Cache

	WriterURIs []string
	Exporter   export.Exporter

	Authenticator auth.Authenticator

	Templates []fs.FS

	CustomHandlers map[string]http.HandlerFunc
	
	Verbose bool
}

func SettingsFromConfig(ctx context.Context, cfg *Config) (*Settings, error) {
	return nil, fmt.Errorf("Not implemented")
}
