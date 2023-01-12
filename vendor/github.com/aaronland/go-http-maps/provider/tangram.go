package provider

import (
	"context"
	"fmt"
	"github.com/aaronland/go-http-leaflet"
	"github.com/aaronland/go-http-tangramjs"
	tilepack_http "github.com/tilezen/go-tilepacks/http"
	"github.com/tilezen/go-tilepacks/tilepack"
	"io"
	"log"
	"net/http"
	"net/url"
)

const TANGRAM_SCHEME string = "tangram"

type TangramProvider struct {
	Provider
	leafletOptions *leaflet.LeafletOptions
	tangramOptions *tangramjs.TangramJSOptions
	tilezenOptions *TilezenOptions
	logger         *log.Logger
}

func init() {

	tangramjs.APPEND_LEAFLET_RESOURCES = false
	tangramjs.APPEND_LEAFLET_ASSETS = false

	ctx := context.Background()
	RegisterProvider(ctx, TANGRAM_SCHEME, NewTangramProvider)
}

func TangramJSOptionsFromURL(u *url.URL) (*tangramjs.TangramJSOptions, error) {

	opts := tangramjs.DefaultTangramJSOptions()

	q := u.Query()

	opts.NextzenOptions.APIKey = q.Get("nextzen-apikey")

	q_style_url := q.Get("nextzen-style-url")

	if q_style_url != "" {
		opts.NextzenOptions.StyleURL = q_style_url
	}

	q_tile_url := q.Get("nextzen-tile-url")

	if q_style_url != "" {
		opts.NextzenOptions.TileURL = q_tile_url
	}

	return opts, nil
}

func NewTangramProvider(ctx context.Context, uri string) (Provider, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	leaflet_opts, err := LeafletOptionsFromURL(u)

	if err != nil {
		return nil, fmt.Errorf("Failed to create leaflet options, %w", err)
	}

	tangram_opts, err := TangramJSOptionsFromURL(u)

	if err != nil {
		return nil, fmt.Errorf("Failed to create tilezen options, %w", err)
	}

	tilezen_opts, err := TilezenOptionsFromURL(u)

	if err != nil {
		return nil, fmt.Errorf("Failed to create tilezen options, %w", err)
	}

	logger := log.New(io.Discard, "", 0)

	p := &TangramProvider{
		leafletOptions: leaflet_opts,
		tangramOptions: tangram_opts,
		tilezenOptions: tilezen_opts,
		logger:         logger,
	}

	return p, nil
}

func (p *TangramProvider) Scheme() string {
	return TANGRAM_SCHEME
}

func (p *TangramProvider) AppendResourcesHandler(handler http.Handler) http.Handler {
	return p.AppendResourcesHandlerWithPrefix(handler, "")
}

func (p *TangramProvider) AppendResourcesHandlerWithPrefix(handler http.Handler, prefix string) http.Handler {
	handler = leaflet.AppendResourcesHandlerWithPrefix(handler, p.leafletOptions, prefix)
	handler = tangramjs.AppendResourcesHandlerWithPrefix(handler, p.tangramOptions, prefix)
	return handler
}

func (p *TangramProvider) AppendAssetHandlers(mux *http.ServeMux) error {
	return p.AppendAssetHandlersWithPrefix(mux, "")
}

func (p *TangramProvider) AppendAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {

	err := leaflet.AppendAssetHandlersWithPrefix(mux, prefix)

	if err != nil {
		return fmt.Errorf("Failed to append leaflet asset handler, %w", err)
	}

	err = tangramjs.AppendAssetHandlersWithPrefix(mux, prefix)

	if err != nil {
		return fmt.Errorf("Failed to append tangram asset handler, %w", err)
	}

	if p.tilezenOptions.EnableTilepack {

		tilepack_reader, err := tilepack.NewMbtilesReader(p.tilezenOptions.TilepackPath)

		if err != nil {
			return fmt.Errorf("Failed to create tilepack reader, %w", err)
		}

		tilepack_url := p.tilezenOptions.TilepackURL

		if prefix != "" {

			tilepack_url, err = url.JoinPath(prefix, tilepack_url)

			if err != nil {
				return fmt.Errorf("Failed to join path with %s and %s", prefix, tilepack_url)
			}
		}

		tilepack_handler := tilepack_http.MbtilesHandler(tilepack_reader)
		mux.Handle(tilepack_url, tilepack_handler)
	}

	return nil
}

func (p *TangramProvider) SetLogger(logger *log.Logger) error {
	p.logger = logger
	return nil
}
