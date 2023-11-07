GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/whosonfirst-browser cmd/whosonfirst-browser/main.go

debug:
	go run -mod $(GOMOD) cmd/whosonfirst-browser/main.go \
		-enable-all \
		-map-provider leaflet \
		-leaflet-tile-url https://tile.openstreetmap.org/{z}/{x}/{y}.png \
		-javascript-at-eof \
		-rollup-assets \
		-reader-uri $(READER) \
		-authenticator-uri null:// \
		-enable-edit \
		-verbose

debug-tilepack:
	make cli && bin/whosonfirst-browser -enable-all -nextzen-tilepack-database $(TILEPACK) -reader-uri $(READER)

debug-docker:
	docker run -it -p 8080:8080 whosonfirst-browser /usr/local/bin/whosonfirst-browser -server-uri 'http://0.0.0.0:8080' -enable-all -nextzen-api-key $(APIKEY)

lambda:
	@make lambda-browser

lambda-browser:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f browser.zip; then rm -f browser.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/whosonfirst-browser/main.go
	zip browser.zip bootstrap
	rm -f bootstrap

docker:
	docker build -t whosonfirst-browser .

bump-version:
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g'
