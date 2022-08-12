# go-whosonfirst-image

Use `go-whosonfirst-svg` and `oksvg` to render Who's On First features as raster images.

## Example

```
import (
	"github.com/whosonfirst/go-whosonfirst-image"
	"image/png"
	"io"
	"os"
)

func main() {

     	path := "/path/to/feature.geojson"

	r, _ := os.Open(path)
	defer r.Close()
	
	body, err := io.ReadAll(r)

	opts := image.NewDefaultOptions()
	img, _ := image.FeatureToImage(body, opts)

	png.Encode(opts.Writer, img)
}
```

_Error handling removed for brevity._

## Tools

### wof-feature-to-png

```
./bin/wof-feature-to-png /usr/local/data/whosonfirst-data-constituency-us/data/110/874/663/7/1108746637.geojson > 1108746637.png
```

Would produce:

![](images/1108746637.png)

As in: https://spelunker.whosonfirst.org/id/1108746637/

## See also

* https://github.com/srwiley/oksvg
* https://github.com/whosonfirst/go-whosonfirst-svg
