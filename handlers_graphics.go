package browser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-browser/v7/http/www"	
)

func pngHandlerFunc(ctx context.Context) (http.Handler, error) {

	sizes := www.DefaultRasterSizes()

	png_opts := &www.RasterHandlerOptions{
		Sizes:  sizes,
		Format: "png",
		Reader: wof_reader,
		Logger: logger,
	}

	png_handler, err := www.RasterHandler(png_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create raster/png handler, %w", err)
	}

	return png_handler, nil
}

func svgHandlerFunc(ctx context.Context) (http.Handler, error) {

	sizes := www.DefaultSVGSizes()

	svg_opts := &www.SVGHandlerOptions{
		Sizes:  sizes,
		Reader: wof_reader,
		Logger: logger,
	}

	svg_handler, err := www.SVGHandler(svg_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create SVG handler, %w", err)
	}

	if cors_wrapper != nil {
		svg_handler = cors_wrapper.Handler(svg_handler)
	}

	return svg_handler, nil
}
