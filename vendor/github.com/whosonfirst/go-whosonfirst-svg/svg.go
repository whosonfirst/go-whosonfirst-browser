package svg

import (
	"fmt"
	"github.com/fapian/geojson2svg/pkg/geojson2svg"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
	"github.com/whosonfirst/go-whosonfirst-spr/util"
	"io"
	"os"
	"strings"
)

type Options struct {
	Width  float64
	Height float64
	Writer io.Writer
}

func NewDefaultOptions() *Options {

	opts := Options{
		Width:  1024.0,
		Height: 1024.0,
		Writer: os.Stdout,
	}

	return &opts
}

func FeatureToSVG(f geojson.Feature, opts *Options) error {

	geom, err := geometry.ToString(f)

	if err != nil {
		return err
	}

	s := geojson2svg.New()

	err = s.AddGeometry(geom)

	if err != nil {
		return err
	}

	spr, err := f.SPR()

	if err != nil {
		return err
	}

	attrs, err := util.SPRToMap(spr)

	if err != nil {
		return err
	}

	attrs["viewBox"] = fmt.Sprintf("0 0 %0.2f %0.2f", opts.Width, opts.Height)
	
	namespaces := map[string]string{
		"xmlns": "http://www.w3.org/2000/svg",
	}

	for k, _ := range attrs {

		parts := strings.Split(k, ":")

		if len(parts) != 2 {
			continue
		}

		prefix := parts[0]

		_, ok := namespaces[prefix]

		if ok {
			continue
		}

		ns := fmt.Sprintf("xmlns:%s", prefix)
		uri := fmt.Sprintf("x-urn:namespaces#%s", prefix)

		namespaces[ns] = uri
	}

	for ns, uri := range namespaces {
		attrs[ns] = uri
	}

	s_opts := geojson2svg.WithAttributes(attrs)

	rsp := s.Draw(opts.Width, opts.Height, s_opts)
	_, err = opts.Writer.Write([]byte(rsp))

	if err != nil {
		return err
	}

	return nil
}
