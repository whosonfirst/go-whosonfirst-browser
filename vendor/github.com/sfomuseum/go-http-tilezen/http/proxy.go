package http

import (
	"context"
	"github.com/sfomuseum/go-tilezen"
	"github.com/whosonfirst/go-cache"
	"io"
	gohttp "net/http"
	"time"
)

type TilezenProxyHandlerOptions struct {
	Cache   cache.Cache
	Timeout time.Duration
}

func TilezenProxyHandler(proxy_opts *TilezenProxyHandlerOptions) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		path := req.URL.Path

		tile, err := tilezen.ParseURI(path)

		if err != nil {
			gohttp.Error(rsp, "Invalid path", gohttp.StatusBadRequest)
			return
		}

		q := req.URL.Query()

		api_key := q.Get("api_key")

		if api_key == "" {
			gohttp.Error(rsp, "Missing API key", gohttp.StatusBadRequest)
			return
		}

		tilezen_opts := &tilezen.Options{
			ApiKey: api_key,
		}

		ctx, cancel := context.WithTimeout(context.Background(), proxy_opts.Timeout)
		defer cancel()

		t_rsp, err := tilezen.FetchTileWithCache(ctx, proxy_opts.Cache, tile, tilezen_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		defer t_rsp.Close()

		_, err = io.Copy(rsp, t_rsp)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	return gohttp.HandlerFunc(fn), nil
}
