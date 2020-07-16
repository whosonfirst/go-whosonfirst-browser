package tilezen

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jtacoma/uritemplates"
	"github.com/paulmach/orb/clip"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-cache"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
)

const DEFAULT_URITEMPLATE string = "https://tile.nextzen.org/tilezen/vector/{version}/{size}/{layer}/{z}/{x}/{y}.{format}?api_key={apikey}"
const DEFAULT_VERSION string = "v1"
const DEFAULT_SIZE string = "512"
const DEFAULT_LAYER = "all"
const DEFAULT_FORMAT = "mvt"
const MAX_ZOOM int = 16

func NewTile(z int, x int, y int) (*Tile, error) {

	tile := &Tile{
		Z:       z,
		X:       x,
		Y:       y,
		Version: DEFAULT_VERSION,
		Size:    DEFAULT_SIZE,
		Layer:   DEFAULT_LAYER,
		Format:  DEFAULT_FORMAT,
	}

	return tile, nil
}

type Tile struct {
	X       int
	Y       int
	Z       int
	Version string
	Size    string
	Layer   string
	Format  string
}

func (t *Tile) URI() string {
	return fmt.Sprintf("%d/%d/%d.%s", t.Z, t.X, t.Y, t.Format)
}

func (t *Tile) String() string {
	return t.URI()
}

type Options struct {
	ApiKey      string
	Origin      string
	Debug       bool
	URITemplate string
}

func IsOverZoom(z int) bool {

	if z > MAX_ZOOM {
		return true
	}

	return false
}

// this does not account for version, size or layer yet

func ParseURI(uri string) (*Tile, error) {

	re_path, err := regexp.Compile(`(?:(.*)/)?(\d+)/(\d+)/(\d+).(\w+)$`)

	if err != nil {
		return nil, err
	}

	if !re_path.MatchString(uri) {
		return nil, errors.New("Invalid URI")
	}

	m := re_path.FindStringSubmatch(uri)

	z, err := strconv.Atoi(m[2])

	if err != nil {
		return nil, err
	}

	x, err := strconv.Atoi(m[3])

	if err != nil {
		return nil, err
	}

	y, err := strconv.Atoi(m[4])

	if err != nil {
		return nil, err
	}

	format := m[5]

	tile, err := NewTile(z, x, y)

	if err != nil {
		return nil, err
	}

	tile.Format = format
	return tile, nil
}

func FetchTileWithCache(ctx context.Context, tile_cache cache.Cache, tile *Tile, opts *Options) (io.ReadCloser, error) {

	cache_key := tile.URI()

	t_rsp, err := tile_cache.Get(ctx, cache_key)

	if err != nil {

		if !cache.IsCacheMiss(err) {
			return nil, err
		}

		t_rsp, err = FetchTile(tile, opts)

		if err != nil {
			return nil, err
		}

		t_rsp, err = tile_cache.Set(ctx, cache_key, t_rsp)

		if err != nil {
			return nil, err
		}
	}

	return t_rsp, nil
}

