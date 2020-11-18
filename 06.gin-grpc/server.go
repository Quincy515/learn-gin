package main

import (
	"fmt"
	"gin-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile("keys/test.pem", "keys/test.key")
	if err != nil {
		log.Fatal(err)
	}
	rpcServer := grpc.NewServer(grpc.Creds(creds))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))

	//listen, _ := net.Listen("tcp", ":8081")
	//rpcServer.Serve(listen)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request)
		rpcServer.ServeHTTP(writer, request)
	})
	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	httpServer.ListenAndServeTLS("keys/test.pem", "keys/test.key")
}
