package emitter

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/filters"
	"io"
)

func init() {
	ctx := context.Background()
	RegisterEmitter(ctx, "geojsonl", NewGeoJSONLEmitter)
}

type GeojsonLEmitter struct {
	Emitter
	filters filters.Filters
}

func NewGeoJSONLEmitter(ctx context.Context, uri string) (Emitter, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, err
	}

	idx := &GeojsonLEmitter{
		filters: f,
	}

	return idx, nil
}

func (idx *GeojsonLEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(ctx, uri)

	if err != nil {
		return err
	}

	defer fh.Close()

	// see this - we're using ReadLine because it's entirely possible
	// that the raw GeoJSON (LS) will be too long for bufio.Scanner
	// see also - https://golang.org/pkg/bufio/#Reader.ReadLine
	// (20170822/thisisaaronland)

	reader := bufio.NewReader(fh)
	raw := bytes.NewBuffer([]byte(""))

	i := 0

	for {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		path := fmt.Sprintf("%s#%d", uri, i)
		i += 1

		fragment, is_prefix, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		raw.Write(fragment)

		if is_prefix {
			continue
		}

		br := bytes.NewReader(raw.Bytes())
		fh, err := ioutil.NewReadSeekCloser(br)

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
		err = index_cb(ctx, fh)

		if err != nil {
			return err
		}

		raw.Reset()
	}

	return nil
}
