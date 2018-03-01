#!/bin/sh

TAG=${DRONE_TAG=latest}
SHA=${DRONE_COMMIT_SHA}

LDFLAGS="-extldflags '-static' -X main.version=${TAG} -X main.commit=${SHA}"

set -e
set -x

go build -ldflags ${LDFLAGS} \
	-o release/linux/arm64/drone-autoscaler \
	github.com/drone/autoscaler/cmd/drone-autoscaler
