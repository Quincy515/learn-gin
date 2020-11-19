package main

import (
	"context"
	"fmt"
	"gin-grpc/helper"
	"gin-grpc/services"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
)

func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	userClient := services.NewUserServiceClient(conn)
	var i int32
	req := services.UserScoreRequest{}
	req.Users = make([]*services.UserInfo, 0)

	for i = 1; i < 20; i++ {
		req.Users = append(req.Users, &services.UserInfo{UserId: i})
	}
	stream, err := userClient.GetUserScoreByServerStream(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break // 读取到结尾就 break
		}
		if err != nil { // 读取失败，就停止程序运行
			log.Fatal(err)
		}
		fmt.Println(res.Users)
	}
	//t := timestamp.Timestamp{Seconds: time.Now().Unix()}
	//orderClient := services.NewOrderServiceClient(conn)
	//res, _ := orderClient.NewOrder(ctx, &services.OrderRequest{
	//	OrderMain: &services.OrderMain{
	//		OrderId:    1001,
	//		OrderNo:    "20201118",
	//		OrderMoney: 90,
	//		OrderTime:  &t,
	//	},
	//})
	//fmt.Println(res)
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
