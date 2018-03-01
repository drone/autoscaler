#!/bin/sh

set -e
set -x

TAG=${DRONE_TAG=latest}
SHA=${DRONE_COMMIT_SHA:0:7}

LDFLAGS="-extldflags '-static' -X main.version=${TAG} -X main.commit=${SHA}"

go build -ldflags ${LDFLAGS} \
	-o release/linux/arm64/drone-autoscaler \
	github.com/drone/autoscaler/cmd/drone-autoscaler
