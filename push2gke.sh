#!/bin/bash
IMAGE_TAG=gcr.io/homin-dev/gb:latest 
docker buildx build --platform linux/amd64 -t $IMAGE_TAG .
docker push $IMAGE_TAG