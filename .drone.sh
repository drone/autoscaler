#!/bin/sh

set -e
set -x

go build \
    -ldflags '-extldflags "-static"' \
	-ldflags "-X main.version=${DRONE_TAG=latest}" \
	-ldflags "-X main.commit=${DRONE_COMMIT_SHA}" \
	-o release/linux/arm64/drone-autoscaler \
	github.com/drone/autoscaler/cmd/drone-autoscaler
