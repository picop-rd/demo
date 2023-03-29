#!/bin/bash -ex
PROXY_CONTROLLER_URL=$1

# proxy-a
curl ${PROXY_CONTROLLER_URL}/proxy/service-a/activate -XPUT
# proxy-b
curl ${PROXY_CONTROLLER_URL}/proxy/service-b/activate -XPUT
# proxy-c
curl ${PROXY_CONTROLLER_URL}/proxy/service-c/activate -XPUT
# proxy-mysql
curl ${PROXY_CONTROLLER_URL}/proxy/service-mysql/activate -XPUT

