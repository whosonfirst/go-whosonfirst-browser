package emitter

import (
	"bufio"
	"context"
	"github.com/whosonfirst/go-whosonfirst-iterate/filters"
)

func init() {
	ctx := context.Background()
	RegisterEmitter(ctx, "filelist", NewFileListEmitter)
}

type FileListEmitter struct {
	Emitter
	filters filters.Filters
}

func NewFileListEmitter(ctx context.Context, uri string) (Emitter, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, err
	}

	idx := &FileListEmitter{
		filters: f,
	}

	return idx, nil
}

func (idx *FileListEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(ctx, uri)

	if err != nil {
		return err
	}

	defer fh.Close()

	scanner := bufio.NewScanner(fh)

	for scanner.Scan() {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		path := scanner.Text()

		fh, err := ReaderWithPath(ctx, path)

		if err != nil {
			return err
		}

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

		err = index_cb(ctx, fh)

		if err != nil {
			return err
		}
	}

	err = scanner.Err()

	if err != nil {
		return err
	}

	return nil
}
