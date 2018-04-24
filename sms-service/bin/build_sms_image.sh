#!/bin/bash
DIRNAME=`dirname $0`
DOCKER_BUILD_DIR=`cd $DIRNAME/; pwd`
echo "DOCKER_BUILD_DIR=${DOCKER_BUILD_DIR}"
cd ${DOCKER_BUILD_DIR}

BUILD_ARGS="--no-cache"
ORG="onap"
VERSION="2.0.0"
PROJECT="aaf"
IMAGE="sms"
DOCKER_REPOSITORY="nexus3.onap.org:10003"
IMAGE_NAME="${DOCKER_REPOSITORY}/${ORG}/${PROJECT}/${IMAGE}"
TIMESTAMP=$(date +"%Y%m%dT%H%M%S")

if [ $HTTP_PROXY ]; then
    BUILD_ARGS+=" --build-arg HTTP_PROXY=${HTTP_PROXY}"
fi
if [ $HTTPS_PROXY ]; then
    BUILD_ARGS+=" --build-arg HTTPS_PROXY=${HTTPS_PROXY}"
fi

function generate_binary {
    pushd ../src/sms
    make build
    popd
    cp ../target/sms .
}

function remove_binary {
    rm sms
}

function build_image {
    echo "Start build docker image: ${IMAGE_NAME}"
    docker build ${BUILD_ARGS} -t ${IMAGE_NAME}:latest -f smsdockerfile .
}

function push_image_tag {
    TAG_NAME=$1
    echo "Start push ${TAG_NAME}"
    docker tag ${IMAGE_NAME}:latest ${TAG_NAME}
    docker push ${TAG_NAME}
}

function push_image {
    echo "Start push ${IMAGE_NAME}:latest"
    docker push ${IMAGE_NAME}:latest

    push_image_tag ${IMAGE_NAME}:${VERSION}-SNAPSHOT-latest
}

generate_binary
build_image
push_image
remove_binary