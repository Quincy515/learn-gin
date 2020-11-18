protoc --go_out=plugins=grpc:../services Prod.proto
protoc --go_out=plugins=grpc:../services Orders.proto
protoc \
  -I . \
  -I ${GOPATH}/src \
  -I ${GOPATH}/src/github.com/protoc-gen-validate \
  --go_out=plugins=grpc:../services \
  --validate_out=lang=go:../services \
  Models.proto

protoc --grpc-gateway_out=logtostderr=true:../services Prod.proto
protoc --grpc-gateway_out=logtostderr=true:../services Orders.proto


