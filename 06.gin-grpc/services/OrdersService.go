package services

import (
	"context"
	"fmt"
)

type OrderService struct{}

func (this *OrderService) NewOrder(ctx context.Context, orderRequest *OrderRequest) (*OrderResponse, error) {
	fmt.Println(orderRequest.OrderMain)
	return &OrderResponse{
		Status:  "ok",
		Message: "success",
	}, nil
}
