package main

import (
	"gin-grpc/helper"
	"gin-grpc/services"
	"google.golang.org/grpc"
	"net"
)

func main() {
	rpcServer := grpc.NewServer(grpc.Creds(helper.GetServerCreds()))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))   // 商品服务
	services.RegisterOrderServiceServer(rpcServer, new(services.OrderService)) // 订单服务

	listen, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(listen)

	//mux := http.NewServeMux()
	//mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Println(request)
	//	rpcServer.ServeHTTP(writer, request)
	//})
	//httpServer := &http.Server{
	//	Addr:    ":8081",
	//	Handler: mux,
	//}
	//err := httpServer.ListenAndServeTLS("keys/client.pem", "keys/client.key")
	//if err != nil {
	//	log.Fatal(err)
	//}
}
