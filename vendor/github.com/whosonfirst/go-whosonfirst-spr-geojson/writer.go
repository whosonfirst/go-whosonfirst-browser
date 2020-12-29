package geojson

import (
	"context"
	"github.com/whosonfirst/go-reader"
	"io"
	"sync/atomic"
)

type FeatureCollectionWriter struct {
	reader reader.Reader
	writer io.Writer
	index  uint32
}

func NewFeatureCollectionWriter(r reader.Reader, wr io.Writer) (*FeatureCollectionWriter, error) {

	fc := &FeatureCollectionWriter{
		reader: r,
		writer: wr,
		index:  uint32(0),
	}

	return fc, nil
}

func (wr *FeatureCollectionWriter) Begin() error {
	_, err := io.WriteString(wr.writer, `{"type":"FeatureCollection", "features": [`)
	return err
}

func (wr *FeatureCollectionWriter) WriteFeature(ctx context.Context, rel_path string) error {

	if atomic.LoadUint32(&wr.index) > 0 {

		_, err := io.WriteString(wr.writer, `,`)

		if err != nil {
			return err
		}
	}

	atomic.AddUint32(&wr.index, 1)

	fh, err := wr.reader.Read(ctx, rel_path)

	if err != nil {
		return err
	}

	defer fh.Close()

	// Something something something strip white space from fh here...

	_, err = io.Copy(wr.writer, fh)
	return err
}

func (wr *FeatureCollectionWriter) End() error {
	_, err := io.WriteString(wr.writer, `]}`)
	return err
}
