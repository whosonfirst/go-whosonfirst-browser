package api

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-whosonfirst-browser/v6/http"
	"github.com/whosonfirst/go-whosonfirst-browser/v6/writer"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type publishFeatureOptions struct {
	Logger     *log.Logger
	WriterURIs []string
	Exporter   export.Exporter
	Cache      cache.Cache
	URI        *http.URI
}

func publishFeature(ctx context.Context, opts *publishFeatureOptions, body []byte) ([]byte, error) {

	repo, err := properties.Repo(body)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive repo from body, %w", err)
	}

	wr_opts := &writer.WriterOptions{
		WriterURIs: opts.WriterURIs,
		Logger:     opts.Logger,
		Repo:       repo,
	}

	wr, err := writer.NewWriter(ctx, wr_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new writer, %w", err)
	}

	exp_body, err := opts.Exporter.Export(ctx, body)

	if err != nil {
		return nil, fmt.Errorf("Failed to export body, %w", err)
	}

	var rel_path string

	if opts.URI == nil {

		wof_id, err := properties.Id(exp_body)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive ID from body, %w", err)
		}

		// To do: alternate geometries

		rel_path, err = uri.Id2RelPath(wof_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive rel path, %w", err)
		}

	} else {

		rel_path, err = uri.Id2RelPath(opts.URI.Id, opts.URI.URIArgs)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive rel path, %w", err)
		}
	}

	br := bytes.NewReader(exp_body)

	_, err = wr.Write(ctx, rel_path, br)

	if err != nil {
		return nil, fmt.Errorf("Failed to write body, %w", err)
	}

	err = opts.Cache.Unset(ctx, rel_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to unset cache for '%s', %w", rel_path, err)
	}

	return exp_body, nil
}
