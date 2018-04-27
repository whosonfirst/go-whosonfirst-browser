# go-whosonfirst-image

Use `go-whosonfirst-svg` and `oksvg` to render Who's On First features as raster images.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

It's early days, still.

## Example

```
import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-image"
	"image/png"
	"os"
)

func main() {

     	path := "/path/to/feature.geojson"
	
	f, _ := feature.LoadFeatureFromFile(path)

	opts := image.NewDefaultOptions()
	img, _ := image.FeatureToImage(f, opts)

	png.Encode(opts.Writer, img)
}
```

_Error handling removed for brevity._

## Tools

### wof-feature2png

```
./bin/wof-feature2png /usr/local/data/whosonfirst-data-constituency-us/data/110/874/663/7/1108746637.geojson > 1108746637.png
```

Would produce:

![](images/1108746637.png)

As in: https://spelunker.whosonfirst.org/id/1108746637/

## See also

* https://github.com/srwiley/oksvg
* https://github.com/whosonfirst/go-whosonfirst-svg
