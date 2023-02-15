package writer

import (
	_ "github.com/whosonfirst/go-writer-github/v3"
)

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	wof_writer "github.com/whosonfirst/go-writer/v3"
)

type WriterOptions struct {
	WriterURIs []string
	Logger     *log.Logger
	Repo       string
}

func NewWriter(ctx context.Context, opts *WriterOptions) (wof_writer.Writer, error) {

	writers := make([]wof_writer.Writer, len(opts.WriterURIs))

	for idx, wr_uri := range opts.WriterURIs {

		// because github_writer.EnsureGitHubAccessToken (in application/browser/browser.go)
		// will URL-encode '{repo}'

		wr_uri, err := url.QueryUnescape(wr_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to unescape '%s', %w", wr_uri, err)
		}

		if opts.Repo != "" && strings.Contains(wr_uri, "{repo}") {

			wr_uri = strings.Replace(wr_uri, "{repo}", opts.Repo, -1)
		}

		wr, err := wof_writer.NewWriter(ctx, wr_uri)

		if err != nil {
			// Note: Don't leak writer URI which may have a GH token in it
			return nil, fmt.Errorf("Failed to create writer at offset %d, %w", idx, err)
		}

		writers[idx] = wr
	}

	multi_opts := &wof_writer.MultiWriterOptions{
		Logger:  opts.Logger,
		Writers: writers,
	}

	multi_wr, err := wof_writer.NewMultiWriterWithOptions(ctx, multi_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create multi writer, %w", err)
	}

	return multi_wr, nil
}
