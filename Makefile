CWD=$(shell pwd)

go-bindata:
	mkdir -p cmd/go-bindata
	mkdir -p cmd/go-bindata-assetfs
	curl -s -o cmd/go-bindata/main.go https://raw.githubusercontent.com/whosonfirst/go-bindata/master/cmd/go-bindata/main.go
	curl -s -o cmd/go-bindata-assetfs/main.go https://raw.githubusercontent.com/whosonfirst/go-bindata-assetfs/master/cmd/go-bindata-assetfs/main.go

debug:
	@make bake
	go run -mod vendor cmd/browser/main.go -enable-all -nextzen-api-key $(APIKEY)

bake: bake-static bake-templates

bake-static:
	go build -o bin/go-bindata cmd/go-bindata/main.go
	go build -o bin/go-bindata-assetfs cmd/go-bindata-assetfs/main.go
	rm -f www/static/*~ www/static/css/*~ www/static/javascript/*~
	@PATH=$(PATH):$(CWD)/bin bin/go-bindata-assetfs -prefix www -pkg http www/static/javascript www/static/css www/static/fonts

bake-templates:
	mv bindata.go http/assetfs.go
	rm -rf templates/html/*~
	bin/go-bindata -pkg templates -o assets/templates/html.go www/templates/html

