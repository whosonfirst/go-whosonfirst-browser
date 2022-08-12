package svg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	geojson_svg "github.com/whosonfirst/go-geojson-svg"
	"github.com/whosonfirst/go-whosonfirst-feature/geometry"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-spr/v2/util"
	"io"
	"os"
	"strings"
)

type StyleFunction func(f []byte) (map[string]string, error)

type Options struct {
	Width         float64
	Height        float64
	Mercator      bool
	Writer        io.Writer
	StyleFunction StyleFunction
}

func NewDefaultOptions() *Options {

	f := NewDefaultStyleFunction()

	opts := Options{
		Width:         1024.0,
		Height:        1024.0,
		Writer:        os.Stdout,
		Mercator:      false,
		StyleFunction: f,
	}

	return &opts
}

func NewDefaultStyleFunction() StyleFunction {

	style_func := func(boyd []byte) (map[string]string, error) {
		attrs := make(map[string]string)
		return attrs, nil
	}

	return style_func
}

func NewDopplrStyleFunction() StyleFunction {

	default_styles := NewDefaultStyleFunction()

	style_func := func(f []byte) (map[string]string, error) {

		attrs, err := default_styles(f)

		if err != nil {
			return nil, err
		}

		pt, err := properties.Placetype(f)

		if err != nil {
			return nil, err
		}

		fill := fmt.Sprintf("fill: %s", str2hex(pt))

		styles := make([]string, 0)
		styles = append(styles, fill)

		attrs["style"] = strings.Join(styles, ";")

		return attrs, nil
	}

	return style_func
}

func NewFillStyleFunction(colour string) StyleFunction {

	default_styles := NewDefaultStyleFunction()

	style_func := func(f []byte) (map[string]string, error) {

		attrs, err := default_styles(f)

		if err != nil {
			return nil, err
		}

		fill := fmt.Sprintf("fill: %s", colour)
		attrs["style"] = fill

		return attrs, nil
	}

	return style_func
}

func FeatureToSVG(f []byte, opts *Options) error {

	geom, err := geometry.Geometry(f)

	if err != nil {
		return err
	}

	orb_geom := geom.Geometry()
	mbr := orb_geom.Bound()

	mbr_w := mbr.Max[0] - mbr.Min[0]
	mbr_h := mbr.Max[1] - mbr.Min[1]

	w := opts.Width
	h := opts.Height

	if mbr_w == mbr_h {
		// pass
	} else if mbr_w > mbr_h {
		h = h * (mbr_h / mbr_w)
	} else {
		w = w * (mbr_w / mbr_h)
	}

	enc_geom, err := geom.MarshalJSON()

	if err != nil {
		return err
	}

	s := geojson_svg.New()
	s.Mercator = opts.Mercator

	err = s.AddGeometry(string(enc_geom))

	if err != nil {
		return err
	}

	wof_spr, err := spr.WhosOnFirstSPR(f)

	if err != nil {
		return err
	}

	attrs, err := util.SPRToMap(wof_spr)

	if err != nil {
		return err
	}

	style_attrs, err := opts.StyleFunction(f)

	if err != nil {
		return err
	}

	for k, v := range style_attrs {

		// TO DO: consult this: https://github.com/srwiley/oksvg/blob/master/doc/SVG_Element_List.txt

		/*
			ok := false

			switch k {
			case "id":
				ok = true
			case "class":
				ok = true
			case "style":
				ok = true
			default:
				// pass
			}

			if !ok {
				msg := fmt.Sprintf("Invalid style attribute '%s'", k)
				return errors.New(msg)
			}
		*/

		attrs[k] = v
	}

	attrs["viewBox"] = fmt.Sprintf("0 0 %0.2f %0.2f", w, h)

	wof_id, err := properties.Id(f)

	if err != nil {
		return err
	}

	wof_pt, err := properties.Placetype(f)

	if err != nil {
		return err
	}

	id := fmt.Sprintf("wof-%s", wof_id)
	attrs["id"] = id

	pt := fmt.Sprintf("wof-%s", wof_pt)

	class, _ := attrs["class"]
	class = fmt.Sprintf("%s %s", class, pt)

	attrs["class"] = strings.Trim(class, " ")

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

	s_opts := geojson_svg.WithAttributes(attrs)

	rsp := s.Draw(w, h, s_opts)
	_, err = opts.Writer.Write([]byte(rsp))

	if err != nil {
		return err
	}

	return nil
}

func str2hex(text string) string {

	hasher := md5.New()
	hasher.Write([]byte(text))

	enc := hex.EncodeToString(hasher.Sum(nil))
	code := enc[0:6]

	return fmt.Sprintf("#%s", code)
}
