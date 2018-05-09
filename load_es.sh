#!/bin/sh

if [ -z "$ES_URL" ] ; then
    echo ES_URL must be set
    exit 1
fi
ES_AUTH_HEADER=""
CURL_DELETE_INDEX_CMD='curl -X DELETE -H "Content-Type: application/x-ndjson"'
CURL_POST_CMD='curl -X POST -H "Content-Type: application/x-ndjson"'
if [ -n "$ES_USER" ] ; then
    echo "Credentials provided - adding auth header for ES"
    CURL_POST_CMD="${CURL_POST_CMD} -H \"Authorization: $(echo $ES_USER:$ES_PASS | base64)\""
    CURL_DELETE_INDEX_CMD="${CURL_DELETE_INDEX_CMD} -H \"Authorization: $(echo $ES_USER:$ES_PASS | base64)\""
else
    echo "No Credentials provided - accessing ES anonymously"
fi
CURL_DELETE_INDEX_CMD="${CURL_DELETE_INDEX_CMD} ${ES_URL}/cities"
CURL_POST_CMD="${CURL_POST_CMD} --data-binary @bulk-cities-es.txt ${ES_URL}/_bulk"
echo "Deleting index if exists"
eval $CURL_DELETE_INDEX_CMD
echo
echo "Loading Data"
eval $CURL_POST_CMD
