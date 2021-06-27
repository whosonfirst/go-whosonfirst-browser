package emitter

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-iterate/filters"
	"os"
	"path/filepath"
)

func init() {
	ctx := context.Background()
	RegisterEmitter(ctx, "directory", NewDirectoryEmitter)
}

type DirectoryEmitter struct {
	Emitter
	filters filters.Filters
}

func NewDirectoryEmitter(ctx context.Context, uri string) (Emitter, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, err
	}

	idx := &DirectoryEmitter{
		filters: f,
	}

	return idx, nil
}

func (idx *DirectoryEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {

	abs_path, err := filepath.Abs(uri)

	if err != nil {
		return err
	}

	crawl_cb := func(path string, info os.FileInfo) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		if info.IsDir() {
			return nil
		}

		fh, err := ReaderWithPath(ctx, path)

		if err != nil {
			return err
		}

		defer fh.Close()

		if idx.filters != nil {

			ok, err := idx.filters.Apply(ctx, fh)

			if err != nil {
				return err
			}

			if !ok {
				return nil
			}

			_, err = fh.Seek(0, 0)

			if err != nil {
				return err
			}
		}

		ctx = AssignPathContext(ctx, path)
		return index_cb(ctx, fh)
	}

	c := crawl.NewCrawler(abs_path)
	return c.Crawl(crawl_cb)
}
