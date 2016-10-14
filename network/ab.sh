#!/bin/bash

docker run ${DOCKER_OPTIONS} daocloud.io/haipeng/ab /usr/bin/ab $*
