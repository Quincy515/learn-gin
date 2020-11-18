package main

import (
	"context"
	"fmt"
	"gin-grpc/helper"
	"gin-grpc/services"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := services.NewProdServiceClient(conn)
	//prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 12})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Info(prodRes.ProdStock)
	res, err := prodClient.GetProdStocks(context.Background(), &services.QuerySize{Size: 10})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Prodres)
	fmt.Println(res.Prodres[2].ProdStock)
}
