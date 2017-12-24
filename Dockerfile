# https://blog.docker.com/2016/09/docker-golang/
# https://blog.golang.org/docker

# docker build -t wof-staticd .

# For example:
# docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='http' -e HTTP_ROOT='https://whosonfirst.mapzen.com/data/' -e MAPZEN_APIKEY='mapzen-****' wof-staticd
#
# Or:
#
#
# docker run -it -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='s3' -e S3_BUCKET='whosonfirst' -e S3_PREFIX='' -e S3_REGION='us-east-1' -e S3_CREDENTIALS='env:' -e AWS_ACCESS_KEY_ID='***' -e AWS_SECRET_ACCESS_KEY='***' -e MAPZEN_APIKEY='mapzen-***' wof-staticd

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

