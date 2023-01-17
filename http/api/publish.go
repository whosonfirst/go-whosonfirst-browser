package api

import (
	"bytes"
	"context"
	"fmt"
	"log"

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

	br := bytes.NewReader(exp_body)

	rel_path, err := uri.Id2RelPath(opts.URI.Id, opts.URI.URIArgs)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive rel path, %w", err)
	}

	_, err = wr.Write(ctx, rel_path, br)

	if err != nil {
		return nil, fmt.Errorf("Failed to write body, %w", err)
	}

	return exp_body, nil
}
