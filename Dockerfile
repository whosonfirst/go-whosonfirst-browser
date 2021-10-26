FROM golang:1.17-alpine as builder

RUN mkdir /build

COPY . /build/go-whosonfirst-browser

RUN apk update && apk upgrade \
    && apk add libc-dev gcc \
    && cd /build/go-whosonfirst-browser \
    && go build -mod vendor -o /usr/local/bin/whosonfirst-browser cmd/whosonfirst-browser/main.go    

FROM alpine:latest

COPY --from=builder /usr/local/bin/whosonfirst-browser /usr/local/bin/

RUN mkdir -p /usr/local/data/tiles

COPY tiles/*.db /usr/local/data/tiles/

RUN apk update && apk upgrade \
    && apk add ca-certificates