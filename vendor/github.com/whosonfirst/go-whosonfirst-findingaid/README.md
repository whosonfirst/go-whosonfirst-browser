# go-whosonfirst-findingaid

A Go language interface for building and querying finding aids of Who's On First documents.

## Documentation

Documentation is incomplete.

## FindingAids

Conceptually a finding aid consists of two parts:

* An indexer which indexes (or catalogs) one or more Who's On First (WOF) records in to a cache. WOF records may be cataloged in full, truncated or otherwise manipulated according to logic implemented by the indexing or caching layers.
* A cache of WOF records, in full or otherwise manipulated, that can resolved using a given WOF ID.

It is generally assumed that a complete catalog of WOF records will be assembled in advance of any query actions but that is not an absolute requirement. For an example of a lazy-loading catalog and query implementation, where all operations are performed at runtime, consult the documentation for the `readercache` chaching layer below.

There can be more than one kind of finding aid. Finding aids can implement their own internal logic for cataloging, caching and querying WOF records. A finding aid need only implement the following interface:

```
type FindingAid interface {
     Resolver
     Indexer
}

type Indexer interface {
	IndexURIs(context.Context, ...string) error
	IndexReader(context.Context, io.Reader) error
}

type Resolver interface {
	ResolveURI(context.Context, string) (interface{}, error)
}
```

Note the ambiguous return value (`interface{}`) for the `ResolveURI` method. Since it impossible to know in advance the response properties of any given finding aid it is left to developers to cast query results in to the appropriate type if necessary.

...

## See also

* https://github.com/whosonfirst/go-cache
* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://en.wikipedia.org/wiki/Finding_aid