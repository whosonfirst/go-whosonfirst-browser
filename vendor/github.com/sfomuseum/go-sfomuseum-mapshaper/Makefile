cli:
	go build -mod vendor -o bin/server cmd/server/main.go

docker:
	docker build -t mapshaper-server .

debug-docker:
	docker run -it -p 8080:8080 -e MAPSHAPER_SERVER_URI=http://0.0.0.0:8080 mapshaper-server
