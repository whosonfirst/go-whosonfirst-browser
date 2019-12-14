# go-whosonfirst-svg

Tools for working with Who's On First and SVG documents.

## Install

You will need to have both `Go` (specifically a version [1.12](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make tools
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### wof2svg

For example:

```
go run -mod vendor cmd/wof2svg/main.go sfo.geojson sfo.svg
```

Will produce:

![](docs/images/sfo.png)

And:

```
go run -mod vendor cmd/wof2svg/main.go --mercator sfo.geojson sfo-merc.svg
```

Will produce:

![](docs/images/sfo-mercator.png)


## See also

* https://github.com/whosonfirst/go-geojson-svg
