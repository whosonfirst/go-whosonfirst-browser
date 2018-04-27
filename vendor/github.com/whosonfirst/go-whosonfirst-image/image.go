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

func FeatureToPNG(f geojson.Feature, fh io.Writer) error {

	img, err := FeatureToImage(f)

	if err != nil {
		return err
	}

	err = png.Encode(fh, img)

	if err != nil {
		return err
	}

	return nil
}

func FeatureToImage(f geojson.Feature) (image.Image, error) {

	tmpfile, err := ioutil.TempFile("", "svg")

	if err != nil {
		return nil, err
	}

	defer os.Remove(tmpfile.Name())

	// log.Println("TMP", tmpfile.Name())

	opts := svg.NewDefaultOptions()
	opts.Writer = tmpfile

	err = svg.FeatureToSVG(f, opts)

	if err != nil {
		return nil, err
	}

	icon, err := oksvg.ReadIcon(tmpfile.Name(), oksvg.StrictErrorMode)

	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(image.Rect(0, 0, int(icon.ViewBox.W), int(icon.ViewBox.H)))
	painter := rasterx.NewRGBAPainter(img)
	raster := rasterx.NewDasher(int(icon.ViewBox.W), int(icon.ViewBox.H))

	icon.Draw(raster, painter, 1.0)

	return img, nil
}
