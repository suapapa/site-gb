#!/bin/bash

set -e

IMAGE_TAG=gcr.io/homin-dev/gb 
docker buildx build --platform linux/amd64 --build-arg=PROGRAM_VER=$1 -t $IMAGE_TAG:$1 .
docker push $IMAGE_TAG:$1

IMAGE_TAG_LATEST=$IMAGE_TAG:latest
docker tag $IMAGE_TAG:$1 $IMAGE_TAG_LATEST
docker push $IMAGE_TAG_LATEST

git tag -a $1 -m "add tag for $1"
git push --tags
