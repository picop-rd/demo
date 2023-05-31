#!/bin/bash -ex
PROXY_CONTROLLER_URL=$1

curl ${PROXY_CONTROLLER_URL}/routes -XPUT -H "Content-Type: application/json" -d "
[
	{\"proxy_id\":\"service-a\",     \"env_id\": \"feature-1\", \"destination\": \"service-a-feature-1.service-a.svc.cluster.local:80\"},
	{\"proxy_id\":\"service-mysql\", \"env_id\": \"feature-1\", \"destination\": \"192.168.0.3:13307\"},
	{\"proxy_id\":\"service-b\",     \"env_id\": \"feature-2\", \"destination\": \"service-b-feature-2.service-b.svc.cluster.local:80\"},
	{\"proxy_id\":\"service-mysql\", \"env_id\": \"feature-2\", \"destination\": \"192.168.0.3:13308\"}
]
"
