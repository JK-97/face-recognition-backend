#! /bin/bash

VERSION="1.0.0"
BUILDFILE="Dockerfile"
ARCH=`uname -i`

if [ "${ARCH}" == "x86_64" ]; then
    BUILDFILE="Dockerfile"
    ARCH="x8664"

elif [ "${ARCH}" == "aarch64" ]; then
    BUILDFILE="Dockerfile.arm64"
    ARCH="arm64v8"
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
