CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src/github.com/whosonfirst/go-whosonfirst-mysql; then rm -rf src/github.com/whosonfirst/go-whosonfirst-mysql; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-mysql
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-mysql/database
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-mysql/tables
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-mysql/utils
	cp -r database src/github.com/whosonfirst/go-whosonfirst-mysql/
	cp -r tables src/github.com/whosonfirst/go-whosonfirst-mysql/
	cp -r utils src/github.com/whosonfirst/go-whosonfirst-mysql/
	cp -r *.go src/github.com/whosonfirst/go-whosonfirst-mysql/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/go-sql-driver/mysql"
	@GOPATH=$(GOPATH) go get -u "github.com/twpayne/go-geom"
	@GOPATH=$(GOPATH) go get -u "github.com/go-ini/ini"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-index"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-log"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-names"
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/warning src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/tidwall src/github.com/
	rm -rf src/github.com/jteeuwen/go-bindata/testdata
	rm -rf src/github.com/whosonfirst/go-whosonfirst-index/vendor/github.com/whosonfirst/go-whosonfirst-mysql/

vendor-deps: rmdeps deps
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt database/*.go
	go fmt tables/*.go
	go fmt utils/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-mysql-index cmd/wof-mysql-index.go
