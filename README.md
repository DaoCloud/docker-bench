# docker-bench

## examples
    go run scale-out.go common.go -r 2 -n 5 -i daocloud.io/nginx
    go run scale-in.go common.go -r 2 -n 5 -i daocloud.io/nginx
    go run create-service.go common.go -r 5 -n 1 -i daocloud.io/nginx
    go run create-container.go common.go -c 2 -r 3 -i daocloud.io/nginx
