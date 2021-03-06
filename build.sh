#! /bin/bash

VERSION=${1}
ARCH=${2}
BUILDFILE="Dockerfile"

if [ "x${VERSION}" -eq "x" ]; then
    VERSION="1.0.0"
fi

if [ "x${ARCH}" -eq "x" ]; then
    ARCH=`uname -i`
fi

if [ "${ARCH}" == "x86_64" ]; then
    BUILDFILE="Dockerfile"
    ARCH="x8664"
    make build
elif [ "${ARCH}" == "aarch64" ]; then
    BUILDFILE="Dockerfile.arm64"
    ARCH="arm64v8"
    make build-arm
else
    echo "unkown cpu arch"
    exit 1
fi

repo="registry.jiangxingai.com:5000/face-recognition-backend:${ARCH}_cpu_${VERSION}"
latest_repo="registry.jiangxingai.com:5000/face-recognition-backend:${ARCH}_cpu_latest"

sudo docker build . -t ${repo} -f ${BUILDFILE}
sudo docker tag ${repo} ${latest_repo}
sudo docker push ${repo}
sudo docker push ${latest_repo}
