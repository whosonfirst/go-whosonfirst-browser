# https://blog.docker.com/2016/09/docker-golang/
# https://blog.golang.org/docker

# docker build -t wof-staticd .

# See the way we're passing around ENV variables and calling '-s3-credentials "env:' because by the
# time this gets invoked in Docker-land we are a different "person" (20171215/thisisaaronland)
# docker run -p 6161:8080 -e HOST='0.0.0.0' -e S3_BUCKET='example.com' -e S3_PREFIX='' -e S3_REGION='us-east-1' -e S3_CREDENTIALS='env:' -e AWS_ACCESS_KEY_ID='your-aws-key' -e AWS_SECRET_ACCESS_KEY='your-aws-secret' -e MAPZEN_APIKEY-'your-mapzen-apikey' wof-staticd

FROM golang

ADD . /go-whosonfirst-static

RUN cd /go-whosonfirst-static; make bin

EXPOSE 8080

# HOW DO I MAKE THIS A CONDITIONAL ?

# CMD /go-whosonfirst-render/bin/wof-staticd -host ${HOST} -source fs -fs-root "${FS_ROOT}" -mapzen-apikey "${MAPZEN_APIKEY}"
# CMD /go-whosonfirst-render/bin/wof-staticd -host ${HOST} -source http -http-root "${HTTP_ROOT}" -mapzen-apikey "${MAPZEN_APIKEY}"

CMD /go-whosonfirst-static/bin/wof-staticd -host ${HOST} -source s3 -s3-bucket "${S3_BUCKET}" -s3-prefix "${S3_PREFIX}" -s3-region "${S3_REGION}" -s3-credentials "${S3_CREDENTIALS}" -mapzen-apikey "${MAPZEN_APIKEY}"
