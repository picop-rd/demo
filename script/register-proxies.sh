#!/bin/bash -ex
PROXY_CONTROLLER_URL=$1

# proxy-a
curl ${PROXY_CONTROLLER_URL}/proxy/service-a/register -XPUT -H "Content-Type: application/json" -d "{\"endpoint\":\"http://proxy-a.service-a.svc.cluster.local:9000\"}"
# proxy-b
curl ${PROXY_CONTROLLER_URL}/proxy/service-b/register -XPUT -H "Content-Type: application/json" -d "{\"endpoint\":\"http://proxy-b.service-b.svc.cluster.local:9000\"}"
# proxy-c
curl ${PROXY_CONTROLLER_URL}/proxy/service-c/register -XPUT -H "Content-Type: application/json" -d "{\"endpoint\":\"http://proxy-c.service-c.svc.cluster.local:9000\"}"
# proxy-mysql
curl ${PROXY_CONTROLLER_URL}/proxy/service-mysql/register -XPUT -H "Content-Type: application/json" -d "{\"endpoint\":\"http://proxy-mysql.service-mysql.svc.cluster.local:9000\"}"

