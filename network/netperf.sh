#!/bin/bash

docker run ${DOCKER_OPTIONS} daocloud.io/haipeng/netperf /usr/bin/netperf $*
