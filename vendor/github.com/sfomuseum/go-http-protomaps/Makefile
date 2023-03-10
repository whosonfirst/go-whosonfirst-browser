cli:
	go build -mod vendor -o bin/example cmd/example/main.go

debug:
	go run -mod vendor cmd/example/main.go -javascript-at-eof

protomaps-js:
	curl -s -L -o static/javascript/protomaps.js https://unpkg.com/protomaps@latest/dist/protomaps.js 
	curl -s -L -o static/javascript/protomaps.min.js https://unpkg.com/protomaps@latest/dist/protomaps.min.js 
