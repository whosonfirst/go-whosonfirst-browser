CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-render
	# cp -r cache src/github.com/whosonfirst/go-whosonfirst-render/
	cp -r assets src/github.com/whosonfirst/go-whosonfirst-render/
	cp -r http src/github.com/whosonfirst/go-whosonfirst-render/
	cp -r reader src/github.com/whosonfirst/go-whosonfirst-render/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-sdk-go"
	@GOPATH=$(GOPATH) go get -u "github.com/jteeuwen/go-bindata/"
	# @GOPATH=$(GOPATH) go get -u "github.com/patrickmn/go-cache"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-sanitize"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	rm -rf src/github.com/jteeuwen/go-bindata/testdata

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

assets: self
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/jteeuwen/go-bindata/go-bindata/
	rm -rf templates/*/*~
	rm -rf assets
	mkdir -p assets/html
	@GOPATH=$(GOPATH) bin/go-bindata -pkg html -o assets/html/html.go templates/html

fmt:
	# go fmt cache/*.go
	go fmt cmd/*.go
	go fmt assets/*.go
	go fmt http/*.go
	go fmt reader/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-renderd cmd/wof-renderd.go
