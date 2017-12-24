#!/bin/sh

# This assumes 'ADD . /go-whosonfirst-static' which is defined in the Dockerfile

STATICD="/go-whosonfirst-static/bin/wof-staticd"
ARGS=""

if [ "${HOST}" != "" ]
then
    ARGS="${ARGS} -host ${HOST}"
fi

if [ "${MAPZEN_APIKEY}" != "" ]
then
    ARGS="${ARGS} -mapzen-apikey ${MAPZEN_APIKEY}"
fi

if [ "${TEST}" != "" ]
then
    ARGS="${ARGS} -test ${TEST}"
fi

if [ "${SOURCE}" = "http" ]
then

    if [ "${HTTP_ROOT}" = "" ]
    then
	echo "missing HTTP_ROOT"
	exit 1
    fi
    
    ARGS="${ARGS} -source http -http-root ${HTTP_ROOT}"
fi

if [ "${SOURCE}" = "s3" ]
then
    
    if [ "${S3_BUCKET}" = "" ]
    then
	echo "missing S3_BUCKET"
	exit 1
    fi

    ARGS="${ARGS} -source s3 -s3-bucket ${S3_BUCKET}"

    if [ "${S3_PREFIX}" != "" ]
    then
	ARGS="${ARGS} -s3-prefix ${S3_PREFIX}"
    fi

    if [ "${S3_REGION}" != "" ]
    then
	ARGS="${ARGS} -s3-region ${S3_REGION}"
    fi

    if [ "${S3_CREDENTIALS}" != "" ]
    then
	ARGS="${ARGS} -s3-credentials ${S3_CREDENTIALS}"
    fi
fi    
    
if [ "${SOURCE}" = "fs" ]
then
    
    if [ "${FS_ROOT}" = "" ]
    then
	echo "missing FS_ROOT"
	exit 1
    fi

    if [ ! -d ${FS_ROOT} ]
    then
	echo "FS_ROOT '${FS_ROOT}' does not exist"
	exit 1
    fi
    
    ARGS="${ARGS} -source fs -fs-root ${FS_ROOT}"
fi

# echo ${STATICD} ${ARGS}

${STATICD} ${ARGS}

if [ $? -ne 0 ]
then
   echo "command '${STATICD} ${ARGS}' failed"
   exit 1
fi

exit 0
