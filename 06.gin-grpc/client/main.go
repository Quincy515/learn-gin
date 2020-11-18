package main

import (
	"context"
	"fmt"
	"gin-grpc/helper"
	"gin-grpc/services"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	t := timestamp.Timestamp{Seconds: time.Now().Unix()}
	orderClient := services.NewOrderServiceClient(conn)
	res, _ := orderClient.NewOrder(ctx, &services.OrderMain{
		OrderId:    1001,
		OrderNo:    "20201118",
		OrderMoney: 90,
		OrderTime:  &t,
	})
	fmt.Println(res)
	//prodClient := services.NewProdServiceClient(conn)
	//// 获取商品库存
	////prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 12, ProdArea: services.ProdAreas_B})
	//prodRes, err := prodClient.GetProdInfo(context.Background(), &services.ProdRequest{ProdId: 12})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Info(prodRes)
	//res, err := prodClient.GetProdStocks(context.Background(), &services.QuerySize{Size: 10})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(res.Prodres)
	//fmt.Println(res.Prodres[2].ProdStock)
}
