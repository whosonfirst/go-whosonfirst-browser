# go-whosonfirst-writer

## Important

Work in progress. Documentation to follow

## Interfaces

### WriterInitializationFunc

```
type WriterInitializationFunc func(ctx context.Context, uri string) (Writer, error)
```

### Writer

```
type Writer interface {
	Write(context.Context, string, io.ReadSeeker) (int64, error)
	Close() error
	WriterURI(string) string
}
```

## See also

* https://github.com/whosonfirst/go-reader
