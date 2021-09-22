package orderx

import (
	"context"
	"github.com/PhongVX/micro-protos/order"
)

type GServiceI interface {
	InsertOrder(ctx context.Context, in *order.InsertOrderRequest) (*order.InsertOrderResponse, error)
	InsertOrderDetail(context.Context, *order.InsertOrderDetailRequest) (*order.InsertOrderDetailResponse, error)
}

func NewGService(repos RepositoryI) *GService{
	return &GService{
		repos: repos,
	}
}

func (s *GService)InsertOrder(ctx context.Context, o *order.InsertOrderRequest) (*order.InsertOrderResponse, error) {
	return s.repos.InsertOrder(ctx, o)
}

func (s *GService) InsertOrderDetail(ctx context.Context, o *order.InsertOrderDetailRequest) (*order.InsertOrderDetailResponse, error){
	return s.repos.InsertOrderDetail(ctx, o)
}