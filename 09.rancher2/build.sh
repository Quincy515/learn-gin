echo "start build (linux,amd64)"
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
go build -o build/myserver main.go
echo "complete build (linux,amd64)";