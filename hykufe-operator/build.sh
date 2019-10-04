#!/bin/bash
set -x
# [ "$#" -eq 1 ] || { echo 'version is required e.g. v0.0.1'; exit 1; }
images=$(echo $IMAGES | tr " " "\n")

operator-sdk generate k8s
operator-sdk generate openapi
operator-sdk build $images
echo 'push to docker hub'
docker push $images
