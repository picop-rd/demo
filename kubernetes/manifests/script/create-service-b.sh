#!/bin/bash -e

environment=$1

cat ./template/service-b.yaml | sed -e "s/\$(ENV)/$environment/g"
