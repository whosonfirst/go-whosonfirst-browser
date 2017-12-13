# go-whosonfirst-render

## Important

Stop. This is too soon for you. Really. You should assume that everything about this package (including the name) will change.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Known knowns and other things "to figure out"

* This package is already starting to get littered with both rendering and delivering code, for example all of the static assets necessary to show a WOF document rendered as a HTML document. It is probably the case that we should have two packages with the "delivery" package modifying the rendered HTML with instance-specific CSS and the like. Or not.

* The "reader" code (and by extension the caching layer) in this package should probably be moved in to its own package.