func FetchTile(t *Tile, opts *Options) (io.ReadCloser, error) {

	z := t.Z
	x := t.X
	y := t.Y

	fetch_z := z
	fetch_x := x
	fetch_y := y

	// see notes below about whether or not we keep the overzooming code
	// in this package or in tile/rasterzen.go (20190606/thisisaaronland)

	overzoom := IsOverZoom(z)

	if overzoom && t.Format != "json" {
		return nil, errors.New("Overzooming is only supported for `.json` tiles")
	}

	if overzoom {

		max := MAX_ZOOM
		mag := z - max

		ux := uint(x) >> uint(mag)
		uy := uint(y) >> uint(mag)

		fetch_z = max
		fetch_x = int(ux)
		fetch_y = int(uy)
	}

	layer := "all"

	values := make(map[string]interface{})
	values["version"] = t.Version
	values["layer"] = t.Layer
	values["size"] = t.Size
	values["format"] = t.Format
	values["apikey"] = opts.ApiKey
	values["z"] = fetch_z
	values["x"] = fetch_x
	values["y"] = fetch_y

	template := DEFAULT_URITEMPLATE

	if opts.URITemplate != "" {
		template = opts.URITemplate
	}

	endpoint, err := uritemplates.Parse(template)

	if err != nil {
		return nil, err
	}

	url, err := endpoint.Expand(values)

	if err != nil {
		return nil, err
	}

	cl := new(http.Client)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	if opts.Origin != "" {
		req.Header.Set("Origin", opts.Origin)
	}

	if opts.Debug {

		dump, err := httputil.DumpRequest(req, false)

		if err != nil {
			return nil, err
		}

		log.Println(string(dump))
	}

	rsp, err := cl.Do(req)

	if err != nil {
		return nil, err
	}

	if opts.Debug {
		log.Println(url, rsp.Status)
	}

	if rsp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Nextzen returned a non-200 response fetching %s/%d/%d/%d : '%s'", layer, z, x, y, rsp.Status))
	}

	rsp_body := rsp.Body

	// overzooming works until it doesn't - specifically there are
	// weird offsets that I don't understand - examples include:
	// ./bin/rasterd -www -www-debug -no-cache -nextzen-debug -nextzen-apikey {KEY}
	// http://localhost:8080/#18/37.61800/-122.38301
	// http://localhost:8080/#19/37.61780/-122.38800
	// http://localhost:8080/svg/19/83903/202936.svg?api_key={KEY}
	// (20190606/thisisaaronland)

	if overzoom {

		// it would be good to cache rsp_body (aka the Z16 tile) here or maybe
		// we just move all of this logic in to tile/rasterzen.go...
		// (20190606/thisisaaronland)

		cropped_rsp, err := CropTile(z, x, y, rsp_body)

		if err != nil {
			return nil, err
		}

		rsp_body = cropped_rsp
	}

	return rsp_body, nil
}

// crop all the elements in fh to the bounds of (z, x, y)

func CropTile(z int, x int, y int, fh io.ReadCloser) (io.ReadCloser, error) {

	zm := maptile.Zoom(uint32(z))
	tl := maptile.New(uint32(x), uint32(y), zm)

	bounds := tl.Bound()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	cropped_tile := make(map[string]interface{})

	type CroppedResponse struct {
		Layer             string
		FeatureCollection *geojson.FeatureCollection
	}

	done_ch := make(chan bool)
	err_ch := make(chan error)
	rsp_ch := make(chan CroppedResponse)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Layers is defined in layers.go
	// PLEASE FIX ME TO DO ALL LAYERS AND REMOVE layers.go

	for _, layer_name := range Layers {

		go func(layer_name string) {

			defer func() {
				done_ch <- true
			}()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			fc_rsp := gjson.GetBytes(body, layer_name)

			if !fc_rsp.Exists() {
				return
			}

			fc_str := fc_rsp.String()

			fc, err := geojson.UnmarshalFeatureCollection([]byte(fc_str))

			if err != nil {
				err_ch <- err
				return
			}

			cropped_fc := geojson.NewFeatureCollection()

			for _, f := range fc.Features {

				geom := f.Geometry
				clipped_geom := clip.Geometry(bounds, geom)

				if clipped_geom == nil {
					continue
				}

				f.Geometry = clipped_geom
				cropped_fc.Append(f)
			}

			if len(cropped_fc.Features) > 0 {

				rsp := CroppedResponse{
					Layer:             layer_name,
					FeatureCollection: cropped_fc,
				}

				rsp_ch <- rsp
			}

		}(layer_name)
	}

	remaining := len(Layers)

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case err := <-err_ch:
			return nil, err
		case rsp := <-rsp_ch:
			cropped_tile[rsp.Layer] = rsp.FeatureCollection
		}
	}

	cropped_body, err := json.Marshal(cropped_tile)

	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(cropped_body)
	return ioutil.NopCloser(r), nil
}
