#!/bin/bash
DIRNAME=`dirname $0`
DOCKER_BUILD_DIR=`cd $DIRNAME/; pwd`
cd ${DOCKER_BUILD_DIR}

(cd ../src/sms && make build)
cp ../target/sms .

sudo ./build_image.sh
