package services

import "context"

type ProdService struct{}

func (this *ProdService) GetProdStock(ctx context.Context, request *ProdRequest) (*ProdResponse, error) {
	var stock int32 = 0
	if request.ProdArea == ProdAreas_A {
		stock = 39
	} else if request.ProdArea == ProdAreas_B {
		stock = 41
	} else {
		stock = 20
	}
	return &ProdResponse{ProdStock: stock}, nil
}

func (this *ProdService) GetProdStocks(context.Context, *QuerySize) (*ProdResponseList, error) {
	Prodres := []*ProdResponse{
		&ProdResponse{ProdStock: 28},
		&ProdResponse{ProdStock: 29},
		&ProdResponse{ProdStock: 30},
		&ProdResponse{ProdStock: 31},
	}
	return &ProdResponseList{Prodres: Prodres}, nil
}
