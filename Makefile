cli:
	go build -mod vendor -o bin/whosonfirst-browser cmd/whosonfirst-browser/main.go

debug:
	go run -mod vendor cmd/whosonfirst-browser/main.go -enable-graphics -enable-data -enable-html -proxy-tiles -nextzen-api-key $(APIKEY)

debug-docker:
	docker run -it -p 8080:8080 whosonfirst-browser /usr/local/bin/whosonfirst-browser -server-uri 'http://0.0.0.0:8080' -enable-all -nextzen-api-key $(APIKEY)

lambda:
	@make lambda-browser

lambda-browser:
	if test -f main; then rm -f main; fi
	if test -f browser.zip; then rm -f browser.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/whosonfirst-browser/main.go
	zip browser.zip main
	rm -f main

docker:
	docker build -t whosonfirst-browser .

bump-version:
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/whosonfirst\/go-whosonfirst-browser\/$(PREVIOUS)/github.com\/whosonfirst\/go-whosonfirst-browser\/$(NEW)/g'
