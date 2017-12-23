CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-static
	# cp -r cache src/github.com/whosonfirst/go-whosonfirst-static/
	cp -r assets src/github.com/whosonfirst/go-whosonfirst-static/
	cp -r http src/github.com/whosonfirst/go-whosonfirst-static/
	cp -r utils src/github.com/whosonfirst/go-whosonfirst-static/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

docker-build:
	docker build -t wof-static .

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/jteeuwen/go-bindata/"
	@GOPATH=$(GOPATH) go get -u "github.com/elazarl/go-bindata-assetfs/"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-http-mapzenjs"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-spr"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-svg"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
	rm -rf src/github.com/jteeuwen/go-bindata/testdata
	rm -rf src/github.com/whosonfirst/go-whosonfirst-svg/vendor/github.com/whosonfirst/go-whosonfirst-geojson-v2
	rm -rf src/github.com/whosonfirst/go-whosonfirst-svg/vendor/github.com/whosonfirst/go-whosonfirst-flags
	rm -rf src/github.com/whosonfirst/go-whosonfirst-svg/vendor/github.com/whosonfirst/go-whosonfirst-spr
	rm -rf src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/go-whosonfirst-flags

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
	rm -f static/css/*~ static/javascript/*~ static/tangram/*~ static/fonts/*~
	@PATH=$(PATH):$(CWD)/bin bin/go-bindata-assetfs -pkg http static/javascript static/css static/tangram static/fonts
	rm http/static_fs.go
	mv bindata_assetfs.go http/static_fs.go

build: static assets bin

wof: wof-fonts wof-css

wof-fonts:
	if test ! -d static/fonts; then mkdir -p static/fonts; fi
	curl -s -o static/fonts/Poppins-Light.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Poppins-Light.ttf
	curl -s -o static/fonts/Poppins-Medium.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Poppins-Medium.ttf
	curl -s -o static/fonts/Poppins-SemiBold.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Poppins-SemiBold.ttf
	curl -s -o static/fonts/Roboto-Light.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Roboto-Light.ttf
	curl -s -o static/fonts/Roboto-LightItalic.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Roboto-LightItalic.ttf
	curl -s -o static/fonts/Roboto-Regular.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Roboto-Regular.ttf
	curl -s -o static/fonts/Roboto-Mono-Light.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Roboto-Mono-Light.ttf
	curl -s -o static/fonts/glyphicons-halflings-regular.eot https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/glyphicons-halflings-regular.eot
	curl -s -o static/fonts/glyphicons-halflings-regular.svg https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/glyphicons-halflings-regular.svg
	curl -s -o static/fonts/glyphicons-halflings-regular.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/glyphicons-halflings-regular.ttf
	curl -s -o static/fonts/glyphicons-halflings-regular.woff https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/glyphicons-halflings-regular.woff

wof-css:
	if test ! -d static/css; then mkdir -p static/css; fi
	curl -s -o static/css/whosonfirst.www.css https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/css/mapzen.whosonfirst.css

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
	go fmt utils/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-staticd cmd/wof-staticd.go

debug: build
	bin/wof-staticd -port 8080 -source http -http-root https://data.whosonfirst.org -mapzen-apikey ${MAPZEN_APIKEY}
