package export

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-browser/schema"
	"github.com/whosonfirst/go-whosonfirst-export"
	"github.com/whosonfirst/go-whosonfirst-export/options"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/go-writer"
	"io/ioutil"
	"log"
)

// TBD: passing update/commit messages
// TBD: how strict to be about validating properties
// TBD: validating geometries
// (20191230/thisisaaronland)

func ExportFeature(ctx context.Context, wr writer.Writer, body []byte) ([]byte, error) {

	props_rsp := gjson.GetBytes(body, "properties")

	if !props_rsp.Exists() {
		return nil, errors.New("Missing properties")
	}

	_, err := schema.HasValidProperties(props_rsp.Value())

	if err != nil {
		// return nil, err
		log.Printf("warning: %s\n", err.Error())
	}

	ex_opts, err := options.NewDefaultOptions()

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)

	err = export.Export(body, ex_opts, bw)

	if err != nil {
		return nil, err
	}

	bw.Flush()

	exported_body := buf.Bytes()

	id_rsp := gjson.GetBytes(exported_body, "properties.wof:id")

	if !id_rsp.Exists() {
		return nil, errors.New("Missing wof:id")
	}

	id := id_rsp.Int()

	rel_path, err := uri.Id2RelPath(id)

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(exported_body)
	fh := ioutil.NopCloser(br)

	err = wr.Write(ctx, rel_path, fh)

	if err != nil {
		return nil, err
	}

	return exported_body, nil
}
