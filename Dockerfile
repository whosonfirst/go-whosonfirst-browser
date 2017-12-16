# https://blog.docker.com/2016/09/docker-golang/
# https://blog.golang.org/docker

# docker build -t wof-staticd .

# For example:
# docker run -p 6161:8080 -e HOST='0.0.0.0' -e SOURCE='http' -e HTTP_ROOT='https://whosonfirst.mapzen.com/data/' -e MAPZEN_APIKEY='mapzen-****' wof-staticd

# Or:
#
# docker run -p 6161:8080 -e HOST='0.0.0.0' -e S3_BUCKET='example.com' -e S3_PREFIX='' -e S3_REGION='us-east-1' -e S3_CREDENTIALS='iam:' -e MAPZEN_APIKEY-'your-mapzen-apikey' wof-staticd

FROM golang

ADD . /go-whosonfirst-static

RUN cd /go-whosonfirst-static; make bin

EXPOSE 8080

ENTRYPOINT /go-whosonfirst-static/docker/entrypoint.sh

