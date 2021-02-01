package scripts

const BuildScript = `
docker run --rm \
-v /home/custer/myweb:/app \
-v /home/custer/gopath:/go \
-w /app/src \
-e GOPROXY=https://goproxy.cn \
golang:1.15.7-alpine3.13 \
go build -o ../myserver main.go
`
