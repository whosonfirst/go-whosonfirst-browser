CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test ! -d src; then mkdir src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-image
	cp *.go src/github.com/whosonfirst/go-whosonfirst-image/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/srwiley/oksvg"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-svg"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/warning"
	mv src/github.com/whosonfirst/go-whosonfirst-svg/vendor/github.com/whosonfirst/go-whosonfirst-geojson-v2 src/github.com/whosonfirst
	mv src/github.com/whosonfirst/go-whosonfirst-svg/vendor/github.com/whosonfirst/go-whosonfirst-flags src/github.com/whosonfirst
	mv src/github.com/whosonfirst/go-whosonfirst-svg/vendor/github.com/whosonfirst/go-whosonfirst-spr src/github.com/whosonfirst

# if you're wondering about the 'rm -rf' stuff below it's because Go is
# weird... https://vanduuren.xyz/2017/golang-vendoring-interface-confusion/
# (20170912/thisisaaronland)

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src


fmt:
	go fmt *.go
	go fmt cmd/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-feature2png cmd/wof-feature2png.go
