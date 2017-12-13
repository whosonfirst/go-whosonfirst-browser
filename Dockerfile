# https://blog.docker.com/2016/09/docker-golang/
# https://blog.golang.org/docker

# docker build -t wof-renderd .
# docker run -p 6161:8080 -e AWS_ACCESS_KEY_ID='your-aws-key' -e AWS_SECRET_ACCESS_KEY='your-aws-secret' -e HOST='0.0.0.0' -e S3_BUCKET='example.com' -e S3_PREFIX='' -e S3_REGION='us-east-1' -e S3_CREDENTIALS='env:' wof-renderd

FROM golang

ADD . /go-whosonfirst-render

RUN cd /go-whosonfirst-render; make bin

EXPOSE 8080

# PLEASE FOR TO ADDING CONDITIONAL FLAGS IN CASE WE'RE RUNNING THIS AGAINST A "-source fs"
# SETUP OR EQUIVALENT... (20171213/thisisaaronland)

CMD /go-whosonfirst-render/bin/wof-renderd -host ${HOST} -source s3 -s3-bucket "${S3_BUCKET}" -s3-prefix "${S3_PREFIX}" -s3-region "${S3_REGION}" -s3-credentials "${S3_CREDENTIALS}"
