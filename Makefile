CWD=$(shell pwd)

cli:
	go build -mod vendor -o bin/whosonfirst-browser cmd/whosonfirst-browser/main.go

go-bindata:
	mkdir -p cmd/go-bindata
	mkdir -p cmd/go-bindata-assetfs
	curl -s -o cmd/go-bindata/main.go https://raw.githubusercontent.com/whosonfirst/go-bindata/master/cmd/go-bindata/main.go
	curl -s -o cmd/go-bindata-assetfs/main.go https://raw.githubusercontent.com/whosonfirst/go-bindata-assetfs/master/cmd/go-bindata-assetfs/main.go

debug:
	@make bake
	go run -mod vendor cmd/whosonfirst-browser/main.go -enable-all -proxy-tiles -nextzen-api-key $(APIKEY)

lambda:
	@make lambda-browser

lambda-browser:
	if test -f main; then rm -f main; fi
	if test -f browser.zip; then rm -f browser.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/whosonfirst-browser/main.go
	zip browser.zip main
	rm -f main

bake: bake-static bake-templates

bake-static:
	go build -mod vendor -o bin/go-bindata cmd/go-bindata/main.go
	go build -mod vendor -o bin/go-bindata-assetfs cmd/go-bindata-assetfs/main.go
	rm -f www/static/*~ www/static/css/*~ www/static/javascript/*~
	@PATH=$(PATH):$(CWD)/bin bin/go-bindata-assetfs -prefix www -pkg http www/static/javascript www/static/css www/static/fonts

bake-templates:
	mv bindata.go http/assetfs.go
	rm -rf templates/html/*~
	bin/go-bindata -pkg templates -o assets/templates/html.go www/templates/html

docker:
	docker build -t whosonfirst-browser .

bump-version:
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g'
