#!/bin/bash -e

environment=$1

cat ./template/service-a.yaml | sed -e "s/\$(ENV)/$environment/g"
