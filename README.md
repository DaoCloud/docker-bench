# docker-bench

## Container & Service
    go run scale-out.go common.go -r 2 -n 5 -i daocloud.io/nginx
    go run scale-in.go common.go -r 2 -n 5 -i daocloud.io/nginx
    go run create-service.go common.go -r 5 -n 1 -i daocloud.io/nginx
    go run create-container.go common.go -c 2 -r 3 -i daocloud.io/nginx

## Network
    ./ping.sh -c 5 8.8.8.8

    DOCKER_OPTIONS="--network host" ./nginx.sh
    DOCKER_OPTIONS="--network host" ./ab.sh -c 5 -n 20 127.0.0.1/4K
    DOCKER_OPTIONS="--network host" ./ab.sh -c 5 -n 20 127.0.0.1/4M
    DOCKER_OPTIONS="--network host" ./ab.sh -c 5 -n 20 127.0.0.1/40M

    DOCKER_OPTIONS="--network host" ./netserver.sh
    DOCKER_OPTIONS="--network host" ./netperf.sh -H 127.0.0.1 -l 5 -f M
