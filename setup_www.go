package browser

import (
	"context"
	"fmt"
	"sync"

	"github.com/sfomuseum/go-template/html"
	"github.com/sfomuseum/go-template/text"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/templates/javascript"
)

var setupWWWOnce sync.Once
var setupWWWError error

func setupWWW() {

	ctx := context.Background()
	var err error

	html_t, err = html.LoadTemplates(ctx, cfg.Templates...)

	if err != nil {
		setupWWWError = fmt.Errorf("Failed to load templates, %w", err)
		return
	}

	js_t, err = text.LoadTemplatesMatching(ctx, "*.js", javascript.FS)

	if err != nil {
		setupWWWError = fmt.Errorf("Failed to load JS templates, %w", err)
		return
	}

	setupStaticOnce.Do(setupStatic)

	if setupStaticError != nil {
		setupWWWError = fmt.Errorf("Failed to configure static setup, %w", setupStaticError)
	}

	setupMapsOnce.Do(setupStatic)

	if setupMapsError != nil {
		setupWWWError = fmt.Errorf("Failed to configure static setup, %w", setupMapsError)
	}

	setupCORSOnce.Do(setupCORS)
}
