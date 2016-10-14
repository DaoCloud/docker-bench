#!/bin/bash

docker run ${DOCKER_OPTIONS} daocloud.io/haipeng/busybox ping $*
