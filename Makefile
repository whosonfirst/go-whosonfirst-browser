cli:
	go build -mod vendor -ldflags="-s -w" -o bin/whosonfirst-browser cmd/whosonfirst-browser/main.go

debug:
	go run cmd/whosonfirst-browser/main.go -enable-all -enable-edit \
	-reader-uri 'findingaid://https/static.sfomuseum.org/findingaid?template=https://raw.githubusercontent.com/sfomuseum-data/{repo}/main/data/' \
	-cache-uri null:// \
	-map-provider leaflet \
	-leaflet-tile-url 'https://tile.openstreetmap.org/{z}/{x}/{y}.png' \
	-leaflet-enable-draw \
	-authenticator-uri null:// \
	-writer-uri stdout:// \
	-spatial-database-uri 'pmtiles://?tiles=s3blob%3A%2F%2Fsfomuseum-tiles%3Fregion%3Dus-west-2%26prefix%3Dpoint-in-polygon%2F%26credentials%3Dsession&database=whosonfirst_sfom_13&enable-cache=true&zoom=13&layer=whosonfirst_sfom_13'

x-debug:
	# @make cli
	./bin/whosonfirst-browser -enable-all -map-provider tangram -nextzen-apikey $(APIKEY) -reader-uri $(READER)

debug-tilepack:
	make cli && bin/whosonfirst-browser -enable-all -nextzen-tilepack-database $(TILEPACK) -reader-uri $(READER)

debug-docker:
	docker run -it -p 8080:8080 whosonfirst-browser /usr/local/bin/whosonfirst-browser -server-uri 'http://0.0.0.0:8080' -enable-all -nextzen-api-key $(APIKEY)

lambda:
	@make lambda-browser

lambda-browser:
	if test -f main; then rm -f main; fi
	if test -f browser.zip; then rm -f browser.zip; fi
	GOOS=linux go build -mod vendor -ldflags="-s -w" -o main cmd/whosonfirst-browser/main.go
	zip browser.zip main
	rm -f main

docker:
	docker build -t whosonfirst-browser .

bump-version:
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g'
