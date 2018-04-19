#!/bin/sh

# This assumes 'ADD . /go-whosonfirst-static' which is defined in the Dockerfile

STATICD="/go-whosonfirst-static/bin/wof-staticd"
ARGS=""

if [ "${HOST}" != "" ]
then
    ARGS="${ARGS} -host ${HOST}"
fi

if [ "${NEXTZEN_APIKEY}" != "" ]
then
    ARGS="${ARGS} -mapzen-apikey ${NEXTZEN_APIKEY}"
fi

if [ "${TEST}" != "" ]
then
    ARGS="${ARGS} -test ${TEST}"
fi

if [ "${SOURCE}" = "" ]
then
    echo "missing SOURCE"
    exit 1
fi

if [ "${SOURCE_DSN}" = "" ]
then
    echo "missing SOURCE_DSN"
    exit 1
fi

ARGS="${ARGS} -source ${SOURCE} -source-dsn ${SOURCE_DSN}"

if [ "${CACHE}" != "" ]
then

    ARGS="${ARGS} -cache ${CACHE}"

    if [ "${CACHE_ARGS}" != "" ]
    then

	for CA in ${CACHE_ARGS}
	do
	    ARGS="${ARGS} -cache-arg ${CA}"
	done
    fi
    
fi

if [ "${DEBUG}" != "" ]
then	       
    ARGS="${ARGS} -debug"
fi
   
# echo ${STATICD} ${ARGS}

${STATICD} ${ARGS}

if [ $? -ne 0 ]
then
   echo "command '${STATICD} ${ARGS}' failed"
   exit 1
fi

exit 0
