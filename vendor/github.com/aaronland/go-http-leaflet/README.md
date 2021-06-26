# go-http-leaflet

`go-http-leaflet` is an HTTP middleware package for including Leaflet.js (v1.7.1) assets in web applications.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-http-leaflet.svg)](https://pkg.go.dev/github.com/aaronland/go-http-leaflet)

`go-http-leaflet` is an HTTP middleware package for including Leaflet.js assets in web applications. It exports two principal methods:

* `leaflet.AppendAssetHandlers(*http.ServeMux)` which is used to append HTTP handlers to a `http.ServeMux` instance for serving Leaflet JavaScript files, and related assets.
* `leaflet.AppendResourcesHandler(http.Handler, *LeafletOptions)` which is used to rewrite any HTML produced by previous handler to include the necessary markup to load Leaflet

This package doesn't specify any code or methods for how Leaflet.js is used. It only provides method for making Leaflet.js available to existing applications.

By default this package only appends assets and resources for Leaflet.js but it also includes the necessary assets to enable the use of the [leaflet-hash](https://github.com/mlevans/leaflet-hash), [Leaflet.fullscreen](https://github.com/Leaflet/Leaflet.fullscreen) and [Leaflet.Draw](https://leaflet.github.io/Leaflet.draw/docs/leaflet-draw-latest.html) plugins. These are enabled by invoking the corresponding `EnableHash`, `EnableFullscreen` and `EnableDraw` methods on the `leaflet.LeafletOptions` instance.

## Example

```
package main

import (
        "embed"
	"flag"
	"github.com/aaronland/go-http-leaflet"
	"log"
	"net/http"
)

//go:embed *.html
var FS embed.FS

func ExampleHandler(templates *template.Template) (http.Handler, error) {

	t := templates.Lookup("example")

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		t.Execute(rsp)
		return
	}

	return http.HandlerFunc(fn), nil
}

func main() {

	t, _ := template.ParseFS(FS, "*.html")

	http.NewServeMux()

	leaflet_opts := leaflet.DefaultLeafletOptions()
	leaflet_opts.EnableHash()
	leaflet_opts.EnableFullscreen()
	leaflet_opts.EnableDraw()

	example_handler, _ := ExampleHandler(t)
	example_handler = leaflet.AppendResourcesHandler(example_handler, leaflet_opts)

	mux.Handle("/", example_handler)

	endpoint := "localhost:8080"
	log.Printf("Listening for requests on %s\n", endpoint)

	http.ListenAndServe(endpoint, mux)
}
```

_Error handling omitted for brevity._

You can see an example of this application by running the [cmd/example](cmd/example/main.go) application. You can do so by invoking the `example` Makefile target. For example:

```
$> make example
go run -mod vendor cmd/example/main.go -enable-hash -enable-fullscreen -enable-draw -tile-url 'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
2021/05/05 10:53:53 Listening for requests on localhost:8080
```

The when you open the URL `http://localhost:8080` in a web browser you should see the following:

![](docs/images/go-http-leaflet-example.png)

### Notes

The default `leaflet.draw.css` causes sprites to be misaligned, at least when using Firefox. I have found the following CSS declaration to fix the problem:

```
.leaflet-draw-toolbar a {
	background-size: 300px 30px !important;
}
```

## See also

* https://leafletjs.com
* https://github.com/mlevans/leaflet-hash
* https://github.com/Leaflet/Leaflet.fullscreen
* https://leaflet.github.io/Leaflet.draw/docs/leaflet-draw-latest.html