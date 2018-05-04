CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-readwrite-mysql
	cp -r reader src/github.com/whosonfirst/go-whosonfirst-readwrite-mysql/
	cp -r writer src/github.com/whosonfirst/go-whosonfirst-readwrite-mysql/
	cp -r utils src/github.com/whosonfirst/go-whosonfirst-readwrite-mysql/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-mysql"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
	rm -rf src/github.com/whosonfirst/go-whosonfirst-mysql/vendor/github.com/whosonfirst/go-whosonfirst-index/

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt reader/*.go
	go fmt writer/*.go
	go fmt utils/*.go
	go fmt cmd/*.go

bin: 	self
	rm -rf bin/*
	GOPATH=$(GOPATH) go build -o bin/wof-mysql-readerd cmd/wof-mysql-readerd.go
