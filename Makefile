PWD := $(shell pwd)
APP := attach
PKG := github.com/ckeyer/$(APP)
CMS_PKG := github.com/ckeyer/commons
GO := CGO_ENABLED=0 GOBIN=${PWD}/bundles go
HASH := $(shell which sha1sum || which shasum)

DEV_IMAGE := ckeyer/dev:go
IMAGE_NAME := ckeyer/$(APP):$(GIT_BRANCH)

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
VERSION := $(shell cat VERSION.txt)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_AT := $(shell date "+%Y-%m-%dT%H:%M:%SZ%z")
PACKAGE_NAME := $(APP)$(VERSION).$(OS)-$(ARCH)

LD_FLAGS := -X $(CMS_PKG)/version.version=$(VERSION) \
 -X $(CMS_PKG)/version.gitCommit=$(GIT_COMMIT) \
 -X $(CMS_PKG)/version.buildAt=$(BUILD_AT) -w

init:
	echo $(IMAGE_NAME)
	which govendor || go get github.com/kardianos/govendor
	govendor sync

local: generate
	$(GO) install -a -ldflags="$(LD_FLAGS)" .

generate:
	$(GO) generate ./protos

build:
	docker run --rm \
	 --name $(APP)-building \
	 -e CGO_ENABLED=0 \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) make local

run:
	docker run --rm \
	 --name $(APP)-dev-running \
	 -e CGO_ENABLED=0 \
	 -p 8089:8080 \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) ${GO} run ./main.go

image: build
	docker build -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

test:
	${GO} test -ldflags="$(LD_FLAGS)" $$(go list ./... |grep -v "vendor")

dev: dev-server

dev-server:
	docker run --rm -it \
	 --name $(APP)-dev \
	 -p 8089:8080 \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) bash

dev-client:
	docker run --rm -it \
	 --name $(APP)-dev-client \
	 -v /var/run/docker.sock:/var/run/docker.sock \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) bash
