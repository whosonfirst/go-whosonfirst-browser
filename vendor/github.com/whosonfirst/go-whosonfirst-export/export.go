package export

import (
	"bytes"
	"encoding/json"
	"io"

	format "github.com/tomtaylor/go-whosonfirst-format"
	"github.com/whosonfirst/go-whosonfirst-export/options"
	"github.com/whosonfirst/go-whosonfirst-export/properties"
)

type Feature struct {
	Type       string      `json:"type"`
	Id         int64       `json:"id"`
	Properties interface{} `json:"properties"`
	Bbox       []float64   `json:"bbox,omitempty"`
	Geometry   interface{} `json:"geometry"`
}

func Export(feature []byte, opts options.Options, wr io.Writer) error {

	var err error

	feature, err = Prepare(feature, opts)

	if err != nil {
		return err
	}

	feature, err = Format(feature, opts)

	if err != nil {
		return err
	}

	r := bytes.NewReader(feature)
	_, err = io.Copy(wr, r)

	return err
}

func Prepare(feature []byte, opts options.Options) ([]byte, error) {

	var err error

	feature, err = properties.EnsureWOFId(feature, opts.IDProvider())

	if err != nil {
		return nil, err
	}

	feature, err = properties.EnsureRequired(feature)

	if err != nil {
		return nil, err
	}

	feature, err = properties.EnsureEDTF(feature)

	if err != nil {
		return nil, err
	}

	feature, err = properties.EnsureParentId(feature)

	if err != nil {
		return nil, err
	}

	feature, err = properties.EnsureBelongsTo(feature)

	if err != nil {
		return nil, err
	}

	feature, err = properties.EnsureSupersedes(feature)

	if err != nil {
		return nil, err
	}

	feature, err = properties.EnsureSupersededBy(feature)

	if err != nil {
		return nil, err
	}

	feature, err = properties.EnsureTimestamps(feature)

	if err != nil {
		return nil, err
	}

	return feature, nil
}

func Format(feature []byte, opts options.Options) ([]byte, error) {
	var f format.Feature
	json.Unmarshal(feature, &f)
	return format.FormatFeature(&f)
}
