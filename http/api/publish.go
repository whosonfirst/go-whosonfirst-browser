package api

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	aa_log "github.com/aaronland/go-log"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/http"
	"github.com/whosonfirst/go-whosonfirst-browser/v7/writer"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type publishFeatureOptions struct {
	Logger      *log.Logger
	WriterURIs  []string
	Exporter    export.Exporter
	Cache       cache.Cache
	URI         *http.URI
	Account     *auth.Account
	Title       string
	Description string
	Branch      string
}

func publishFeature(ctx context.Context, opts *publishFeatureOptions, body []byte) ([]byte, error) {

	repo, err := properties.Repo(body)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive repo from body, %w", err)
	}

	// START OF on-the-fly values for https://github.com/whosonfirst/go-writer-github based writers

	author := opts.Account.Name
	owner := opts.Account.Name
	email := "editor@localhost"
	title := opts.Title
	description := opts.Description
	branch := opts.Branch

	to_replace := map[string]string{
		"{author}":      author,
		"{owner}":       owner,
		"{email}":       email,
		"{title}":       title,
		"{description}": description,
		"{branch}":      branch,
	}

	// Note how we are creating a fresh (local) set of writer URIs so we don't need to worry
	// about two different instances of a handler changing each other's opts.WriterURIs list

	writer_uris := make([]string, len(opts.WriterURIs))

	for idx, wr_uri := range opts.WriterURIs {

		for match_str, replace_str := range to_replace {
			wr_uri = strings.Replace(wr_uri, match_str, replace_str, -1)
		}

		writer_uris[idx] = wr_uri

		aa_log.Debug(opts.Logger, "Writer URI %d %s", idx, wr_uri)
	}

	// END OF on-the-fly values for https://github.com/whosonfirst/go-writer-github based writers

	wr_opts := &writer.WriterOptions{
		WriterURIs: writer_uris,
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
		return nil, fmt.Errorf("Failed to write body for %s, %w", rel_path, err)
	}

	err = wr.Close(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to close writer for %s, %w", rel_path, err)
	}

	err = opts.Cache.Unset(ctx, rel_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to unset cache for '%s', %w", rel_path, err)
	}

	return exp_body, nil
}
