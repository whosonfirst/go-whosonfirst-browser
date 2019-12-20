fmt:
	go fmt *.go
	go fmt application/*.go
	go fmt client/*.go
	go fmt cmd/*.go
	go fmt engine/*.go
	go fmt http/*.go
	go fmt server/*.go
	go fmt service/*.go
	go fmt utils/*.go

tools:
	if test ! -d bin; then mkdir bin; fi
	go build -o bin/int cmd/int/main.go
	go build -o bin/intd-client cmd/intd-client/main.go
	go build -o bin/intd-server cmd/intd-server/main.go