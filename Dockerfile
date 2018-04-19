# https://blog.docker.com/2016/09/docker-golang/
# https://blog.golang.org/docker

# docker build -t wof-static .

# For example:
# docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='http' -e SOURCE_DSN='https://whosonfirst.mapzen.com/data/' -e NEXTZEN_APIKEY='mapzen-****' wof-staticd
#
# Or:
#
# docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='s3' -e SOURCE_DSN='bucket=whosonfirst region=us-east-1 credentials=env:' -e AWS_ACCESS_KEY_ID='***' -e AWS_SECRET_ACCESS_KEY='***' -e NEXTZEN_APIKEY='mapzen-***' wof-staticd
#
# Or even still, with caching:
# docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='http' -e SOURCE_DSN='https://whosonfirst.mapzen.com/data/' -e CACHE='gocache' -e CACHE_ARGS='DefaultExpiration=300 CleanupInterval=600' -e DEBUG='debug' -e NEXTZEN_APIKEY=${NEXTZEN_APIKEY} wof-static

# build phase - see also:
# https://medium.com/travis-on-docker/multi-stage-docker-builds-for-creating-tiny-go-images-e0e1867efe5a
# https://medium.com/travis-on-docker/triple-stage-docker-builds-with-go-and-angular-1b7d2006cb88

FROM golang:alpine AS build-env

# https://github.com/gliderlabs/docker-alpine/issues/24

RUN apk add --update alpine-sdk

ADD . /go-whosonfirst-static

RUN cd /go-whosonfirst-static; make bin

# bundle phase - note the way we need certificates

FROM alpine

RUN apk add --update ca-certificates

WORKDIR /go-whosonfirst-static/bin/

COPY --from=build-env /go-whosonfirst-static/bin/wof-staticd /go-whosonfirst-static/bin/wof-staticd
COPY --from=build-env /go-whosonfirst-static/docker/entrypoint.sh /go-whosonfirst-static/bin/entrypoint.sh

EXPOSE 8080

ENTRYPOINT /go-whosonfirst-static/bin/entrypoint.sh

