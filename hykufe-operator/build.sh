#!/bin/bash

version=$1

operator-sdk generate k8s
operator-sdk generate openapi
operator-sdk build yoowj7472/hykufe-operator:${version}
echo 'push to docker hub'
docker push yoowj7472/hykufe-operator:${version}
