package emitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/filters"
	"io"
)

func init() {
	ctx := context.Background()
	RegisterEmitter(ctx, "featurecollection", NewFeatureCollectionEmitter)
}

type FeatureCollectionEmitter struct {
	Emitter
	filters filters.Filters
}

func NewFeatureCollectionEmitter(ctx context.Context, uri string) (Emitter, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, err
	}

	i := &FeatureCollectionEmitter{
		filters: f,
	}

	return i, nil
}

func (idx *FeatureCollectionEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(ctx, uri)

	if err != nil {
		return err
	}

	defer fh.Close()

	body, err := io.ReadAll(fh)

	if err != nil {
		return err
	}

	type FC struct {
		Type     string
		Features []interface{}
	}

	var collection FC

	err = json.Unmarshal(body, &collection)

	if err != nil {
		return err
	}

	for i, f := range collection.Features {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		feature, err := json.Marshal(f)

		if err != nil {
			return err
		}

		br := bytes.NewReader(feature)
		fh, err := ioutil.NewReadSeekCloser(br)

		if err != nil {
			return err
		}

		if idx.filters != nil {

			ok, err := idx.filters.Apply(ctx, fh)

			if err != nil {
				return err
			}

			if !ok {
				continue
			}

			_, err = fh.Seek(0, 0)

			if err != nil {
				return err
			}
		}

		path := fmt.Sprintf("%s#%d", uri, i)
		ctx = AssignPathContext(ctx, path)

		err = index_cb(ctx, fh)

		if err != nil {
			return err
		}
	}

	return nil
}
