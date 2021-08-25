cli:
	go build -mod vendor -o bin/lookupd cmd/lookupd/main.go
	go build -mod vendor -o bin/catalog cmd/catalog/main.go
	go build -mod vendor -o bin/resolve cmd/resolve/main.go

debug:
	go run -mod vendor cmd/lookupd/main.go

lambda-handlers:
	@make lambda-server

lambda-server:	
	if test -f main; then rm -f main; fi
	if test -f lookupd.zip; then rm -f lookupd.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/lookupd/main.go
	zip lookupd.zip main
	rm -f main

docker-lookupd:
	docker build -f Dockerfile.lookupd -t findingaid-lookupd .		
