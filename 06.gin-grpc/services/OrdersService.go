package services

import (
	"context"
	"fmt"
)

type OrderService struct{}

func (this *OrderService) NewOrder(ctx context.Context, orderRequest *OrderRequest) (*OrderResponse, error) {
	err := orderRequest.OrderMain.Validate()
	if err != nil {
		return &OrderResponse{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}
	fmt.Println(orderRequest.OrderMain)
	return &OrderResponse{
		Status:  "ok",
		Message: "success",
	}, nil
}
