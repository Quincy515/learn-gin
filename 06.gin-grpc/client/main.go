package main

import (
	"context"
	"gin-grpc/services"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("keys/test.pem", "*.custer.fun")
	if err != nil {
		log.Fatal(err)
	}
	conn,err:=grpc.Dial(":8081",grpc.WithTransportCredentials(creds))

//	creds, err := credentials.NewClientTLSFromFile("keys/grpc.crt", "custer.fun")
//	if err != nil {
//		log.Fatal(err)
//	}
//	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := services.NewProdServiceClient(conn)
	prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 12})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(prodRes.ProdStock)
}
