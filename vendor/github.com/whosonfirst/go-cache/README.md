# go-cache

Go interface for caching things

_This package replaces [go-whosonfirst-cache](https://github.com/whosonfirst/go-whosonfirst-cache) which will be retired soon._

## Install

You will need to have both `Go` (version [Go 1.12](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make tools
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Interfaces

```
type Cache interface {
     	Name() string
	Open(context.Context, string) error
	Close(context.Context) error
	Get(context.Context, string) (io.ReadCloser, error)
	Set(context.Context, string, io.ReadCloser) (io.ReadCloser, error)
	Unset(context.Context, string) error
	Hits() int64
	Misses() int64
	Evictions() int64
	Size() int64
	SizeWithContext(context.Context) int64
}
```

## See also

