FROM golang:1.12-alpine as builder

RUN mkdir /build

COPY . /build/go-whosonfirst-browser

RUN apk update && apk upgrade \
    && apk add make git \
    && cd /build/go-whosonfirst-browser \
    && make bake \
    && go build -mod vendor -o /usr/local/bin/whosonfirst-browser cmd/whosonfirst-browser/main.go    

FROM alpine:latest

COPY --from=builder /usr/local/bin/whosonfirst-browser /usr/local/bin/

RUN apk update && apk upgrade \
    && apk add ca-certificates