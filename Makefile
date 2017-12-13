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
	@GOPATH=$(GOPATH) go get -u "github.com/elazarl/go-bindata-assetfs/"
	# @GOPATH=$(GOPATH) go get -u "github.com/patrickmn/go-cache"
	# @GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-sanitize"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-http-mapzenjs"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-http-rewrite"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
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

static: self
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/jteeuwen/go-bindata/go-bindata/
	@GOPATH=$(GOPATH) go build -o bin/go-bindata-assetfs vendor/github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs/main.go
	rm -f static/css/*~ static/javascript/*~ static/tangram/*~
	@PATH=$(PATH):$(CWD)/bin bin/go-bindata-assetfs -pkg http static/javascript static/css static/tangram
	rm http/static_fs.go
	mv bindata_assetfs.go http/static_fs.go

build: static assets bin

maps: mapzenjs tangram refill crosshairs

tangram:
	if test ! -d static/tangram; then mkdir -p static/tangram; fi
	curl -s -o static/javascript/tangram.js https://mapzen.com/tangram/tangram.debug.js
	curl -s -o static/javascript/tangram.min.js https://mapzen.com/tangram/tangram.min.js

refill:
	if test ! -d static/tangram; then mkdir -p static/tangram; fi
	curl -s -o static/tangram/refill-style.zip https://mapzen.com/carto/refill-style/refill-style.zip

mapzenjs:
	if test ! -d static/javascript; then mkdir -p static/javascript; fi
	if test ! -d static/css; then mkdir -p static/css; fi
	curl -s -o static/css/mapzen.js.css https://mapzen.com/js/mapzen.css
	curl -s -o static/javascript/mapzen.js https://mapzen.com/js/mapzen.js
	curl -s -o static/javascript/mapzen.min.js https://mapzen.com/js/mapzen.min.js

crosshairs:
	if test ! -d static/javascript; then mkdir -p static/javascript; fi
	curl -s -o static/javascript/slippymap.crosshairs.js https://raw.githubusercontent.com/whosonfirst/js-slippymap-crosshairs/master/src/slippymap.crosshairs.js	

fmt:
	# go fmt cache/*.go
	go fmt cmd/*.go
	go fmt assets/*/*.go
	go fmt http/*.go
	go fmt reader/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-renderd cmd/wof-renderd.go
