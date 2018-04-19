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

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/zendesk/go-bindata/"
	@GOPATH=$(GOPATH) go get -u "github.com/elazarl/go-bindata-assetfs/"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-http-mapzenjs"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-s3/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-fs/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-http/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-spr"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-svg"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
	rm -rf src/github.com/zendesk/go-bindata/testdata
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
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/zendesk/go-bindata/go-bindata/
	rm -rf templates/*/*~
	rm -rf assets
	mkdir -p assets/html
	@GOPATH=$(GOPATH) bin/go-bindata -pkg html -o assets/html/html.go templates/html

static: self
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/zendesk/go-bindata/go-bindata/
	@GOPATH=$(GOPATH) go build -o bin/go-bindata-assetfs vendor/github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs/main.go
	rm -f static/css/*~ static/javascript/*~ static/tangram/*~ static/fonts/*~
	@PATH=$(PATH):$(CWD)/bin bin/go-bindata-assetfs -pkg http static/javascript static/css static/tangram static/fonts
	if test -f http/static_fs.go; then rm http/static_fs.go; fi
	mv bindata.go http/static_fs.go

build:
	@make static
	@make assets
	@make bin

wof: wof-fonts wof-css

wof-fonts:
	if test ! -d static/fonts; then mkdir -p static/fonts; fi
	curl -s -o static/fonts/Poppins-Light.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Poppins-Light.ttf
	curl -s -o static/fonts/Poppins-Medium.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Poppins-Medium.ttf
	curl -s -o static/fonts/Poppins-SemiBold.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Poppins-SemiBold.ttf
	curl -s -o static/fonts/Roboto-Light.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Roboto-Light.ttf
	curl -s -o static/fonts/Roboto-LightItalic.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Roboto-LightItalic.ttf
	curl -s -o static/fonts/Roboto-Regular.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/Roboto-Regular.ttf
	curl -s -o static/fonts/RobotoMono-Light.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/RobotoMono-Light.ttf
	curl -s -o static/fonts/glyphicons-halflings-regular.eot https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/glyphicons-halflings-regular.eot
	curl -s -o static/fonts/glyphicons-halflings-regular.svg https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/glyphicons-halflings-regular.svg
	curl -s -o static/fonts/glyphicons-halflings-regular.ttf https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/glyphicons-halflings-regular.ttf
	curl -s -o static/fonts/glyphicons-halflings-regular.woff https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/fonts/glyphicons-halflings-regular.woff

wof-css:
	if test ! -d static/css; then mkdir -p static/css; fi
	curl -s -o static/css/whosonfirst.www.css https://raw.githubusercontent.com/whosonfirst/whosonfirst-www/master/www/css/mapzen.whosonfirst.css

maps: mapzenjs tangram refill crosshairs

localforage:
	curl -s -o static/javascript/localforage.js https://raw.githubusercontent.com/mozilla/localForage/master/dist/localforage.js
	curl -s -o static/javascript/localforage.min.js https://raw.githubusercontent.com/mozilla/localForage/master/dist/localforage.min.js

tangram:
	if test ! -d static/tangram; then mkdir -p static/tangram; fi
	curl -s -o static/javascript/tangram.js https://www.nextzen.org/tangram/tangram.debug.js
	curl -s -o static/javascript/tangram.min.js https://www.nextzen.org/tangram/tangram.min.js

refill:
	if test ! -d static/tangram; then mkdir -p static/tangram; fi
	curl -s -o static/tangram/refill-style.zip https://www.nextzen.org/carto/refill-style/refill-style.zip

# mapzen.js has not been moved over to nextzen yet (https://www.nextzen.org/)

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
	bin/wof-staticd -port 8080 -source http -http-root https://data.whosonfirst.org -cache lru -cache-arg 'CacheSize=500' -debug -mapzen-apikey ${NEXTZEN_APIKEY}

debug-local: build
	bin/wof-staticd -port 8080 -source fs -fs-root /usr/local/data/whosonfirst-data/data -cache bigcache -cache-arg HardMaxCacheSize=100 -cache-arg MaxEntrySize=1024 -debug -mapzen-apikey ${NEXTZEN_APIKEY}

docker-build:
	docker build -t wof-static .

docker-debug: docker-build
	docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='http' -e HTTP_ROOT='https://whosonfirst.mapzen.com/data/' -e CACHE='gocache' -e CACHE_ARGS='DefaultExpiration=300 CleanupInterval=600' -e DEBUG='debug' -e NEXTZEN_APIKEY=${NEXTZEN_APIKEY} wof-static