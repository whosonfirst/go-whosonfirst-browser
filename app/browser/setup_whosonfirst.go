package browser

import (
	"context"
	"fmt"
	"sync"

	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-cachereader"
	github_reader "github.com/whosonfirst/go-reader-github"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	github_writer "github.com/whosonfirst/go-writer-github/v3"
)

var setupWhosOnFirstReaderOnce sync.Once
var setupWhosOnFirstReaderError error

var setupWhosOnFirstWriterOnce sync.Once
var setupWhosOnFirstWriterError error

func setupWhosOnFirstReader() {

	ctx := context.Background()

	reader_uris := cfg.ReaderURIs

	if cfg.GitHubReaderAccessTokenURI != "" {

		if cfg.GitHubReaderAccessTokenURI == "" {
			cfg.GitHubReaderAccessTokenURI = cfg.GitHubAccessTokenURI
		}

		for idx, r_uri := range reader_uris {

			r_uri, err := github_reader.EnsureGitHubAccessToken(ctx, r_uri, cfg.GitHubReaderAccessTokenURI)

			if err != nil {
				setupWhosOnFirstReaderError = fmt.Errorf("Failed to ensure GitHub access token for '%s', %w", r_uri, err)
			}

			reader_uris[idx] = r_uri
		}
	}

	browser_reader, err := reader.NewMultiReaderFromURIs(ctx, reader_uris...)

	if err != nil {
		setupWhosOnFirstReaderError = fmt.Errorf("Failed to create reader, %w", err)
		return
	}

	browser_cache, err := cache.NewCache(ctx, cfg.CacheURI)

	if err != nil {
		setupWhosOnFirstReaderError = fmt.Errorf("Failed to create new cache, %w", err)
		return
	}

	cr_opts := &cachereader.CacheReaderOptions{
		Reader: browser_reader,
		Cache:  browser_cache,
	}

	cr, err := cachereader.NewCacheReaderWithOptions(ctx, cr_opts)

	if err != nil {
		setupWhosOnFirstReaderError = fmt.Errorf("Failed to create cache reader, %w", err)
		return
	}

	wof_reader = cr
	wof_cache = browser_cache
}

func setupWhosOnFirstWriter() {

	ctx := context.Background()
	var err error

	wof_exporter, err = export.NewExporter(ctx, cfg.ExporterURI)

	if err != nil {
		setupWhosOnFirstReaderError = fmt.Errorf("Failed to create new exporter, %w", err)
		return
	}

	wof_writer_uris = cfg.WriterURIs

	if cfg.GitHubReaderAccessTokenURI != "" {

		if cfg.GitHubWriterAccessTokenURI == "" {
			cfg.GitHubWriterAccessTokenURI = cfg.GitHubAccessTokenURI
		}

		for idx, wr_uri := range writer_uris {

			wr_uri, err := github_writer.EnsureGitHubAccessToken(ctx, wr_uri, cfg.GitHubReaderAccessTokenURI)

			if err != nil {
				setupWhosOnFirstReaderError = fmt.Errorf("Failed to ensure GitHub access token for '%s', %w", wr_uri, err)
			}

			wof_writer_uris[idx] = wr_uri
		}
	}

}
