package emitter

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-iterate/filters"
)

func init() {
	ctx := context.Background()
	RegisterEmitter(ctx, "file", NewFileEmitter)
}

type FileEmitter struct {
	Emitter
	filters filters.Filters
}

func NewFileEmitter(ctx context.Context, uri string) (Emitter, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, err
	}

	idx := &FileEmitter{
		filters: f,
	}

	return idx, nil
}

func (idx *FileEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(ctx, uri)

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

	ctx = AssignPathContext(ctx, uri)
	return index_cb(ctx, fh)
}
