#!/bin/bash
set -x
[ "$#" -eq 1 ] || { echo 'version is required e.g. v0.0.1'; exit 1; }
version=$1

operator-sdk generate k8s
operator-sdk generate openapi
operator-sdk build yoowj7472/hykufe-operator:${version}
echo 'push to docker hub'
docker push yoowj7472/hykufe-operator:${version}
