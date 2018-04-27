package http

import (
       "errors"
       "fmt"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-static/utils"
	"github.com/whosonfirst/go-whosonfirst-image"
	gohttp "net/http"
	"strings"
)

func RasterHandler(r reader.Reader, format string) (gohttp.Handler, error) {

     format = strings.ToLower(format)

     switch format {
     	    case "png":
	    	 // pass
            default:
		return nil, errors.New("Invalid or unsupported raster format")
		}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		f, err, status := utils.FeatureFromRequest(req, r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), status)
			return
		}

		content_type := fmt.Sprintf("image/%s", format)
		rsp.Header().Set("Content-Type", content_type)

		// THIS INTERFACE _WILL_ CHANGE SO DON'T GET TOO USED TO IT
		// (20180427/thisisaaronland)

		image.FeatureToPNG(f, rsp)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
