package main

import (
	"context"
	"gin-grpc/helper"
	"gin-grpc/services"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func main() {
	gwmux := runtime.NewServeMux()                                                   // 创建路由
	opt := []grpc.DialOption{grpc.WithTransportCredentials(helper.GetClientCreds())} // 指定客户端请求时使用的证书
	err := services.RegisterProdServiceHandlerFromEndpoint(
		context.Background(), gwmux, "localhost:8081", opt)
	if err != nil {
		log.Fatal(err)
	}
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: gwmux,
	}
	httpServer.ListenAndServe()
}
