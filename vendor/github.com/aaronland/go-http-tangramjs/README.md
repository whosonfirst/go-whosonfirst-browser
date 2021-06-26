# go-http-tangramjs

`go-http-tangramjs` is an HTTP middleware package for including Tangram.js (v0.21.1) assets in web applications.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-http-tangramjs.svg)](https://pkg.go.dev/github.com/sfomuseum/go-http-tangramjs)

`go-http-tangramjs` is an HTTP middleware package for including Tangram.js assets in web applications. It exports two principal methods: 

* `tangramjs.AppendAssetHandlers(*http.ServeMux)` which is used to append HTTP handlers to a `http.ServeMux` instance for serving Tangramjs JavaScript files, and related assets.
* `tangramjs.AppendResourcesHandler(http.Handler, *TangramJSOptions)` which is used to rewrite any HTML produced by previous handler to include the necessary markup to load Tangramjs JavaScript files and related assets.

## Example

```
package main

import (
	"embed"
	"github.com/aaronland/go-http-tangramjs"
	"html/template"
	"log"
	"net/http"
)

//go:embed *.html
var FS embed.FS

func ExampleHandler(templates *template.Template) (http.Handler, error) {

	t := templates.Lookup("example")

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		err := t.Execute(rsp, nil)
		return
	}

	return http.HandlerFunc(fn), nil
}

func main() {

	api_key := "****"
	style_url := "/tangram/refill-style.zip"

	t, _ := template.ParseFS(FS, "*.html")

	mux := http.NewServeMux()

	tangramjs.AppendAssetHandlers(mux)
	
	map_handler, _:= ExampleHandler(t)

	tangramjs_opts := tangramjs.DefaultTangramJSOptions()
	tangramjs_opts.NextzenOptions.APIKey = api_key
	tangramjs_opts.NextzenOptions.StyleURL = style_url

	map_handler = tangramjs.AppendResourcesHandler(map_handler, tangramjs_opts)

	mux.Handle("/", map_handler)

	endpoint := "localhost:8080"
	log.Printf("Listening for requests on %s\n", endpoint)

	http.ListenAndServe(endpoint, mux)
}
```

_Error handling omitted for the sake of brevity._

You can see an example of this application by running the [cmd/example](cmd/example/main.go) application. You can do so by invoking the `example` Makefile target. For example:

```
$> make example APIKEY=XXXXXX
go run -mod vendor cmd/example/main.go -api-key XXXXXX
2021/05/05 15:31:06 Listening for requests on localhost:8080
```

The when you open the URL `http://localhost:8080` in a web browser you should see the following:

![](docs/images/go-http-tangramjs-example.png)

## See also

* https://github.com/tangrams/tangram
* https://github.com/aaronland/go-http-leaflet