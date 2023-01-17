package writer

import (
	"context"
	"fmt"
	wof_writer "github.com/whosonfirst/go-writer/v3"
	"log"
	"strings"
)

type WriterOptions struct {
	WriterURIs []string
	Logger     *log.Logger
	Repo       string
}

func NewWriter(ctx context.Context, opts *WriterOptions) (wof_writer.Writer, error) {

	writers := make([]wof_writer.Writer, len(opts.WriterURIs))

	for idx, wr_uri := range opts.WriterURIs {

		if opts.Repo != "" && strings.Contains(wr_uri, "{repo}") {
			wr_uri = strings.Replace(wr_uri, "{repo}", opts.Repo, -1)
		}

		wr, err := wof_writer.NewWriter(ctx, wr_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create writer for '%s', %w", wr_uri, err)
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
