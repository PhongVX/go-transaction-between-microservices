package productx

import (
	"context"
	"github.com/PhongVX/micro-protos/product"
	"orchestrator/internal/transactionx"
)

type RepositoryI interface {
	UpdateProduct(context.Context, *product.UpdateProductRequest) (*product.UpdateProductResponse, error)
}

func NewRepository(txSrv transactionx.ServiceI) *Repository{
	return &Repository{
		txSrv: txSrv,
	}
}

func (r *Repository)UpdateProduct(ctx context.Context, p *product.UpdateProductRequest) (*product.UpdateProductResponse, error){
	tx, err := r.txSrv.GetTxByCorrelationID(p.CorrelationID)
	if err != nil {
		return nil, err
	}
	queryString := `UPDATE product SET quantity=(quantity - $1) WHERE id=$2`
	rs, err := tx.Exec(queryString, p.Quantity, p.ID)
	if err != nil {
		return &product.UpdateProductResponse{RowAffected: int32(0)}, err
	}
	if val, err := rs.RowsAffected(); err == nil && val > 0 {
		return &product.UpdateProductResponse{RowAffected: int32(val)}, nil
	}
	return &product.UpdateProductResponse{RowAffected: int32(0)}, err
}
