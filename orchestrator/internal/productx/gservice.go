package productx

import (
	"context"

	"github.com/PhongVX/micro-protos/product"
)

type GServiceI interface {
	UpdateProduct(context.Context, *product.UpdateProductRequest) (*product.UpdateProductResponse, error)
}

func NewGService(repos RepositoryI) *GService{
	return &GService{
		repos: repos,
	}
}

func (s *GService)UpdateProduct(ctx context.Context, p *product.UpdateProductRequest) (*product.UpdateProductResponse, error){
	return s.repos.UpdateProduct(ctx, p)
}
