.PHONY: build build-alpine clean test help default

BIN_NAME=face-recognition-backend
ARM_BIN_NAME=face-recognition-backend

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "registry.jiangxingai.com:5000/face-recognition-backend"

default: build

help:
	@echo 'Management commands for face-recognition-backend:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make get-deps        runs dep ensure, mostly used for ci.'
	@echo '    make build-alpine    Compile optimized for alpine linux.'
	@echo '    make package         Build final docker image with just the go binary inside'
	@echo '    make tag             Tag image created by package with latest, git commit and version'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make push            Push tagged images to registry'
	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	export CGO_ENABLED=0; \
	go build -ldflags "-X gitlab.jiangxingai.com/luyor/face-recognition-backend/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X gitlab.jiangxingai.com/luyor/face-recognition-backend/version.BuildDate=${BUILD_DATE}" -o bin/${BIN_NAME}

build-arm:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	export CGO_ENABLED=0; \
	export GOARCH=arm64; \
	go build -ldflags "-X gitlab.jiangxingai.com/luyor/face-recognition-backend/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X gitlab.jiangxingai.com/luyor/face-recognition-backend/version.BuildDate=${BUILD_DATE}" -o bin/${ARM_BIN_NAME}

get-deps:
	dep ensure

build-alpine:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags '-w -linkmode external -extldflags "-static" -X gitlab.jiangxingai.com/luyor/face-recognition-backend/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X gitlab.jiangxingai.com/luyor/face-recognition-backend/version.BuildDate=${BUILD_DATE}' -o bin/${ARM_BIN_NAME}

package:
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build -t $(IMAGE_NAME):local .

tag: package
	@echo "Tagging: x86_v1.0.0 ${VERSION} $(GIT_COMMIT)"
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):x86_v1.0.0

push: tag
	@echo "Pushing docker image to registry: x86_v1.0.0 ${VERSION} $(GIT_COMMIT)"
	docker push $(IMAGE_NAME):x86_v1.0.0

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test ./...

