CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-readwrite-s3
	cp -r config src/github.com/whosonfirst/go-whosonfirst-readwrite-s3/
	cp -r reader src/github.com/whosonfirst/go-whosonfirst-readwrite-s3/
	cp -r writer src/github.com/whosonfirst/go-whosonfirst-readwrite-s3/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite/..."
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-sdk-go"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt config/*.go
	go fmt reader/*.go
	go fmt writer/*.go

bin: 	self
	GOPATH=$(GOPATH) go build -o bin/wof-s3-copy cmd/wof-s3-copy.go
