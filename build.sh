#!/usr/bin/env bash

GIT_REVISION=$(git rev-parse --short --verify HEAD)
TIME=$(date -u +%Y%m%d.%H%M%S)
VERSION=1.0.${GIT_REVISION}.${TIME}

build_artifacts () {
  local os=$1
  local arch=$2
  GO111MODULE=on GOOS=$os GOARCH=$arch go build -ldflags "-X main.versionString=${VERSION}" ./cmd/crypta2/
  file crypta2
  tar -czvf crypta2-"$os"-x64.tar.gz crypta2
  rm crypta2
}

build_artifacts linux amd64
build_artifacts darwin amd64
