package geojson

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	go_geojson "github.com/paulmach/go.geojson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"io"
	_ "log"
)

type AsFeatureCollectionOptions struct {
	Reader           reader.Reader
	Writer           io.Writer
	SPRPathResolver  SPRPathResolver
	JSONPathResolver JSONPathResolver
}

type ToFeatureCollectionOptions struct {
	SPRPathResolver  SPRPathResolver
	JSONPathResolver JSONPathResolver
	Reader           reader.Reader
}

func ToFeatureCollection(ctx context.Context, rsp spr.StandardPlacesResults, opts *ToFeatureCollectionOptions) (*go_geojson.FeatureCollection, error) {

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	as_opts := &AsFeatureCollectionOptions{
		Reader:           opts.Reader,
		Writer:           wr,
		SPRPathResolver:  opts.SPRPathResolver,
		JSONPathResolver: opts.JSONPathResolver,
	}

	err := AsFeatureCollection(ctx, rsp, as_opts)

	if err != nil {
		return nil, err
	}

	wr.Flush()

	return go_geojson.UnmarshalFeatureCollection(buf.Bytes())
}

func ToFeatureCollectionWithJSON(ctx context.Context, body []byte, opts *ToFeatureCollectionOptions) (*go_geojson.FeatureCollection, error) {

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	as_opts := &AsFeatureCollectionOptions{
		Reader:           opts.Reader,
		Writer:           wr,
		SPRPathResolver:  opts.SPRPathResolver,
		JSONPathResolver: opts.JSONPathResolver,
	}

	err := AsFeatureCollectionWithJSON(ctx, body, as_opts)

	if err != nil {
		return nil, err
	}

	wr.Flush()

	return go_geojson.UnmarshalFeatureCollection(buf.Bytes())
}

func AsFeatureCollection(ctx context.Context, rsp spr.StandardPlacesResults, opts *AsFeatureCollectionOptions) error {

	r := opts.Reader
	wr := opts.Writer

	fc, err := NewFeatureCollectionWriter(r, wr)

	if err != nil {
		return err
	}

	err = fc.Begin()

	if err != nil {
		return err
	}

	for _, pl := range rsp.Results() {

		var path string

		if opts.SPRPathResolver != nil {

			p, err := opts.SPRPathResolver(ctx, pl)

			if err != nil {
				return err
			}

			path = p

		} else {
			path = pl.Path()
		}

		if path == "" {
			return fmt.Errorf("Unable to determine path for ID '%s'", pl.Id())
		}

		err = fc.WriteFeature(ctx, path)

		if err != nil {
			return err
		}
	}

	return fc.End()
}

func AsFeatureCollectionWithJSON(ctx context.Context, body []byte, opts *AsFeatureCollectionOptions) error {

	r := opts.Reader
	wr := opts.Writer

	if opts.JSONPathResolver == nil {
		return errors.New("Missing JSONPathResolver function")
	}

	paths, err := opts.JSONPathResolver(ctx, body)

	if err != nil {
		return err
	}

	fc, err := NewFeatureCollectionWriter(r, wr)

	if err != nil {
		return err
	}

	err = fc.Begin()

	if err != nil {
		return err
	}

	for _, p := range paths {

		err := fc.WriteFeature(ctx, p)

		if err != nil {
			return err
		}
	}

	return fc.End()
}
