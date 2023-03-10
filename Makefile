GOMOD=vendor

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/whosonfirst-browser cmd/whosonfirst-browser/main.go

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
	GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -o main cmd/whosonfirst-browser/main.go
	zip browser.zip main
	rm -f main

# https://github.com/tdewolff/minify

rollup:
	@make rollup-css
	@make rollup-js

rollup-js:
	if test -f static/javascript/whosonfirst.browser.rollup.js; then rm static/javascript/whosonfirst.browser.rollup.js; fi
	minify static/javascript/*.js -o static/javascript/whosonfirst.browser.rollup.js

rollup-css:
	if test -f static/css/whosonfirst.browser.rollup.css; then rm static/css/whosonfirst.browser.rollup.css; fi
	minify static/css/*.css -o static/css/whosonfirst.browser.rollup.css

docker:
	docker build -t whosonfirst-browser .

bump-version:
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g'
