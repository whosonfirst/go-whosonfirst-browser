# go-artisanal-integers

No, really.


## Install

You will need to have both `Go` and the `make` programs installed on your computer. Assuming you do just type:

```
make tools
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Usage

## Interfaces

### Client

```
type Client interface {
	NextInt() (int64, error)
}
```

### Engine

An "engine" is the interface between your code and an underlying data model (typically a database) for minting artisanal integers. The interface looks like this:

```
type Engine interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
	SetLastInt(int64) error
	SetKey(string) error
	SetOffset(int64) error
	SetIncrement(int64) error
	Close() error
}
```

### Service

```
type Service interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
}
```

### Server

```
type Server interface {
	ListenAndServe(Service) error
}
```

## Tools

Everything is in flux. This will be updated soon.

## Engines

* https://github.com/aaronland/go-artisanal-integers-mysql
* https://github.com/aaronland/go-artisanal-integers-redis
* https://github.com/aaronland/go-artisanal-integers-rqlite
* https://github.com/aaronland/go-artisanal-integers-summitdb

## See also

* http://www.brooklynintegers.com/
* http://www.londonintegers.com/
* http://www.neverendingbooks.org/artisanal-integers
* https://nelsonslog.wordpress.com/2012/07/29/artisinal-integers/
* https://nelsonslog.wordpress.com/2012/08/25/artisinal-integers-part-2/
* http://www.aaronland.info/weblog/2012/12/01/coffee-and-wifi/#timepixels
* https://mapzen.com/blog/mapzen-acquires-mission-integers/
* http://code.flickr.net/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/
