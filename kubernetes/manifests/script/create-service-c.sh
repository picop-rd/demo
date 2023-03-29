#!/bin/bash -e

environment=$1

cat ./template/service-c.yaml | sed -e "s/\$(ENV)/$environment/g"
