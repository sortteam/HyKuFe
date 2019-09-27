#!/bin/bash
set -x
[ "$#" -eq 1 ] || { echo 'version is required e.g. v0.0.1'; exit 1; }
version=$1

export GO111MODULE=on

operator-sdk generate k8s
operator-sdk generate openapi
operator-sdk build hykufe/hg-operator:${version}
echo 'push to docker hub'
docker push hykufe/hg-operator:${version}
