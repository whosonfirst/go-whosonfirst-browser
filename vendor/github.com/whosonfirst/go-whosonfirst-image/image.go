package image

import (
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-svg"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	_ "log"
	"os"
)

type Options struct {
	Width  int
	Height int
	Writer io.Writer
}

func NewDefaultOptions() *Options {

	opts := Options{
		Width:  1024,
		Height: 1024,
		Writer: os.Stdout,
	}

	return &opts
}

func FeatureToPNG(f geojson.Feature, opts *Options) error {

	img, err := FeatureToImage(f, opts)

	if err != nil {
		return err
	}

	err = png.Encode(opts.Writer, img)

	if err != nil {
		return err
	}

	return nil
}

func FeatureToImage(f geojson.Feature, opts *Options) (image.Image, error) {

	tmpfile, err := ioutil.TempFile("", "svg")

	if err != nil {
		return nil, err
	}

	defer func() {

		_, err := os.Stat(tmpfile.Name())

		if !os.IsNotExist(err) {
			os.Remove(tmpfile.Name())
		}
	}()

	// log.Println("TMP", tmpfile.Name())

	svg_opts := svg.NewDefaultOptions()
	svg_opts.Writer = tmpfile
	svg_opts.Height = float64(opts.Height)
	svg_opts.Width = float64(opts.Width)

	err = svg.FeatureToSVG(f, svg_opts)

	if err != nil {
		return nil, err
	}

	icon, err := oksvg.ReadIcon(tmpfile.Name(), oksvg.StrictErrorMode)

	if err != nil {
		return nil, err
	}

	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)

	icon.Draw(raster, 1.0)

	return img, nil
}
